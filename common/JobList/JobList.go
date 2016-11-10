package JobList

import (
	"container/list"
	"log"
	"time"

	"../types"
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

func (JB *JobList) EnQ(I *types.Instance) bool {
	JB.Q.PushBack(I)

	if JB.Len() == 1 && JB.IsMonitor {

		JB.NewCh <- true
	}

	return true
}

func (JB *JobList) DeQ() *types.Instance {
	front := JB.Q.Front()

	if front == nil {
		return nil
	}

	if JB.Len() == 1 && JB.IsMonitor {
		JB.EmptyCh <- true
	}

	return JB.Q.Remove(front).Value.(*types.Instance)
}

func (JB *JobList) Len() { return JB.Q.Len() }

//EventHandler This should be started as a goroutine, it takes two function pointers as argumetns
// When the JobList is new and an entry is made it calls NewEvent,
//When the JobList becomes empty then it automatically calls the EmptyEvent after 5 seconds
//This is useful to call Supress and Unsuppress framework messages
func (JB *JObList) EventHandler(NewEvent func() bool, EmptyEvent func() bool, Frequency time.Duration) {

	JB.IsMonitor = true
	timeoutCh := time.After(Frequency)
	var JobEmpty bool

	for {

		select {

		case <-JB.NewCh:
			log.Printf("JOBLIST: Call NewEvent()")
			NewEvent()
			JobEmpty = false

		case <-JB.EmptyCh:
			log.Printf("JOBLIST: Empty")
			JobEmpty = true

		case <-timeoutCh:
			if JobEmpty {
				log.Printf("JOBLIST: Call EmptyEvent()")
				EmptyEvent()
			}
			timeoutCh = time.After(Frequency)
		}
	}
	JB.IsMonitor = false
}
