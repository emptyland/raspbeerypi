package model

const (
	kMtabPath = "/etc/mtab"
)

type stateResponse struct {
	Load           float32 `json:"load"`
	CPUPercent     float32 `json:"cpuPercent"`
	CPUTemperature float32 `json:"cpuTemperature"`
}

type memoryResponse struct {
	Total    int `json:"total"`
	Used     int `json:"used"`
	SwapUsed int `json:"swapUsed"`
}

type diskUsageVO struct {
	MountPoint string `json:"mountPoint"`
	FSType     string `json:"fsType"`
	Total      int64  `json:"total"`
	Used       int64  `json:"used"`
}

type diskUsageResponse struct {
	Entries []diskUsageVO `json:"entries"`
}
