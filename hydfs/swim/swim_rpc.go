package swim

import (
	"context"
	"cs425/mp3/shared"
	"fmt"
	"net"
	"os"
	"sort"

	"google.golang.org/grpc"
)

type SwimRPCserver struct {
	shared.UnimplementedIntroducerServer
}

func startGRPCServer(host string) {
	lis, err := net.Listen("tcp", host)
	if err != nil {
		swim_log.Fatalf("[ERROR] failed to listen: %v", err)
	}
	s := grpc.NewServer()
	shared.RegisterIntroducerServer(s, &SwimRPCserver{})
	fmt.Printf("[INFO] SWIM gRPC server started at %v\n", lis.Addr())
	if err := s.Serve(lis); err != nil {
		swim_log.Fatalf("[ERROR] failed to serve: %v", err)
	}
}

func (s *SwimRPCserver) RequestMembershipInfo(ctx context.Context, member *shared.MemberInfo) (*shared.MemberContainer, error) {
	swim_log.Printf("Node %d requested to join", next_id)
	Members.Mu.Lock()
	// failed_members.mu.Lock()
	// gossips.mu.Lock()
	// defer gossips.mu.Unlock()
	// defer failed_members.mu.Unlock()
	defer Members.Mu.Unlock()
	// create array of members to send to new member
	memberlist := make([]*shared.MemberInfo, len(Members.MemberMap)+1)
	memberlist[0] = cur_member
	i := 1
	for _, v := range Members.MemberMap {
		memberlist[i] = v
		i++
	}
	fail_list := make([]*shared.MemberInfo, len(failed_members.MemberMap))
	i = 0
	for _, v := range failed_members.MemberMap {
		fail_list[i] = v
		i++
	}
	member.ID = next_id
	member.Hash = getRandomID()
	member_container := &shared.MemberContainer{
		ID:         member.ID,
		Hash:       member.Hash,
		MemberList: memberlist,
		FailList:   fail_list,
	}
	next_id++

	// add new member to current memberlist
	Members.MemberMap[member_container.ID] = member
	gossips.gossipMap[member_container.ID] = &shared.Gossip{
		Member: member,
		TTL:    TTL,
	}
	member_container.EnableSus = ENABLE_SUSPECT
	member_change_chan <- struct{}{}
	membershipList()
	return member_container, nil
}

func (s *SwimRPCserver) SWIMcmd(ctx context.Context, cmd *shared.SWIMIn) (*shared.SWIMOut, error) {
	command := cmd.GetCmd()
	response := shared.SWIMOut{}
	my_ID := (int)(cur_member.ID)
	switch command {
	case "enable_sus":
		ENABLE_SUSPECT = true
		return &shared.SWIMOut{Output: fmt.Sprintf("Node %d Suspect Enabled\n", cur_member.ID)}, nil
	case "disable_sus":
		ENABLE_SUSPECT = false
		return &shared.SWIMOut{Output: fmt.Sprintf("Node %d Suspect Disabled\n", cur_member.ID)}, nil
	case "list_self":
		response.Output = fmt.Sprintf("Member ID is %d", my_ID)
	case "list_mem":
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
		response.Output = fmt.Sprintf("Node %d Membership List: %v\n", cur_member.ID, member_slice)
	case "leave":
		os.Exit(0)
	case "status_sus":
		response.Output = fmt.Sprintf("Suspicion is %t", ENABLE_SUSPECT)
	case "drop_rate":
		DROP_RATE = cmd.GetDropRate()
	}
	return &response, nil
}
