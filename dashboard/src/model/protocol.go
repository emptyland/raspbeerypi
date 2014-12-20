package model

const (
	kMtabPath    = "/etc/mtab"
	kLoadAvgPath = "/proc/loadavg"
	kCPUStatPath = "/proc/stat"
	kMemInfoPath = "/proc/meminfo"
	kTemp0Path   = "/sys/class/thermal/thermal_zone0/temp"
)

type stateResponse struct {
	LoadAvg        []float32 `json:"loadAvg"`
	CPUPercent     float32   `json:"cpuPercent"`
	CPUTemperature float32   `json:"cpuTemperature"`
}

type memoryResponse struct {
	Total    int64 `json:"total"`
	Used     int64 `json:"used"`
	SwapUsed int64 `json:"swapUsed"`
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

type jobContentVO struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Desc   string `json:"desc"`
	Code   string `json:"code"`
	Lang   string `json:"lang"`
	Cron   string `json:"cron"`
	Enable bool   `json:"enable"`
}

type jobContentResponse struct {
	Entries []jobContentVO `json:"entries"`
}

type jobContentRequest jobContentResponse

type operationResponse struct {
	Ok  bool   `json:"ok"`
	Msg string `json:"msg"`
}
