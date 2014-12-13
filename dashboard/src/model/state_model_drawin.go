package model

import (
	"api"
)

type StateModel struct {
	mocked stateResponse
}

var _ = (api.Model)(&StateModel{})

func (self *StateModel) Access(appKey string, token string) bool {
	return true
}

func (self *StateModel) GetApiState(res *stateResponse) error {
	self.mocked.Load += 0.2
	self.mocked.CPUPercent += 0.3
	self.mocked.CPUTemperature += 0.5

	res.Load = self.mocked.Load
	res.CPUPercent = self.mocked.CPUPercent
	res.CPUTemperature = self.mocked.CPUTemperature

	return nil
}
