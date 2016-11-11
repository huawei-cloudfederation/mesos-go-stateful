package types

type Config struct {
	UserName      string //Supply a username
	FrameworkName string //Supply a frameworkname
	Master        string //MesosMaster's endpoint zk://mesos.master/2181 or 10.11.12.13:5050
	ExecutorPath  string //Executor's Path from where to distribute
	DBType        string //Type of the database etcd/zk/mysql/consul etcd.,
	DBEndPoint    string //Endpoint of the database
	LogFile       string //Name of the logfile
	ArtifactIP    string //The IP to which we should bind to for distributing the executor among the interfaces
	ArtifactPort  string //The port to which we should bind to for distributing the executor
	HTTPPort      string //Defaults to 8080 if otherwise specify explicitly
	WInfo         WI
}

type WI struct {
	CPU     int
	Mem     int
	Disk    int
	Network string
	Image   string // Image should be downloaded
}

//NewDefaultConfig Default Constructor to create a config file
func NewDefaultConfig() *Config {
	var Cfg Config
	Cfg.UserName = "ubuntu"
	Cfg.FrameworkName = "MrRedis"
	Cfg.Master = "127.0.0.1:5050"
	Cfg.ExecutorPath = "./WorkloadExecutor"
	Cfg.DBType = "etcd"
	Cfg.DBEndPoint = "127.0.0.1:2379"
	Cfg.LogFile = "stderr"
	Cfg.ArtifactIP = "127.0.0.1"
	Cfg.ArtifactPort = "5454"
	Cfg.HTTPPort = "5656"
	Cfg.WInfo = WI{CPU: 1, Mem: 1, Disk: 1, Network: "bridge", Image: "redis:3.0-alpine"}
	return &Cfg
}

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
