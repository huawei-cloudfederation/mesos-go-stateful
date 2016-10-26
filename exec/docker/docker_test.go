//To run this test please fullfill the following system requirements first
//docker deamon should run in background
//copy the sameersbn/postgresql image to local-machine

package docker

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestMain(M *testing.M) {

	//Run the tests
	M.Run()

}

func TestRun(T *testing.T) {
	var dc Dcontainer

	name := "test"
	image := "sameersbn/postgresql"
	mem := int64(1)
	cmd := []string{}
	logfile := "testLog"

	err := dc.Run(name, image, cmd, mem, logfile)

	fmt.Println(err)

	if err != nil {
		//If its some other error then fail
		T.Fail()
	}
	dc.GetStats()

}

func TestRunImagePullFailed(T *testing.T) {
	var dc Dcontainer

	name := "test"
	image := "testimage"
	mem := int64(1)
	cmd := []string{}
	logfile := "testlog"

	err := dc.Run(name, image, cmd, mem, logfile)

	fmt.Println(err)

}

func TestRunImagePullError(T *testing.T) {
	var dc Dcontainer

	name := "test"
	image := ""
	mem := int64(1)
	cmd := []string{}
	logfile := "testlog"

	err := dc.Run(name, image, cmd, mem, logfile)

	fmt.Println(err)

}

func TestRunLogFileError(T *testing.T) {
	var dc Dcontainer

	name := "test"
	image := "testimage"
	mem := int64(1)
	cmd := []string{}
	logfile := ""

	err := dc.Run(name, image, cmd, mem, logfile)

	fmt.Println(err)

	if err != nil && strings.Contains(err.Error(), "Unable to open the logfileopen : no such file or directory") {
		//If its some other error then fail
		T.Fail()
	}

}

func TestRunCreateContainerError(T *testing.T) {
	var dc Dcontainer

	name := "test"
	image := "sameersbn/postgresql"
	mem := int64(1)
	cmd := []string{}
	logfile := "testLog"

	err := dc.Run(name, image, cmd, mem, logfile)

	fmt.Println(err)

	if err == nil {
		//If its some other error then fail
		T.Fail()
	}

}

func TestRunStartContainerError(T *testing.T) {
	var dc Dcontainer

	name := "test"
	image := "sameersbn/postgresql"
	mem := int64(1)
	cmd := []string{"docker"}
	logfile := "testLog"

	err := dc.Run(name, image, cmd, mem, logfile)

	fmt.Println(err)

	if err == nil {
		//If its some other error then fail
		T.Fail()
	}

}

func TestWaitWithoutID(T *testing.T) {
	var dc Dcontainer

	dc.ID = ""

	val := dc.Wait()

	if val != -1 {
		T.Fail()
	}
}

func TestClose(T *testing.T) {
	var dc Dcontainer

	dc.LogFd, _ = os.Create("test2")

	dc.Close(false)

}

func TestKillError(T *testing.T) {
	var dc Dcontainer

	dc.ID = ""

	err := dc.Kill()

	fmt.Println(err)

	if err != nil && !strings.Contains(err.Error(), "Invalid Container") {
		//If its some other error then fail
		T.Fail()
	}

}
