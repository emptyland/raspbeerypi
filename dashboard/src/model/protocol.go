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
	User   string `json:"user"`
	Code   string `json:"code"`
	Lang   string `json:"lang"`
	Cron   string `json:"cron"`
	Enable bool   `json:"enable"`
}

type jobContentResponse struct {
	Entries []jobContentVO `json:"entries"`
}

type jobContentRequest jobContentResponse
type jobContentPersistented jobContentResponse

type jobOperationResponse struct {
	Ok     bool     `json:"ok"`
	Code   int      `json:"code"`
	Output []string `json:"output"`
}

type operationResponse struct {
	Ok  bool   `json:"ok"`
	Msg string `json:"msg"`
}

type JobEnvDef struct {
	Env      map[string]string   `json:"env"`
	Pwd      string              `json:"pwd"`
	Metadata string              `json:"metadata"`
	Crontab  string              `json:"crontab"`
	Lang     map[string][]string `json:"lang"`
}

type FileEnvDef struct {
	Home string `json:"home"`
}

type fileStatReponse struct {
	Entries []fileStatVO `json:"entries"`
}

type fileStatVO struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Mode    int16  `json:"mode"`
	Contain int    `json:"contain"`
	Size    int64  `json:"size"`
	ModTime string `json:"modTime"`
	IsDir   bool   `json:"isDir"`
}
