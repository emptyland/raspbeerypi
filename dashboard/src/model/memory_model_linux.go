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
    return nil
}
