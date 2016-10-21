
package  TaskMon 

import (
        "fmt"
        "testing"
//        "strings"
//	"os"
	 "../docker"
)

func TestMain(M *testing.M) {

        //Run the tests
        M.Run()

}

func TestLaunchWorkload(T *testing.T) {
        var  tm  TaskMon 
	tm.Image = "hello-world"
	tm.IP = "127.0.0.1"
	tm.Port = 2375
	tm.Container = &docker.Dcontainer{}
	

	err := tm.launchWorkload(false,"127.0.0.1","2375")

	fmt.Println(err)


/*	if err != nil  {
		//If its some other error then fail
		T.Fail()
	}*/


}
