package JobList

import (
	"testing"
	"time"

	"github.com/huawei-cloudfederation/mesos-go-stateful/common/types"
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
		return true
	}
	EmptyEvent := func() bool {
		EmptyFlag = true
		return true
	}

	go JB.EventHandler(NewEvent, EmptyEvent, time.Second)
	time.Sleep(time.Microsecond * 100)

	JB.EnQ(&types.Instance{})

	time.Sleep(time.Microsecond * 100)

	if NewFlag == false {
		t.Logf("value NewFlag =%v Should be =true", NewFlag)
		t.Fail()
	}

	JB.DeQ()

	time.Sleep(time.Microsecond * 100)

	if EmptyFlag == true {
		//Should happen only after one second
		t.Logf("Before Frequency value EmptyFlag =%v Should be =false", EmptyFlag)
		t.Fail()
	}

	time.Sleep(time.Second)

	if EmptyFlag == false {
		t.Logf("After Frequency value EmptyFlag =%v Should be =true", EmptyFlag)
		t.Fail()
	}

}
