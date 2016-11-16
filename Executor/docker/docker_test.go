//To run this test please fullfill the following system requirements first
//docker deamon should run in background
//copy the sameersbn/postgresql image to local-machine

package docker

import (
	"os"
	"strings"
	"testing"
)

func TestMain(M *testing.M) {

	//Run the tests
	M.Run()

}

//Run with correct inputs
func TestRun(T *testing.T) {
	var dc Dcontainer

	name := "test"
	image := "sameersbn/postgresql"
	mem := int64(1)
	cmd := []string{}
	logfile := "testLog"

	err := dc.Run(name, image, cmd, mem, logfile)

	if err != nil {
		//If its some other error then fail
		T.Fail()
	}
	dc.GetStats()

}

//Run with incorrect image name
func TestRunImagePullFailed(T *testing.T) {
	var dc Dcontainer

	name := "test"
	image := "testimage"
	mem := int64(1)
	cmd := []string{}
	logfile := "testlog"

	dc.Run(name, image, cmd, mem, logfile)

}

//Run without image name
func TestRunImagePullError(T *testing.T) {
	var dc Dcontainer

	name := "test"
	image := ""
	mem := int64(1)
	cmd := []string{}
	logfile := "testlog"

	dc.Run(name, image, cmd, mem, logfile)

}

//Run without  logfile
func TestRunLogFileError(T *testing.T) {
	var dc Dcontainer

	name := "test"
	image := "testimage"
	mem := int64(1)
	cmd := []string{}
	logfile := ""

	err := dc.Run(name, image, cmd, mem, logfile)

	if err != nil && strings.Contains(err.Error(), "Unable to open the logfileopen : no such file or directory") {
		//If its some other error then fail
		T.Fail()
	}

}

//Run without docker command throws create container error
func TestRunCreateContainerError(T *testing.T) {
	var dc Dcontainer

	name := "test"
	image := "sameersbn/postgresql"
	mem := int64(1)
	cmd := []string{}
	logfile := "testLog"

	err := dc.Run(name, image, cmd, mem, logfile)

	if err == nil {
		//If its some other error then fail
		T.Fail()
	}

}

//Run without docker start command throws container start error
func TestRunStartContainerError(T *testing.T) {
	var dc Dcontainer

	name := "test"
	image := "sameersbn/postgresql"
	mem := int64(1)
	cmd := []string{"docker"}
	logfile := "testLog"

	err := dc.Run(name, image, cmd, mem, logfile)

	if err == nil {
		//If its some other error then fail
		T.Fail()
	}

}

//Wait without container id
func TestWaitWithoutID(T *testing.T) {
	var dc Dcontainer

	dc.ID = ""

	val := dc.Wait()

	if val != -1 {
		T.Fail()
	}
}

//close with input
func TestClose(T *testing.T) {
	var dc Dcontainer

	dc.LogFd, _ = os.Create("test2")

	dc.Close(false)

}

//kill without  container id
func TestKillError(T *testing.T) {
	var dc Dcontainer

	dc.ID = ""

	err := dc.Kill()

	if err != nil && !strings.Contains(err.Error(), "Invalid Container") {
		//If its some other error then fail
		T.Fail()
	}

}
