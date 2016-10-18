package etcd

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"time"

	cli "github.com/coreos/etcd/client"
	"golang.org/x/net/context"
)

func TestMain(M *testing.M) {

	//Run the tests
	M.Run()

}

func TestNew(T *testing.T) {
	var db etcdDB
	New()

	if db.isSetup != false {
		T.Fail()
	}
}

// Login with endpoint
func TestLoginWithEndPoint(T *testing.T) {
	var db etcdDB

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "")
	}))

	defer ts.Close()

	config := ts.URL

	db.Cfg = cli.Config{
		Endpoints: []string{config},
	}

	err := db.Login()
	if err != nil {
		T.FailNow()
	}
}

// Login without endpoint
func TestLogWithoutEndPoint(T *testing.T) {
	var db etcdDB

	db.Cfg = cli.Config{
		Endpoints: []string{},
	}

	err := db.Login()

	if err == nil {
		//Error cannot be nil
		T.Fail()
	}

}

func TestSetUpWithConfig(T *testing.T) {
	var db etcdDB
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"action":"set","node":{"key":"/test","value":"Hello","modifiedIndex":4,"createdIndex":4}}`)
	}))
	defer ts.Close()

	config := ts.URL

	db.BaseDir = "/test"
	db.InstDir = "/test/instance"
	db.ConfDir = "/test/config"

	db.Setup(config)
}

func TestSetUpWithoutConfig(T *testing.T) {

	var db etcdDB

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "A-ok")
	}))
	defer ts.Close()

	config := ts.URL

	err := db.Setup(config)

	fmt.Println(err)

	if err != nil && !strings.Contains(err.Error(), "response is invalid json. The endpoint is probably not valid etcd cluster endpoint") {
		//If its some other error then fail
		T.Fail()
	}

}


func TestIsSetup(T *testing.T) {

	var db etcdDB

	resp := db.IsSetup()

	if resp != false {
		T.Fail()
	}
}

func TestCreateSection(T *testing.T) {

	var db etcdDB

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"action":"set","node":{"key":"/testDir","value":"Hello","modifiedIndex":4,"createdIndex":4}}`)
	}))
	defer ts.Close()

	config := ts.URL

	db.Cfg = cli.Config{
		Endpoints: []string{config},
		Transport: cli.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}

	db.C, _ = cli.New(db.Cfg)

	db.Kapi = cli.NewKeysAPI(db.C)
	db.Ctx = context.Background()

	err := db.CreateSection("/testDir")

	if err != nil {
		T.Fail()
	}
}

func TestCreateSectionError(T *testing.T) {

	var db etcdDB

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "A-ok")
	}))
	defer ts.Close()

	config := ts.URL

	db.Cfg = cli.Config{
		Endpoints: []string{config},
		Transport: cli.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}

	db.C, _ = cli.New(db.Cfg)

	db.Kapi = cli.NewKeysAPI(db.C)
	db.Ctx = context.Background()

	err := db.CreateSection("/testDir")

	fmt.Println(err)

	if err != nil && !strings.Contains(err.Error(), "response is invalid json. The endpoint is probably not valid etcd cluster endpoint") {
		//If its some other error then fail
		T.Fail()
	}
}

func TestSet(T *testing.T) {

	var db etcdDB

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"action":"set","node":{"key":"/testDir","value":"Hello","modifiedIndex":4,"createdIndex":4}}`)
	}))
	defer ts.Close()

	config := ts.URL

	db.Cfg = cli.Config{
		Endpoints: []string{config},
		Transport: cli.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}

	db.C, _ = cli.New(db.Cfg)

	db.Kapi = cli.NewKeysAPI(db.C)
	db.Ctx = context.Background()

	err := db.Set("/testDir", "Hello")
	
	if err != nil{
		T.Fail()
	}

}

func TestSetError(T *testing.T) {

	var db etcdDB

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "A-OK")
	}))
	defer ts.Close()

	config := ts.URL

	db.Cfg = cli.Config{
		Endpoints: []string{config},
		Transport: cli.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}

	db.C, _ = cli.New(db.Cfg)

	db.Kapi = cli.NewKeysAPI(db.C)
	db.Ctx = context.Background()

	err := db.Set("/testDir", "Hello")

	fmt.Println(err)
	if err != nil && !strings.Contains(err.Error(), "response is invalid json. The endpoint is probably not valid etcd cluster endpoint") {
		//If its some other error then fail
		T.Fail()
	}
}

func TestGetValidConfig(T *testing.T) {

	var db etcdDB

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"action":"get","node":{"key":"/testDir","value":"Hello","modifiedIndex":4,"createdIndex":4}}`)
	}))
	defer ts.Close()

	config := ts.URL

	db.Cfg = cli.Config{
		Endpoints: []string{config},
		Transport: cli.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}

	db.C, _ = cli.New(db.Cfg)

	db.Kapi = cli.NewKeysAPI(db.C)
	db.Ctx = context.Background()

	resp, err := db.Get("/testDir")

	fmt.Println(resp, err)
}

func TestGetError(T *testing.T) {

	var db etcdDB

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "A-ok")
	}))
	defer ts.Close()

	config := ts.URL

	db.Cfg = cli.Config{
		Endpoints: []string{config},
		Transport: cli.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}

	db.C, _ = cli.New(db.Cfg)

	db.Kapi = cli.NewKeysAPI(db.C)
	db.Ctx = context.Background()

	resp, err := db.Get("/testDir")

	if resp != "" && !strings.Contains(err.Error(), "response is invalid json. The endpoint is probably not valid etcd cluster endpoint") {
		//If its some other error then fail
		T.Fail()
	}
}

func TestIsDir(T *testing.T) {

	var db etcdDB

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"action":"get","node":{"key":"/testDir","value":"Hello","modifiedIndex":4,"createdIndex":4}}`)
	}))
	defer ts.Close()

	config := ts.URL

	db.Cfg = cli.Config{
		Endpoints: []string{config},
		Transport: cli.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}

	db.C, _ = cli.New(db.Cfg)

	db.Kapi = cli.NewKeysAPI(db.C)
	db.Ctx = context.Background()

	db.IsDir("/testDir")

}

func TestIsDirErro(T *testing.T) {

	var db etcdDB

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "A-ok")
	}))
	defer ts.Close()

	config := ts.URL

	db.Cfg = cli.Config{
		Endpoints: []string{config},
		Transport: cli.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}

	db.C, _ = cli.New(db.Cfg)

	db.Kapi = cli.NewKeysAPI(db.C)
	db.Ctx = context.Background()

	err, resp := db.IsDir("/testDir")

	if resp != false && !strings.Contains(err.Error(), "response is invalid json. The endpoint is probably not valid etcd cluster endpoint") {
		//If its some other error then fail
		T.Fail()
	}

}

func TestIsKey(T *testing.T) {

	var db etcdDB

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"action":"get","node":{"key":"/testDir","value":"Hello","modifiedIndex":4,"createdIndex":4}}`)
	}))
	defer ts.Close()

	config := ts.URL

	db.Cfg = cli.Config{
		Endpoints: []string{config},
		Transport: cli.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}

	db.C, _ = cli.New(db.Cfg)

	db.Kapi = cli.NewKeysAPI(db.C)
	db.Ctx = context.Background()

	 db.IsKey("/testDir")


}

func TestIsKeyError(T *testing.T) {

	var db etcdDB

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "A-ok")
	}))
	defer ts.Close()

	config := ts.URL

	db.Cfg = cli.Config{
		Endpoints: []string{config},
		Transport: cli.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}

	db.C, _ = cli.New(db.Cfg)

	db.Kapi = cli.NewKeysAPI(db.C)
	db.Ctx = context.Background()

	resp, err := db.IsKey("/testDir")

	if resp != false && !strings.Contains(err.Error(), "response is invalid json. The endpoint is probably not valid etcd cluster endpoint") {
		//If its some other error then fail
		T.Fail()
	}

}

func TestUpdate(T *testing.T) {
	var db etcdDB
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"action":"set","node":{"key":"/testDir","value":"Hello","modifiedIndex":4,"createdIndex":4}}`)
	}))

	defer ts.Close()

	config := ts.URL
	db.Cfg = cli.Config{
		Endpoints: []string{config},
		Transport: cli.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}

	db.C, _ = cli.New(db.Cfg)

	db.Kapi = cli.NewKeysAPI(db.C)
	db.Ctx = context.Background()

	err := db.Update("/testDir", "Hello", true)

	if err != nil {
		T.Fail()
	}

}

func TestUpdateError(T *testing.T) {
	var db etcdDB
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "A-ok")
	}))

	defer ts.Close()
	config := ts.URL

	db.Cfg = cli.Config{
		Endpoints: []string{config},
		Transport: cli.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}

	db.C, _ = cli.New(db.Cfg)

	db.Kapi = cli.NewKeysAPI(db.C)
	db.Ctx = context.Background()

	err := db.Update("/testDir", "Hello", true)

	if err != nil && !strings.Contains(err.Error(), "response is invalid json. The endpoint is probably not valid etcd cluster endpoint") {
		//If its some other error then fail
		T.Fail()
	}
}

func TestDeleteSectionError(T *testing.T) {

	var db etcdDB

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "A-ok")
	}))
	defer ts.Close()

	config := ts.URL

	db.Cfg = cli.Config{
		Endpoints: []string{config},
		Transport: cli.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}

	db.C, _ = cli.New(db.Cfg)

	db.Kapi = cli.NewKeysAPI(db.C)
	db.Ctx = context.Background()

	err := db.DeleteSection("/testDir")

	fmt.Println(err)

	if err != nil && !strings.Contains(err.Error(), "response is invalid json. The endpoint is probably not valid etcd cluster endpoint") {
		//If its some other error then fail
		T.Fail()
	}
}

func TestDel(T *testing.T) {

	var db etcdDB

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"action":"delete","node":{"key":"/testDir","value":"Hello","modifiedIndex":4,"createdIndex":4}}`)
	}))
	defer ts.Close()

	config := ts.URL

	db.Cfg = cli.Config{
		Endpoints: []string{config},
		Transport: cli.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}

	db.C, _ = cli.New(db.Cfg)

	db.Kapi = cli.NewKeysAPI(db.C)
	db.Ctx = context.Background()

	db.Del("/testDir")

}

func TestDelError(T *testing.T) {

	var db etcdDB

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "A-ok")
	}))
	defer ts.Close()

	config := ts.URL

	db.Cfg = cli.Config{
		Endpoints: []string{config},
		Transport: cli.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}

	db.C, _ = cli.New(db.Cfg)

	db.Kapi = cli.NewKeysAPI(db.C)
	db.Ctx = context.Background()

	err := db.Del("/testDir")


	if err != nil && !strings.Contains(err.Error(), "response is invalid json. The endpoint is probably not valid etcd cluster endpoint") {
		//If its some other error then fail
		T.Fail()
	}
}

func TestListSection(T *testing.T) {

	var db etcdDB

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"action":"get","node":{"key":"/testDir","value":"Hello","modifiedIndex":4,"createdIndex":4}}`)
	}))
	defer ts.Close()

	config := ts.URL

	db.Cfg = cli.Config{
		Endpoints: []string{config},
		Transport: cli.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}

	db.C, _ = cli.New(db.Cfg)

	db.Kapi = cli.NewKeysAPI(db.C)
	db.Ctx = context.Background()

	db.ListSection("/testDir", true)

}

func TestListSectionError(T *testing.T) {

	var db etcdDB

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "A-ok")
	}))
	defer ts.Close()

	config := ts.URL

	db.Cfg = cli.Config{
		Endpoints: []string{config},
		Transport: cli.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}

	db.C, _ = cli.New(db.Cfg)

	db.Kapi = cli.NewKeysAPI(db.C)
	db.Ctx = context.Background()

	_, err := db.ListSection("/testDir", true)

	if err != nil && !strings.Contains(err.Error(), "response is invalid json. The endpoint is probably not valid etcd cluster endpoint") {
		//If its some other error then fail
		T.Fail()
	}
}

func TestCleanSlateError(T *testing.T) {

	var db etcdDB

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "A-ok")
	}))
	defer ts.Close()

	config := ts.URL

	db.Cfg = cli.Config{
		Endpoints: []string{config},
		Transport: cli.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}

	db.C, _ = cli.New(db.Cfg)

	db.Kapi = cli.NewKeysAPI(db.C)
	db.Ctx = context.Background()

	db.BaseDir = "/testDir"

	err := db.CleanSlate()

	if err != nil && !strings.Contains(err.Error(), "response is invalid json. The endpoint is probably not valid etcd cluster endpoint") {
		//If its some other error then fail
		T.Fail()
	}
}
