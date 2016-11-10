package JobList

import (
	"fmt"
	"testing"

	"../types"
)

var JB *JobList

func TestMain(M *testing.M) {
	JB = NewJobList()

	M.Run()
}

func TestEnQ(t *testing.T) {

	JB.EnQ(&types.Instance{})

	if JB.Len() != 1 {
		t.Fail()
	}
}

func TestDeQ(t *testing.T) {

	I := JB.DeQ()

	if I == nil {
		t.Fail()
	}
}

func TestMonitor(t *testing.T) {
	var NewFlag, EmptyFlag bool

	NewEvent := func() bool {
		NewFlag = true
	}
	EmptyEvent := func() bool {
		EmptyFlag = true
	}

	go JB.EventHandler(NewEvent, EmptyEvent, time.Second)

	JB.EnQ(&types.Instance{})

	if NewFlag == false {
		t.Fail()
	}

	I := JB.DeQ()

	if EmptyFlag == true {
		//Should happen only after one second
		t.Fail()
	}

	time.Sleep(time.Second)

	if EmptyFlag == false {
		t.Fail()
	}

}
