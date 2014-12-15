package model

import (
	"api"
	"bufio"
	"log"
	"os"
	"regexp"
	"strconv"
)

type MemoryModel struct {
	mocked memoryResponse
}

var _ = (api.Model)(&MemoryModel{})
var _ = log.Printf

func (self *MemoryModel) Access(appKey string, token string) bool {
	return true
}

const (
	kMemoryTotal = "MemTotal"
	kMemoryFree  = "MemFree"
	kSwapTotal   = "SwapTotal"
	kSwapFree    = "SwapFree"
)

func (self *MemoryModel) GetApiMemory(res *memoryResponse) error {
	memInfo, err := getMemInfo()
	if err != nil {
		return err
	}

	res.Total = memInfo[kMemoryTotal]
	res.Used = res.Total - memInfo[kMemoryFree]
	res.SwapUsed = memInfo[kSwapTotal] - memInfo[kSwapFree]
	return nil
}

var kLineRe = regexp.MustCompile(`^(.+):\s+(\d+) kB$`)

func getMemInfo() (map[string]int64, error) {
	memInfoFile, err := os.Open(kMemInfoPath)
	if err != nil {
		return nil, err
	}
	defer memInfoFile.Close()

	memInfo := make(map[string]int64)

	reader := bufio.NewReader(memInfoFile)
	for {
		var raw []byte
		if raw, _, err = reader.ReadLine(); err != nil {
			break
		}

		matched := kLineRe.FindAllStringSubmatch(string(raw), -1)
		memInfo[matched[0][1]], err = strconv.ParseInt(matched[0][2], 10, 64)
		if err != nil {
			return nil, err
		}
	}

	return memInfo, nil
}
