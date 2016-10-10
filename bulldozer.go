package bulldozer

import (
	"fmt"
)

type WorkerChannel struct {
	Receiver chan interface{}
	Finisher chan int
	ChanTask Task
}

func NewWorkerChannel(goRoutineCount int, task Task) *WorkerChannel {
	return &WorkerChannel{
		Receiver: make(chan interface{}, goRoutineCount),
		Finisher: make(chan int),
		ChanTask: task,
	}
}

// Send the channel to resp to.
// id is there just for making the logs better of.
func worker(myChan *WorkerChannel, freeChan chan *WorkerChannel, respChan chan interface{}, id int) {
	go func() {
		for {
			select {
			case data := <-myChan.Receiver:
				// Processing task
				result := myChan.ChanTask.Run(data)
				respChan <- result
				// Done let me ask for more work.
				freeChan <- myChan
			case <-myChan.Finisher:
				close(myChan.Receiver)
				close(myChan.Finisher)
				return
			}
		}
	}()
}

// Function to initialize the workers.
// Registers the task with the embarresingly parallel run funtion.
// Creates as many go routines as needed to listen to the tasks.
func InitializeWorkers(workerCount int, respChan chan interface{}, task Task) chan *WorkerChannel {
	freeWorkerChan := make(chan *WorkerChannel, workerCount)
	go func() {
		for i := 0; i < workerCount; i++ {
			workerChan := NewWorkerChannel(workerCount, task)
			worker(workerChan, freeWorkerChan, respChan, i)
			// Everyone is free right now. Ask for some work please !
			freeWorkerChan <- workerChan
		}
	}()
	return freeWorkerChan
}

// Scheduler returns the pipe send data on.
// @args - the workers that are free.
// the channel to call to exit the main program.
func Scheduler(freeWorkerChan chan *WorkerChannel, exitChan chan int, resp chan interface{}, workerCount int) (pipe chan interface{}, finish chan int) {
	pipe = make(chan interface{}, workerCount)
	finish = make(chan int)
	go func() {
		for {
			// pickData only if someone is free.
			freeChan := <-freeWorkerChan
			select {
			case data := <-pipe:
				// Assigning the args for work.
				freeChan.Receiver <- data
			case <-finish:
				// Make sure all the workerChan's are done
				freeChan.Finisher <- 1
				fmt.Println("Closed 0")
				for i := 0; i < workerCount-1; i++ {
					freeChan := <-freeWorkerChan
					freeChan.Finisher <- 1
					fmt.Println("Closed", i+1)
				}
				close(pipe)
				close(finish)
				close(resp)
				exitChan <- 1
				return
			}
		}
	}()
	return
}
