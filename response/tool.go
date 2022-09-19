package response

type SystemInfo struct {
	Host       []any   `json:"host"`
	CpuInfo    []any   `json:"cpuInfo"`
	CpuPercent float64 `json:"cpuPercent"`
	MemInfo    []any   `json:"memInfo"`
	Disk       []any   `json:"disk"`
}

type SystemInfoRes struct {
	Data SystemInfo
}
