package types

//TaskUpdate type used to community with Maintainer goroutine
type TaskUpdate struct {
	Name  string
	State string
	Data  []byte
}

//StatsInfo type is used to store docker stats
type StatsInfo struct {
	StatsTime    string  `json:"read"`
	Network      statnet `json:"network"`
	CStats       cstat   `json:"cpu_stats"`
	MStats       mstat   `json:"memory_stats"`
	BlockIOStats bstat   `json:"blockio_stats"`
}

type statnet struct {
	RxBytes   int64 `json:"rx_bytes"`
	RxPackets int64 `json:"rx_packets"`
	RxErrors  int   `json:"rx_errors"`
	RxDropped int   `json:"rx_dropped"`
	TxBytes   int64 `json:"tx_bytes"`
	TxPackets int64 `json:"tx_packets"`
	TxErrors  int   `json:"tx_errors"`
	TxDropped int   `json:"tx_dropped"`
}

type cstat struct {
	CpuUsage       usage      `json:"cpu_usage"`
	SCpuUsage      int64      `json:"system_cpu_usage"`
	ThrottlingData throttling `json:"throttling_data"`
}

type usage struct {
	TotalUsage        int    `json:"total_usage"`
	PerCpuUsage       string `json:"percpu_usage"`
	UsageInKernelMode int    `json:"usage_in_kernel_mode"`
	UsageInUserMode   int    `json:"usage_in_user_mode"`
}

type throttling struct {
	Periods          int `json:"periods"`
	ThrottledPeriods int `json:"throttled_periods"`
	ThrottledTime    int `json:"throttled_time"`
}

type mstat struct {
	Usage    int64 `json:"usage"`
	MaxUsage int64 `json:"max_usage"`
	Stat     st    `json:"stats"`
	FailCnt  int   `json:"failcnt"`
	Limit    int   `json:"limit"`
}

type bstat struct {
	IOServiceBytesRecursive []string `json:"io_service_bytes_recursive"`
	IOServiceRecursive      []string `json:"io_serviced_recursive"`
	IOQueueRecursive        []string `json:"io_queue_recursive"`
	IOServiceTimeRecursive  []string `json:"io_service_time_recursive"`
	IOWaitTimeRecursive     []string `json:"io_wait_time_recursive"`
	IOMergeTimeRecursive    []string `json:"io_merged_recursive"`
	IOTimeRecursive         []string `json:"io_time_recursive"`
	SectorRecursive         []string `json:"sectors_recursive"`
}
type st struct{}
