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

func TestSetUp(T *testing.T) {
	var db etcdDB
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"action":"set","node":{"key":"/test","value":"Hello","modifiedIndex":4,"createdIndex":4}}`)
	}))
	defer ts.Close()

	config := ts.URL


	db.BaseDir = "/test/"
	db.InstDir = "/test/instance"
	db.ConfDir = "/test/config"


	db.Setup(config)
}

func  TestSetUpError(T *testing.T){

	var db etcdDB

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "A-ok")
	}))
	defer ts.Close()

	 config :=   ts.URL 


	err := db.Setup(config)

	fmt.Println(err)

	if err != nil && !strings.Contains(err.Error(), "response is invalid json. The endpoint is probably not valid etcd cluster endpoint") {
		//If its some other error then fail
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

	db.isSetup = true
	db.BaseDir = "/testDir/"
	db.InstDir = "/test/instance"
	db.ConfDir = "/testDir/config"

	err := db.CreateSection("testDir")

	fmt.Println(err)

}

func TestCreateSectionError(T *testing.T) {

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
	db.InstDir = "/test/instance"
	db.ConfDir = "/test/config"

	err := db.CreateSection("testDir")

	fmt.Println(err)

	if err != nil && !strings.Contains(err.Error(), "response is invalid json. The endpoint is probably not valid etcd cluster endpoint") {
		//If its some other error then fail
		T.Fail()
	}
}


func TestSetError(T *testing.T) {

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

	db.isSetup = true
	db.BaseDir = "/test/"
	db.InstDir = "/test/instance"
	db.ConfDir = "/test/config"

	err := db.Set("testDir","Hello")

	fmt.Println(err)

	if err != nil && !strings.Contains(err.Error(), "response is invalid json. The endpoint is probably not valid etcd cluster endpoint") {
		//If its some other error then fail
		T.Fail()
	}
}

func TestGet(T *testing.T) {

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

	db.isSetup = true
	db.BaseDir = "/testDir/"
	db.InstDir = "/test/instance"
	db.ConfDir = "/testDir/config"

	resp, err:= db.Get("testDir")

	fmt.Println(resp,err)

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

	db.isSetup = true
	db.BaseDir = "/testDir/"
	db.InstDir = "/test/instance"
	db.ConfDir = "/testDir/config"

	resp, err:= db.Get("testDir")

	fmt.Println(resp,err)

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

	db.isSetup = true
	db.BaseDir = "/testDir/"
	db.InstDir = "/test/instance"
	db.ConfDir = "/testDir/config"

	err, resp := db.IsDir("testDir")

	fmt.Println(resp,err)

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

	db.isSetup = true
	db.BaseDir = "/testDir/"
	db.InstDir = "/test/instance"
	db.ConfDir = "/testDir/config"



	err, resp := db.IsDir("testDir")

	fmt.Println(resp,err)

}

func TestIsKeyError(T *testing.T) {

	var db etcdDB

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"action":"get","node":{"errorCode":100,"message":"Key not found"}`)
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
	db.BaseDir = "/testDir/"
	db.InstDir = "/test/instance"
	db.ConfDir = "/testDir/config"

	resp,err := db.IsKey("testDir")

	fmt.Println(resp,err)
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

	db.isSetup = true
	db.BaseDir = "/testDir/"
	db.InstDir = "/test/instance"
	db.ConfDir = "/testDir/config"

	resp,err := db.IsKey("testDir")

	fmt.Println(resp,err)

}

func TestIsKeyNotFoundError(T *testing.T) {

	var db etcdDB

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"action":"get","node":{"errorCode":100,"message":"Key found"}`)
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
	db.BaseDir = "/testDir/"
	db.InstDir = "/test/instance"
	db.ConfDir = "/testDir/config"

	 resp,err := db.IsKey("testDir")

	fmt.Println(resp,err)
	if resp != false && !strings.Contains(err.Error(), "response is invalid json. The endpoint is probably not valid etcd cluster endpoint") {
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

	db.isSetup = true
	db.BaseDir = "/test/"
	db.InstDir = "/test/instance"
	db.ConfDir = "/test/config"

	err := db.DeleteSection("testDir")

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

	db.isSetup = true
	db.BaseDir = "/testDir/"
	db.InstDir = "/test/instance"
	db.ConfDir = "/testDir/config"

	err := db.Del("testDir")

	fmt.Println(err)

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

	db.isSetup = true
	db.BaseDir = "/test/"
	db.InstDir = "/test/instance"
	db.ConfDir = "/test/config"

	err := db.Del("testDir")

	fmt.Println(err)

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

	db.isSetup = true
	db.BaseDir = "/testDir/"
	db.InstDir = "/test/instance"
	db.ConfDir = "/testDir/config"

	resp, err:= db.ListSection("testDir",true)

	fmt.Println(resp,err)

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

	db.isSetup = true
	db.BaseDir = "/testDir/"
	db.InstDir = "/test/instance"
	db.ConfDir = "/testDir/config"

	resp, err:= db.ListSection("testDir",true)

	fmt.Println(resp,err)

	if err != nil && !strings.Contains(err.Error(), "response is invalid json. The endpoint is probably not valid etcd cluster endpoint") {
		//If its some other error then fail
		T.Fail()
	}
}

func TestCleanSlateError(T *testing.T) {

	var db etcdDB

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"action":"delete","node":{}}`)
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
	db.InstDir = "/test/instance"
	db.ConfDir = "/test/config"

	err := db.CleanSlate()

	fmt.Println(err)

	if err != nil && !strings.Contains(err.Error(), "response is invalid json. The endpoint is probably not valid etcd cluster endpoint") {
		//If its some other error then fail
		T.Fail()
	}
}
