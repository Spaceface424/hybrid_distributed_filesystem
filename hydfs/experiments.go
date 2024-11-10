package hydfs

import (
	crypto "crypto/rand"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func loadDataset(num_files int, file_size_kb int) {
	hydfs_base_filename := "experiment_file"
	data := make([]byte, file_size_kb*1000)
	_, err := crypto.Read(data)
	if err != nil {
		hydfs_log.Fatalln("[ERROR] failed to generate random data", err)
	}
	for i := range num_files {
		local_filename := fmt.Sprintf("tmp/experiment/%d_%s.txt", i, hydfs_base_filename)
		hydfs_filename := fmt.Sprintf("%d_%s", i, hydfs_base_filename)
		_ = os.Mkdir("tmp/experiment", 0777)
		file, _ := os.Create(local_filename)
		file.Write(data)
		res, _ := hydfsCreate(local_filename, hydfs_filename)
		if res {
			fmt.Printf("SUCCESS Create %s in HyDFS from %s\n", hydfs_filename, local_filename)
		} else {
			fmt.Println("Create failed")
		}
	}
}

func getExperiment(num_files int, num_gets int, distribution string) {
}

func zipfianSample(upper uint64) uint64 {
	z := rand.NewZipf(rand.New(rand.NewSource(time.Now().UnixNano())), 2, 1, upper)
	return z.Uint64()
}

func uniformSample(upper int) uint64 {
	return uint64(rand.Intn(upper))
}
