package swim

import (
	"cs425/mp3/shared"
	"fmt"
	"io"
	"log"
	"math/rand/v2"
	"net"
	"os"
	"sort"
	"sync"
)

type MemberContainer struct {
	MemberMap map[int32]*shared.MemberInfo
	Mu        sync.Mutex
}

type gossipContainer struct {
	gossipMap map[int32]*shared.Gossip
	mu        sync.Mutex
}

type pingackRequest struct {
	conn         *net.UDPConn
	target_addr  *net.UDPAddr
	sender_id    int32
	target_id    int32
	req_id       int32
	round        int32
	inc_num      int32
	message_type shared.MessageType
}

func setupServer(host string, verbose_flag bool) *net.UDPConn {
	cur_member = &shared.MemberInfo{
		Address: host,
		State:   shared.NodeState_ALIVE,
		IncNum:  0,
	}

	Members = MemberContainer{
		MemberMap: make(map[int32]*shared.MemberInfo),
	}
	failed_members = MemberContainer{
		MemberMap: make(map[int32]*shared.MemberInfo),
	}
	gossips = gossipContainer{
		gossipMap: make(map[int32]*shared.Gossip),
	}
	ack_chan = make(chan *shared.PingAck, 10)
	member_change_chan = make(chan struct{}, 10)
	verbose = verbose_flag
	udpAddr, err := net.ResolveUDPAddr("udp", host)
	if err != nil {
		log.Panic("[ERROR] UDP err:", err)
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Panic("[ERROR] Listen UDP err:", err)
	}
	fmt.Println("UDP server started at", conn.LocalAddr())
	return conn
}

func setupLogging() {
	logFile, err := os.OpenFile(fmt.Sprintf("logs/node.%d.log", cur_member.ID), os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	if err != nil {
		log.Fatalf("[ERROR] Log file error: %v", err)
	}
	defer logFile.Close()
	log.SetPrefix(fmt.Sprintf("Node %d, ", cur_member.ID))
	log.SetFlags(log.Lmsgprefix | log.Ltime | log.Lmicroseconds)
	if verbose {
		mw := io.MultiWriter(os.Stdout, logFile)
		log.SetOutput(mw)
	} else {
		log.SetOutput(logFile)
	}
}

func logMembershipList() {
	member_slice := make([]int32, len(Members.MemberMap)+1)
	member_slice[0] = cur_member.ID
	i := 1
	for k := range Members.MemberMap {
		member_slice[i] = k
		i++
	}
	sort.Slice(member_slice, func(i, j int) bool {
		return member_slice[i] < member_slice[j]
	})
	log.Printf("Membership List: %v", member_slice)
}

func fmtGossip(gossip_buffer *map[int32]*shared.Gossip) string {
	gossip_slice := make([]string, len(*gossip_buffer))
	i := 0
	for id, gossip := range *gossip_buffer {
		gossip_slice[i] = fmt.Sprintf("%d: %s ttl: %d,", id, gossip.Member.State, gossip.TTL)
		i++
	}
	return fmt.Sprint(gossip_slice)
}

func PrintMembershipList() {
	member_slice := make([]*shared.MemberInfo, len(Members.MemberMap)+1)
	member_slice[0] = cur_member
	i := 1
	for _, v := range Members.MemberMap {
		member_slice[i] = v
		i++
	}
	sort.Slice(member_slice, func(i, j int) bool {
		return member_slice[i].ID < member_slice[j].ID
	})

	res := "\n-------------------------------------[MEMBERSHIP LIST]-------------------------------------\n\n"
	for _, v := range member_slice {
		var self string
		if v.ID == cur_member.ID {
			self = "(self)"
		}
		res += fmt.Sprintf("HOSTNAME: %s\tHASH_RING: %d\tID: %d\tSTATE: %s %s\n", v.Address, v.Hash, v.ID, v.State, self)
	}
	res += "\n-------------------------------------------------------------------------------------------\n\n"
	fmt.Print(res)
}

func membershipList() {
	if verbose {
		logMembershipList()
	} else {
		PrintMembershipList()
	}
}

func getRoundRobinTarget() int32 {
	Members.Mu.Lock()
	defer Members.Mu.Unlock()

	for rr_index < len(round_robin) {
		_, ok := Members.MemberMap[round_robin[rr_index]]
		rr_index++
		if !ok {
			continue
		}
		return round_robin[rr_index-1]
	}
	round_robin = shuffleRoundRobin()
	rr_index = 1
	return round_robin[rr_index-1]
}

func shuffleRoundRobin() []int32 {
	new_round_robin := make([]int32, len(Members.MemberMap))
	i := 0
	for id := range Members.MemberMap {
		new_round_robin[i] = id
		i++
	}

	rand.Shuffle(len(new_round_robin), func(i, j int) {
		new_round_robin[i], new_round_robin[j] = new_round_robin[j], new_round_robin[i]
	})
	return new_round_robin
}

func decTTL() {
	for id, gossip := range gossips.gossipMap {
		gossip.TTL--
		if gossip.TTL == 0 {
			delete(gossips.gossipMap, id)
		}
	}
}

func failNode(id int32) {
	failed_members.MemberMap[id] = Members.MemberMap[id]
	failed_members.MemberMap[id].State = shared.NodeState_FAILED

	Members.Mu.Lock()
	delete(Members.MemberMap, id)
	Members.Mu.Unlock()
	member_change_chan <- struct{}{}
}
