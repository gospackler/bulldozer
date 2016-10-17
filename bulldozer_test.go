// Create n workers. The run function appends worker id to the input and returns it.
package bulldozer

import (
	"fmt"
	"testing"
)

type MyTask struct {
}

func (mt *MyTask) Run(data interface{}) interface{} {
	res := fmt.Sprintf("Data is %d", data)
	return res
}

func TestNWorkers(t *testing.T) {
	exitChan := make(chan int)
	workerCount := 3
	mt := &MyTask{}
	respChan := make(chan interface{}, workerCount)
	freeWorkerChan := InitializeWorkers(workerCount, respChan, mt)
	input, finish := Scheduler(freeWorkerChan, exitChan, respChan, workerCount)
	done := make(chan int)
	go func() {
		for i := 0; i < 100000; i++ {
			resp := <-respChan
			fmt.Println(resp)
		}
		done <- 1
	}()

	for i := 0; i < 100000; i++ {
		input <- i
	}

	<-done
	finish <- 1
	<-exitChan
}
