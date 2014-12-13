package model

import (
	"api"
)

type MemoryModel struct {
	mocked memoryResponse
}

var _ = (api.Model)(&MemoryModel{})

func (self *MemoryModel) Access(appKey string, token string) bool {
	return true
}

func (self *MemoryModel) GetApiMemory(res *memoryResponse) error {
	self.mocked.Used += 100
	self.mocked.Total += 300
	self.mocked.SwapUsed += 500

	res.Used = self.mocked.Used
	res.Total = self.mocked.Total
	res.SwapUsed = self.mocked.SwapUsed

	return nil
}
