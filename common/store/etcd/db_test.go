package etcd

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"strings"

        "time"

        cli "github.com/coreos/etcd/client"
        "golang.org/x/net/context"

)

func TestMain(M *testing.M) {

	//Run the tests
	M.Run()

}
//var DC_INVALID_ENDPOINT = "http://127.127.127.127:5050"
/*func TestLogin(T *testing.T) {
	var db etcdDB

	db.BaseDir = "/home/divya/mesos-go-statefull"
	db.ConfDir = "/home/divya/mesos-go-statefull" + "/config"
	db.InstDir = "/home/divya/mesos-go-statefull" + "/inst"

	db.Login()
}
*/
func TestSetUp(T *testing.T) {
	var db etcdDB
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"action":"set","node":{"key":"/testDir","value":"Hello","modifiedIndex":4,"createdIndex":4}}`)
	}))
	defer ts.Close()

	config := ts.URL


	db.BaseDir = "/test/"
	db.ConfDir = "/test/config"
	db.InstDir = "/test/InstDir"


	db.Setup(config)
}

func  TestSetUpError(T *testing.T){

	var db etcdDB

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "A-ok")
	}))
	defer ts.Close()

	 config :=   ts.URL 


//	err := db.Setup(config)

	err := db.Setup(config)

	fmt.Println(err)

	if err != nil && !strings.Contains(err.Error(), "response is invalid json. The endpoint is probably not valid etcd cluster endpoint") {
		//If its some other error then fail
		T.Fail()
	}


}


func TestSetUpCreateSection(T *testing.T) {

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

	db.isSetup = true
	db.BaseDir = "/test/"
	db.ConfDir = "/test/config"
	db.InstDir = "/test/InstDir"

	err := db.CreateSection("testDir")

	fmt.Println(err)

	if err == nil {
		//If its some other error then fail
	}
}

func TestSetUpCreateSectionError(T *testing.T) {

	var db etcdDB

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		fmt.Fprintln(w, `{"action":"set","node":{"key":"/message","value":"Hello","modifiedIndex":4,"createdIndex":4}}`)
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

	db.isSetup = true
	db.BaseDir = "/test/"
	db.ConfDir = "/test/config"
	db.InstDir = "/test/InstDir"

	err := db.CreateSection("testDir")

	fmt.Println(err)

	if err != nil && !strings.Contains(err.Error(), "response is invalid json. The endpoint is probably not valid etcd cluster endpoint") {
		//If its some other error then fail
		T.Fail()
	}
}
