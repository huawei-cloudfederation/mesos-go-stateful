package JobList

import (
	"container/list"

	"github.com/huawei-cloudfederation/mesos-go-stateful/common/logs"
	"time"
)

type JobList struct {
	NewCh     chan bool
	EmptyCh   chan bool
	IsMonitor bool
	Q         *list.List
}

func NewJobList() *JobList {
	var JB JobList

	JB.Q = list.New()
	JB.Q.Init()
	JB.NewCh = make(chan bool)
	JB.EmptyCh = make(chan bool)

	return &JB

}

func (JB *JobList) EnQ(I interface{}) bool {
	JB.Q.PushBack(I)

	if JB.Len() == 1 && JB.IsMonitor {
		JB.NewCh <- true
	}

	return true
}

func (JB *JobList) DeQ() interface{} {
	front := JB.Q.Front()

	if front == nil {
		return nil
	}

	if JB.Len() == 1 && JB.IsMonitor {
		JB.EmptyCh <- true
	}

	I := front.Value.(interface{})

	JB.Q.Remove(front)

	return I
}

func (JB *JobList) Front() *list.Element {
	return JB.Q.Front()
}

func (JB *JobList) Back() *list.Element {
	return JB.Q.Back()
}

func (JB *JobList) Delete(E *list.Element) {

	JB.Q.Remove(E)
}

func (JB *JobList) Len() int { return JB.Q.Len() }

//EventHandler This should be started as a goroutine, it takes two function pointers as argumetns
// When the JobList is new and an entry is made it calls NewEvent,
//When the JobList becomes empty then it automatically calls the EmptyEvent after 5 seconds
//This is useful to call Supress and Unsuppress framework messages
func (JB *JobList) EventHandler(NewEvent func() bool, EmptyEvent func() bool, Frequency time.Duration) {

	JB.IsMonitor = true
	timeoutCh := time.After(Frequency)
	var JobEmpty bool

	for {

		select {

		case <-JB.NewCh:
			logs.Printf("JOBLIST: Call NewEvent()")
			NewEvent()
			JobEmpty = false

		case <-JB.EmptyCh:
			logs.Printf("JOBLIST: Empty")
			JobEmpty = true

		case <-timeoutCh:
			timeoutCh = time.After(Frequency)
			len := JB.Len()
			if len == 0 {
				if JobEmpty {
					logs.Printf("JOBLIST: Call EmptyEvent()")
					EmptyEvent()
				}
			} else {
				logs.Printf("JOBLIST Currently has %d jobs pending", len)
			}
		}
	}
	JB.IsMonitor = false
}
