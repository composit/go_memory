package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func main() {
	for i := 0; ; i++ {
		startedAt := time.Now()
		meminfo, err := ioutil.ReadFile("/proc/meminfo")
		if err != nil {
			log.Fatal(err)
		}

		lines := strings.Split(string(meminfo), "\n")
		memInfo := make(map[string]float64)
		memInfo["MemTotal"] = 0
		memInfo["MemFree"] = 0
		memInfo["Buffers"] = 0
		memInfo["Cached"] = 0
		memInfo["SwapTotal"] = 0
		memInfo["SwapFree"] = 0

		memRegexp := regexp.MustCompile(`^(\w+):\s+(\d+)\s`)

		for _, line := range lines {
			matches := memRegexp.FindStringSubmatch(line)
			if len(matches) == 3 {
				key := matches[1]
				value, err := strconv.ParseFloat(matches[2], 64)
				if err != nil {
					log.Fatal(err)
				}
				memInfo[key] = value
			}
		}

		memTotal := memInfo["MemTotal"] / 1024
		memFree := (memInfo["MemFree"] + memInfo["Buffers"] + memInfo["Cached"]) / 1024
		memUsed := memTotal - memFree
		memTotalFloat := memTotal
		memPercentUsed := memUsed / memTotalFloat * 100
		if err != nil {
			log.Fatal(err)
		}

		swap_total := memInfo["SwapTotal"] / 1024
		swap_free := memInfo["SwapFree"] / 1024
		swap_used := swap_total - swap_free
		swap_percent_used := 0.0
		if swap_total != 0 {
			swap_percent_used = swap_used / swap_total * 100
		}

		//  will be passed at the end to report to Scout
		report_data := make(map[string]float64)

		report_data["size"] = memTotal
		report_data["used"] = memUsed
		report_data["avail"] = memTotal - memUsed
		report_data["used_percent"] = memPercentUsed

		report_data["swap_size"] = swap_total
		report_data["swap_used"] = swap_used
		if swap_total != 0 {
			report_data["swap_used_percent"] = swap_percent_used
		}
		var _ = report_data
		fmt.Println(time.Since(startedAt))

		time.Sleep(1 * time.Second)
	}
}
