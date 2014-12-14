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
	if self.mocked.LoadAvg == nil {
		self.mocked.LoadAvg = make([]float32, 3)
	}
	self.mocked.LoadAvg[0] += 0.2
	self.mocked.LoadAvg[1] += 0.2
	self.mocked.LoadAvg[2] += 0.2

	self.mocked.CPUPercent += 0.3
	self.mocked.CPUTemperature += 0.5

	res.LoadAvg = self.mocked.LoadAvg
	res.CPUPercent = self.mocked.CPUPercent
	res.CPUTemperature = self.mocked.CPUTemperature

	return nil
}
