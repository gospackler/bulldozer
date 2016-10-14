// Create n workers. The run function appends worker id to the input and returns it.
package bulldozer

import (
	"testing"
)

type Worker struct {
	Id         int
	runChan    chan int
	resultChan chan string
}

// While initializing.
// Send back channel to.
// listen for input.
// respond to the input.

// @input is an integer.
// @output is a string corresponding to the integer input.
func (w *Worker) Init() (chan int, chan string) {
	w.runChan = make(chan int)
	w.resultChan = make(chan string)
	go Run()
	return w.runChan, w.resultChan
}

func (w *Worker) Run() {
	data := <-w.runChan
	res := fmt.Sprintf("Worker with Id %d processing input %d", Id, data)
	resultChan <- res
}

func (w *Worker) Finish() {
	close(w.runChan)
	close(w.resultChan)
}

func TestNWorkers(t *testing.T) {
	var workerList []*Worker
	worker1 := &Worker{Id: 9}
	worker2 := &Worker{Id: 19}
	worker3 := &Worker{Id: 29}

	workerList = append(workerList, worker1)
	workerList = append(workerList, worker2)
	workerList = append(workerList, worker3)
	bull := NewBulldozer(workerList)

	for i := 1; i < 100; i++ {
		resp, err := bull.DoWork(i)
		if err != nil {
			t.Error("Could not get anything done", err.Error())
		}
		t.Log("Response for ", i, " ", resp)
	}
}
