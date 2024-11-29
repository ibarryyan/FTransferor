package main

import (
	"encoding/json"
	"log"

	"github.com/shirou/gopsutil/disk"
)

type PartitionStat struct {
	Device     string `json:"device"`
	Mountpoint string `json:"mountpoint"`
	Fstype     string `json:"fstype"`
	Opts       string `json:"opts"`
}

type UsageStat struct {
	Path              string  `json:"path"`
	Fstype            string  `json:"fstype"`
	Total             uint64  `json:"total"`
	Free              uint64  `json:"free"`
	Used              uint64  `json:"used"`
	UsedPercent       float64 `json:"usedPercent"`
	InodesTotal       uint64  `json:"inodesTotal"`
	InodesUsed        uint64  `json:"inodesUsed"`
	InodesFree        uint64  `json:"inodesFree"`
	InodesUsedPercent float64 `json:"inodesUsedPercent"`
}

func FetchDeviceInfo() {
	fn := "FetchDeviceInfo"
	partitionStats := make([]PartitionStat, 0)

	infos, err := disk.Partitions(false)
	if err != nil {
		log.Printf("%s fetching partition info err: %v", fn, err)
		return
	}

	for _, info := range infos {
		var stats PartitionStat
		data, err := json.MarshalIndent(info, "", "  ")
		if err != nil {
			log.Printf("%s marshalling partition err: %v", fn, err)
			continue
		}
		if err := json.Unmarshal(data, &stats); err != nil {
			log.Printf("%s json unmarshal err: %v", fn, err)
			continue
		}
		partitionStats = append(partitionStats, stats)
	}

	log.Printf("Device storage information:")

	for _, stat := range partitionStats {
		var usageStat UsageStat
		info, err := disk.Usage(stat.Device)
		if err != nil {
			log.Printf("%s getting disk usage err: %v", fn, err)
			continue
		}
		data, err := json.MarshalIndent(info, "", "  ")
		if err != nil {
			log.Printf("%s marshalling disk usage err: %v", fn, err)
			continue
		}
		if err := json.Unmarshal(data, &usageStat); err != nil {
			log.Printf("%s json unmarshal err: %v", fn, err)
			continue
		}

		// TODO 存储换算应该是1024
		log.Printf("Partition:%s, Total:%dM, Used:%dM, Over:%dM, Rate:%0.2f%s",
			usageStat.Path, usageStat.Total/1e6, usageStat.Used/1e6, usageStat.Free/1e6, usageStat.UsedPercent, "%")
	}
}
