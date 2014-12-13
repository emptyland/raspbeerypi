package model

import (
	"api"
)

type DiskUsageModel struct {
}

var _ = (api.Model)(&DiskUsageModel{})

func (self *DiskUsageModel) Access(appKey string, token string) bool {
	return true
}

func (self *DiskUsageModel) GetApiDisk(res *diskUsageResponse) error {
	res.Entries = make([]diskUsageVO, 3)

	res.Entries[0] = diskUsageVO{
		MountPoint: "/",
		FSType:     "ext4",
		Total:      100,
		Used:       30,
	}
	res.Entries[1] = diskUsageVO{
		MountPoint: "/boot",
		FSType:     "vfat",
		Total:      100,
		Used:       60,
	}
	res.Entries[2] = diskUsageVO{
		MountPoint: "/tmp",
		FSType:     "tmpfs",
		Total:      100,
		Used:       90,
	}

	return nil
}
