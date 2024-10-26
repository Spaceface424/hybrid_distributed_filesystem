package swim

import (
	"bufio"
	"context"
	"cs425/mp2/shared"
	"fmt"
	"log"
	"math/rand/v2"
	"net"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
)

var (
	ENABLE_SUSPECT          = false
	DROP_RATE       float64 = 0.0
	SUSPECT_TIMEOUT float64 = 2.0
	PROTOCAL_PERIOD float64 = 2.0
	K               int     = 1
	TTL             int32   = 4
	M               int     = 10
)

var (
	next_id        int32 = 0
	members        memberContainer
	failed_members memberContainer
	gossips        gossipContainer
	cur_member     *shared.MemberInfo
	round_robin    []int32
	rr_index       int
	ack_chan       chan *shared.PingAck
	verbose        bool
)

func UDPServer(host string, introducer string, verbose bool) {
	conn := setupServer(host, verbose)
	defer conn.Close()
	requestIntroducer(cur_member, introducer)
	setupLogging()

	membershipList()
	go sendUDPServer(conn, ack_chan)
	go listenUDPServer(conn, ack_chan)

	// user command loop
	scanner := bufio.NewReader(os.Stdin)
	for {
		text, _ := scanner.ReadBytes('\n')
		command := string(text[:len(text)-1])

		switch commandParts := strings.Split(command, " "); commandParts[0] {
		case "list_mem":
			printMembershipList()
		}
	}
}

func listenUDPServer(conn *net.UDPConn, ack_chan chan *shared.PingAck) {
	buf := make([]byte, 1024)
	for {
		n, addr, err := conn.ReadFromUDP(buf)
		pb := &shared.PingAck{}
		err = proto.Unmarshal(buf[:n], pb)
		if err != nil {
			log.Println("Read err:", err)
			continue
		}
		if rand.Float64() < DROP_RATE {
			log.Printf("[INFO] Packet dropped while receiving %s from node %d", pb.Type, pb.SenderId)
			continue
		}
		switch pb.Type {
		case shared.MessageType_ACK:
			log.Printf("[INFO] Received ACK from node %d", pb.SenderId)
			ack_chan <- pb
		case shared.MessageType_PING:
			log.Printf("[INFO] Received PING from node %d", pb.SenderId)
			_, ok := members.memberMap[pb.SenderId]
			_, fail_ok := failed_members.memberMap[pb.SenderId]
			if !ok && !fail_ok {
				members.memberMap[pb.SenderId] = &shared.MemberInfo{
					Address: addr.String(),
					ID:      pb.SenderId,
					State:   shared.NodeState_ALIVE,
					IncNum:  pb.IncNum,
				}
			} else if !ok && fail_ok {
				gossips.gossipMap[pb.SenderId] = &shared.Gossip{
					Member: &shared.MemberInfo{
						Address: addr.String(),
						ID:      pb.SenderId,
						State:   shared.NodeState_FAILED,
						IncNum:  0,
					},
					TTL: 1,
				}
			}
			go sendPingAck(pingackRequest{
				conn:         conn,
				target_addr:  addr,
				target_id:    pb.SenderId,
				sender_id:    cur_member.ID,
				round:        pb.Round,
				message_type: shared.MessageType_ACK,
			})
		case shared.MessageType_IDR_PING:
			log.Printf("[INFO] Received IDR_PING from node %d", pb.SenderId)
			go sendPingAck(pingackRequest{
				conn:         conn,
				target_addr:  addr,
				target_id:    pb.SenderId,
				sender_id:    cur_member.ID,
				round:        pb.Round,
				message_type: shared.MessageType_ACK_REQ,
				req_id:       pb.RequestId,
			})
		case shared.MessageType_PING_REQ:
			if _, ok := members.memberMap[pb.RequestId]; !ok {
				break
			}
			log.Printf("[INFO] Received PING_REQ from node %d to PING node %d", pb.SenderId, pb.RequestId)
			target_addr, _ := net.ResolveUDPAddr("udp", members.memberMap[pb.RequestId].Address)
			go sendPingAck(pingackRequest{
				conn:         conn,
				target_addr:  target_addr,
				target_id:    pb.RequestId,
				sender_id:    cur_member.ID,
				round:        pb.Round,
				message_type: shared.MessageType_IDR_PING,
				req_id:       pb.SenderId,
			})
		case shared.MessageType_ACK_REQ:
			if _, ok := members.memberMap[pb.RequestId]; !ok {
				break
			}
			log.Printf("[INFO] Received ACK_REQ from node %d, fowarding ACK to node %d", pb.SenderId, pb.RequestId)
			target_addr, _ := net.ResolveUDPAddr("udp", members.memberMap[pb.RequestId].Address)
			go sendPingAck(pingackRequest{
				conn:         conn,
				target_addr:  target_addr,
				target_id:    pb.RequestId,
				sender_id:    pb.SenderId,
				round:        pb.Round,
				message_type: shared.MessageType_ACK,
			})

		}
		if len(pb.GossipBuffer) > 0 {
			log.Printf("[INFO] Incoming gossip: %v", fmtGossip(&pb.GossipBuffer))
		}
		go updateMembership(pb)
	}
}

func sendUDPServer(conn *net.UDPConn, ack_chan chan *shared.PingAck) {
	var round int32 = 0
	round_timer := time.NewTimer(time.Duration(PROTOCAL_PERIOD) * time.Second)
	ping_req_timer := time.NewTimer(time.Duration(PROTOCAL_PERIOD/2) * time.Second)
	for {
		log.Printf("[INFO] Round: %d", round)
		logMembershipList()
		decTTL()
		round++
		if len(members.memberMap) == 0 {
			log.Printf("[INFO] No members in group...")
			round_robin = []int32{}
			round_timer.Reset(time.Duration(PROTOCAL_PERIOD) * time.Second)
			<-round_timer.C
			continue
		}

		target_id := getRoundRobinTarget()
		target_addr, err := net.ResolveUDPAddr("udp", members.memberMap[target_id].Address)
		if err != nil {
			log.Panic("Dial UDP err:", err)
		}

		round_timer.Reset(time.Duration(PROTOCAL_PERIOD) * time.Second)
		ping_req_timer.Reset(time.Duration(PROTOCAL_PERIOD/2) * time.Second)
		sendPingAck(pingackRequest{
			conn:         conn,
			target_addr:  target_addr,
			target_id:    target_id,
			sender_id:    cur_member.ID,
			round:        round,
			message_type: shared.MessageType_PING,
			inc_num:      cur_member.IncNum,
		})

		period := true
		for period {
			select {
			case pingack := <-ack_chan:
				if pingack.Round != round {
					continue
				}
				period = false
				<-round_timer.C
			case <-ping_req_timer.C:
				ping_req_timer.Stop()
				sendKPingReq(conn, target_id, round)
				continue
			case <-round_timer.C:
				if ENABLE_SUSPECT {
					pingTimeoutSuspicion(target_id)
				} else {
					pingTimeout(target_id)
				}
				period = false
			}
		}
	}
}

func sendKPingReq(conn *net.UDPConn, target_id int32, round int32) {
	k_i := 0
	for _, req_id := range round_robin {
		if req_id == target_id {
			continue
		}
		member, ok := members.memberMap[req_id]
		if !ok {
			continue
		}
		log.Printf("[INFO] PING not received in time from node %d, sending PING_REQ to node %d", target_id, req_id)
		req_addr, _ := net.ResolveUDPAddr("udp", member.Address)
		go sendPingAck(pingackRequest{
			conn:         conn,
			target_id:    req_id,
			target_addr:  req_addr,
			sender_id:    cur_member.ID,
			round:        round,
			message_type: shared.MessageType_PING_REQ,
			req_id:       target_id,
		})
		k_i++
		if k_i == 2 {
			break
		}
	}
}

func sendPingAck(req pingackRequest) {
	if len(gossips.gossipMap) > 0 {
		log.Printf("[INFO] Sending gossip: %v", fmtGossip(&gossips.gossipMap))
	}
	pb := &shared.PingAck{
		Type:         req.message_type,
		SenderId:     req.sender_id,
		Round:        req.round,
		GossipBuffer: gossips.gossipMap,
		RequestId:    req.req_id,
	}

	data, err := proto.Marshal(pb)
	buf := []byte(data)
	_, err = req.conn.WriteToUDP(buf, req.target_addr)
	if err != nil {
		log.Println(err)
	}
	log.Printf("[INFO] Sending %s to node %d, addr: %s", req.message_type, req.target_id, req.target_addr)
}

func pingTimeout(target_id int32) {
	log.Printf("[INFO] ACK not received in time from node %d, marking as FAILED", target_id)
	fmt.Printf("[INFO] ACK not received in time from node %d, marking as FAILED\n", target_id)

	if _, ok := failed_members.memberMap[target_id]; ok {
		return
	}

	failNode(target_id)
	gossips.gossipMap[target_id] = &shared.Gossip{
		Member: failed_members.memberMap[target_id],
		TTL:    TTL,
	}

	membershipList()
}

func suspectTimeout(id int32, inc_num int32) {
	// wait for suspect_timeout # of rounds then check if node still marked as suspect to fail
	timer := time.NewTimer(time.Duration(SUSPECT_TIMEOUT) * time.Duration(PROTOCAL_PERIOD) * time.Second)
	<-timer.C
	if member, ok := members.memberMap[id]; ok && member.IncNum == inc_num && member.State == shared.NodeState_SUSPECT {
		log.Printf("[INFO] Suspect Timeout for node %d, marking as FAILED", id)
		fmt.Printf("[INFO] Suspect Timeout for node %d, marking as FAILED\n", id)
		failNode(id)
	}
}

func pingTimeoutSuspicion(target_id int32) {
	if _, ok := failed_members.memberMap[target_id]; ok {
		log.Printf("[INFO] Node %d already failed", target_id)
		return
	}

	if members.memberMap[target_id].State == shared.NodeState_SUSPECT {
		pingTimeout(target_id)
		return
	}

	log.Printf("[INFO] ACK not received in time from node %d, marking as SUSPECT", target_id)
	fmt.Printf("[INFO] ACK not received in time from node %d, marking as SUSPECT\n", target_id)
	members.memberMap[target_id].State = shared.NodeState_SUSPECT
	go suspectTimeout(target_id, members.memberMap[target_id].IncNum)

	gossips.gossipMap[target_id] = &shared.Gossip{
		Member: members.memberMap[target_id],
		TTL:    TTL,
	}
	membershipList()
}

func getRandomID() int32 {
	return int32(rand.Uint32() % (2 << M))
}

func requestIntroducer(cur_member *shared.MemberInfo, introducer string) {
	if introducer == "" {
		cur_member.ID = next_id
		cur_member.Hash = getRandomID()
		next_id++
		return
	}

	log.Printf("Sending gRPC call to introducer %s", introducer)
	conn, err := grpc.NewClient(introducer, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("[FATAL] gRPC did not connect: %v", err)
	}
	defer conn.Close()
	client := shared.NewIntroducerClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	response, err := client.RequestMembershipInfo(ctx, cur_member)
	if err != nil {
		log.Fatalf("[FATAL] gRPC call error: %v", err)
	}
	cur_member.ID = response.ID
	cur_member.Hash = response.Hash
	for _, member := range response.MemberList {
		members.memberMap[member.ID] = member
	}
	for _, fail_member := range response.FailList {
		failed_members.memberMap[fail_member.ID] = fail_member
	}
	gossips.gossipMap[cur_member.ID] = &shared.Gossip{
		Member: cur_member,
		TTL:    TTL,
	}
	ENABLE_SUSPECT = response.EnableSus
}

func updateMembership(pingack *shared.PingAck) {
	for ID, new_gossip := range pingack.GossipBuffer {
		cur_gossip, gossip_exists := gossips.gossipMap[ID]
		_, member_exists := members.memberMap[ID]

		if ID == cur_member.ID {
			switch new_gossip.Member.State {
			case shared.NodeState_FAILED:
				log.Fatal("[INFO] False positive detected as failed, exiting...")
			case shared.NodeState_SUSPECT:
				copy_member := &shared.MemberInfo{
					Address: new_gossip.Member.Address,
					ID:      new_gossip.Member.ID,
					State:   shared.NodeState_ALIVE,
					IncNum:  new_gossip.Member.IncNum + 1,
				}
				copy_gossip := &shared.Gossip{
					Member: copy_member,
					TTL:    TTL,
				}
				gossips.gossipMap[ID] = copy_gossip
			}
			continue
		}

		_, failed_member_exists := failed_members.memberMap[ID]
		if failed_member_exists && failed_members.memberMap[ID].State == shared.NodeState_FAILED {
			continue
		}

		if !member_exists && new_gossip.Member.State == shared.NodeState_FAILED {
			continue
		}

		if !member_exists {
			members.memberMap[ID] = new_gossip.Member
			gossips.gossipMap[ID] = new_gossip
			continue
		}

		switch new_gossip.Member.State {
		case shared.NodeState_FAILED:
			failNode(ID)
			gossips.gossipMap[ID] = new_gossip
		case shared.NodeState_SUSPECT:
			if new_gossip.Member.IncNum >= members.memberMap[ID].IncNum {
				if gossip_exists {
					gossips.gossipMap[ID] = &shared.Gossip{
						Member: new_gossip.Member,
						TTL:    max(new_gossip.TTL, cur_gossip.TTL),
					}
				} else {
					gossips.gossipMap[ID] = new_gossip
				}
				members.memberMap[ID] = new_gossip.Member
			}
		case shared.NodeState_ALIVE:
			if new_gossip.Member.IncNum > members.memberMap[ID].IncNum {
				if gossip_exists {
					gossips.gossipMap[ID] = &shared.Gossip{
						Member: new_gossip.Member,
						TTL:    max(new_gossip.TTL, cur_gossip.TTL),
					}
				} else {
					gossips.gossipMap[ID] = new_gossip
				}
				members.memberMap[ID] = new_gossip.Member
			}
		}
	}
}
