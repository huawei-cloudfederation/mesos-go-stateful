package etcd

import (
	"strings"
	"time"

	cli "github.com/coreos/etcd/client"
	"golang.org/x/net/context"
)

//gloabal variable for etcd
var ETCD_BASEDIR, ETCD_INSTDIR, ETCD_CONFDIR string

type etcdDB struct {
	C       cli.Client      //The client context
	Kapi    cli.KeysAPI     //The api context for Get/Set/Delete/Update/Watcher etc.,
	Ctx     context.Context //Context for the connection mostly set to context.Background
	Cfg     cli.Config      //Configuration details of the connection should be loaded from a configuration file
	isSetup bool            //Has this been setup
	BaseDir string          //base dir
	InstDir string          //instance dir
	ConfDir string          //config dir
}

//New Function to create an etcDB object
func New() *etcdDB {
	return &etcdDB{isSetup: false}
}

//Login This implements connecting to the ETCD instance
func (db *etcdDB) Login() error {

	var err error
	db.C, err = cli.New(db.Cfg)
	if err != nil {

		return err
	}
	db.Kapi = cli.NewKeysAPI(db.C)
	db.Ctx = context.Background()

	return nil
}

// Setup will create/establish connection with the etcd store and also setup
// the nessary environment if etcd is running for the first time
// framework will look for the following location in the central store
// /framework/instances...... -> Will have the entries of all the instances
// /framework/Config/....		-> Will have the entries of all the config information
func (db *etcdDB) Setup(config string) error {
	var err error
	db.Cfg = cli.Config{
		Endpoints: []string{config},
		Transport: cli.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second * 10,
	}

	err = db.Login()
	if err != nil {
		return err
	}

	db.BaseDir = "/Workload"

	ETCD_BASEDIR = db.BaseDir
	err = db.CreateSection(db.BaseDir)
	if err != nil && strings.Contains(err.Error(), "Key already exists") != true {
		return err
	}

	ETCD_INSTDIR = ETCD_BASEDIR + "/Instances"
	db.InstDir = ETCD_INSTDIR
	err = db.CreateSection(db.InstDir)
	if err != nil && strings.Contains(err.Error(), "Key already exists") != true {
		return err
	}

	ETCD_CONFDIR = ETCD_BASEDIR + "/Config"
	db.ConfDir = ETCD_CONFDIR
	err = db.CreateSection(db.ConfDir)
	if err != nil && strings.Contains(err.Error(), "Key already exists") != true {
		return err
	}

	db.isSetup = true
	return nil
}

func (db *etcdDB) IsSetup() bool {
	return db.isSetup
}

func (db *etcdDB) Set(Key string, Value string) error {

	_, err := db.Kapi.Set(db.Ctx, Key, string(Value), nil)
	return err

}

func (db *etcdDB) Get(Key string) (string, error) {

	resp, err := db.Kapi.Get(db.Ctx, Key, nil)
	if err != nil {
		return "", err
	}
	return resp.Node.Value, nil
}

func (db *etcdDB) IsDir(Key string) (error, bool) {
	resp, err := db.Kapi.Get(db.Ctx, Key, nil)

	if err != nil {
		return err, false
	}
	return nil, resp.Node.Dir
}

func (db *etcdDB) IsKey(Key string) (bool, error) {
	_, err := db.Kapi.Get(db.Ctx, Key, nil)

	if err != nil {
		if strings.Contains(err.Error(), "Key not found") != true {
			return false, err
		}
		return false, err
	}
	return true, nil
}

func (db *etcdDB) Update(Key string, Value string, Lock bool) error {

	return nil
}

func (db *etcdDB) Del(Key string) error {

	_, err := db.Kapi.Delete(db.Ctx, Key, nil)

	if err != nil {
		return err
	}
	return nil

}

//CreateSection will create a directory in etcd store
func (db *etcdDB) CreateSection(Key string) error {

	_, err := db.Kapi.Set(db.Ctx, Key, "", &cli.SetOptions{Dir: true, PrevExist: cli.PrevNoExist})

	if err != nil {
		return err
	}

	return nil
}

//DeleteSection section will delete a directory optionally delete
func (db *etcdDB) DeleteSection(Key string) error {

	_, err := db.Kapi.Delete(db.Ctx, Key, &cli.DeleteOptions{Dir: true})
	return err
}

//ListSection will list a directory
func (db *etcdDB) ListSection(Key string, Recursive bool) ([]string, error) {

	resp, err := db.Kapi.Get(db.Ctx, Key, &cli.GetOptions{Sort: true})

	if err != nil {
		return nil, err
	}

	retStr := make([]string, len(resp.Node.Nodes))

	for i, n := range resp.Node.Nodes {
		retStr[i] = n.Key
	}

	return retStr, nil
}

func (db *etcdDB) CleanSlate() error {

	_, err := db.Kapi.Delete(db.Ctx, db.BaseDir, &cli.DeleteOptions{Dir: true, Recursive: true})

	return err
}
