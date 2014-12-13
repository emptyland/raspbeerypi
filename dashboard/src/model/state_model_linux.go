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
	return nil
}
