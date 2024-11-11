package hydfs

import (
	crypto "crypto/rand"
	"fmt"
	"math"
	"math/rand"
	"os"
	"sort"
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

func getExperiment(num_files int, num_gets int, distribution func(int) int, percent_append int) {
	times := make([]time.Duration, num_gets)
	local_filename_append := "tmp/business_3.txt"
	hydfs_base_filename := "experiment_file"
	for i := range num_gets {
		target_file_idx := distribution(num_files)
		os.MkdirAll("tmp/experiment_get", 0777)
		local_filename := fmt.Sprintf("tmp/experiment_get/%d_%s.txt", target_file_idx, hydfs_base_filename)
		hydfs_filename := fmt.Sprintf("%d_%s", target_file_idx, hydfs_base_filename)
		start := time.Now()
		if rand.Intn(100) < percent_append {
			res, _ := hydfsAppend(local_filename_append, hydfs_filename)
			if !res {
				fmt.Println("APPEND failed")
			}
		} else {
			res, _ := hydfsGet(hydfs_filename, local_filename)
			if !res {
				fmt.Println("GET failed")
			}
		}
		times[i] = time.Since(start)
	}

	var sum time.Duration
	for _, d := range times {
		sum += d
	}
	avg := sum / time.Duration(num_gets)

	var squaredDiffs float64
	avgFloat := float64(avg)
	for _, d := range times {
		diff := float64(d) - avgFloat
		squaredDiffs += diff * diff
	}
	variance := squaredDiffs / float64(num_gets)
	stdDev := time.Duration(math.Sqrt(variance))

	sort.Slice(times, func(i, j int) bool {
		return times[i] < times[j]
	})

	fmt.Printf("AVERAGE: %v\n", avg)
	fmt.Printf("MEDIAN: %v\n", times[len(times)/2])
	fmt.Printf("STD: %v\n", stdDev)
	fmt.Printf("MAX: %v\n", times[len(times)-1])
	fmt.Printf("MIN: %v\n", times[0])
}

func zipfianSample(upper int) int {
	z := rand.NewZipf(rand.New(rand.NewSource(time.Now().UnixNano())), 2, 1, uint64(upper))
	return int(z.Uint64())
}

func uniformSample(upper int) int {
	return int(rand.Intn(upper))
}
