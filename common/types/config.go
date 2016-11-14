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
	WorkLoad      WLSpec //Definition of basic workload, if this is common to all it can be defined in global config
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
	Cfg.WorkLoad.Default()
	return &Cfg
}
