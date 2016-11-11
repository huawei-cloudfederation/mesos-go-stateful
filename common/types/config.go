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
	WorkLoad      WL     //Definition of basic workload, if this is common to all it can be defined in global config
}

type WL struct {
	CPU     float64 //Number of CPU each workload will require NOTE:an instance may have more than one workload
	Mem     float64 //Memory requirement of each workload  in MB
	Disk    float64 //Disk requirement of each workload in GB
	Network string  //Most preffered networkign layer HOST or DOckerBridge
	Image   string  // Docker Image of this workload
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
	Cfg.WorkLoad = WL{CPU: 1.0, Mem: 100.0, Disk: 1.0, Network: "host", Image: "redis:3.0-alpine"}
	return &Cfg
}
