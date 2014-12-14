package model

import (
	"api"
	"fmt"
	"log"
	"os"
	"time"
)

type StateModel struct {
	mocked stateResponse
}

var _ = (api.Model)(&StateModel{})

func (self *StateModel) Access(appKey string, token string) bool {
	return true
}

func (self *StateModel) GetApiState(res *stateResponse) error {
	var err error

	if res.LoadAvg, err = getLoadAvg(); err != nil {
		return err
	}

	var temp int
	if temp, err = getCPUTemp(0); err != nil {
		return err
	}
	res.CPUTemperature = float32(temp) / 1000.0

	var time1, time2 int64
	if time1, time2, err = getCPUUsage(); err != nil {
		return err
	}
	res.CPUPercent = (1.0 - float32(time1)/float32(time2)) * 100.0

	return nil
}

func getLoadAvg() ([]float32, error) {
	loadAvgFile, err := os.Open(kLoadAvgPath)
	if err != nil {
		log.Printf("can not open %s: %v", kLoadAvgPath, err)
		return nil, err
	}
	defer loadAvgFile.Close()

	loadAvg := make([]float32, 3)
	if _, err = fmt.Fscanf(loadAvgFile, "%f %f %f", &loadAvg[0], &loadAvg[1], &loadAvg[2]); err != nil {
		return nil, err
	}

	return loadAvg, nil
}

func getCPUTemp(cpuId int) (int, error) {
	cpuTempFile, err := os.Open(kTemp0Path)
	if err != nil {
		log.Printf("can not open %s: %v", kTemp0Path, err)
		return 0, err
	}
	defer cpuTempFile.Close()

	cpuTemp := 0
	if _, err = fmt.Fscanf(cpuTempFile, "%d", &cpuTemp); err != nil {
		return 0, err
	}

	return cpuTemp, nil
}

const (
	kUser = iota
	kNice
	kSystem
	kIdle
	kIOWait
	kIRQ
	kSoftIRQ
	kMaxTime
)

func getCPUUsage() (int64, int64, error) {
	var begin []int64
	var err error

	if begin, err = getCPUTime(); err != nil {
		return 0, 0, err
	}

	time.Sleep(100 * time.Millisecond)

	var end []int64
	if end, err = getCPUTime(); err != nil {
		return 0, 0, err
	}

	var total [2]int64
	for _, time := range begin {
		total[0] += time
	}
	for _, time := range end {
		total[1] += time
	}

	time1 := end[kIdle] - begin[kIdle]
	time2 := total[1] - total[0]

	return time1, time2, nil
}

func getCPUTime() ([]int64, error) {
	cpuStatFile, err := os.Open(kCPUStatPath)
	if err != nil {
		log.Printf("can not open %s: %v", kCPUStatPath, err)
		return nil, err
	}
	defer cpuStatFile.Close()

	times := make([]int64, kMaxTime)
	if _, err = fmt.Fscanf(cpuStatFile, "cpu  %d %d %d %d %d %d %d", &times[0], &times[1], &times[2], &times[3], &times[4], &times[5], &times[6]); err != nil {
		return nil, err
	}

	return times, nil
}
