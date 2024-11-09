package hydfs

import (
	"cs425/mp3/shared"
	"fmt"
)

func ls(hydfs_filename string) {
	ch := make(chan *shared.MemberInfo)
	count := 0
	for cur_member := members.Front(); cur_member != nil; cur_member = cur_member.Next() {
		count += 1
		go sendLsRPC(cur_member.Value.(*shared.MemberInfo), hydfs_filename, ch)
	}

	file_members := make([]*shared.MemberInfo, 0)
	for range count {
		v := <-ch
		if v != nil {
			file_members = append(file_members, v)
		}
	}

	res := fmt.Sprintf("------------[LIST FILE %s, HASH %d]-----------\n", hydfs_filename, hashFilename(hydfs_filename))
	for _, v := range file_members {
		res += fmt.Sprintf("HOSTNAME: %s\tHASH_RING: %d\tID: %d\tSTATE: %s\n", v.Address, v.Hash, v.ID, v.State)
	}
	res += "----------------------------------------------\n"
	fmt.Print(res)
}
