package main

import (
	"fmt"
)

type Task interface {
	Run(interface{}) interface{}
}

type MyTask struct {
}

func (mt *MyTask) Run(data interface{}) interface{} {
	res := fmt.Sprintf("Data is %d", data)
	return res
}

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

func initializeWorkers(workerCount int, respChan chan interface{}, task Task) chan *WorkerChannel {
	freeWorkerChan := make(chan *WorkerChannel, workerCount)
	func() {
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
func scheduler(freeWorkerChan chan *WorkerChannel, exitChan chan int, resp chan interface{}, workerCount int) (pipe chan int, finish chan int) {
	pipe = make(chan int, workerCount)
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
				fmt.Println("Closed i")
				for i := 0; i < workerCount-1; i++ {
					freeChan := <-freeWorkerChan
					freeChan.Finisher <- 1
					fmt.Println("Closed ", i)
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

func main() {
	exitChan := make(chan int)
	workerCount := 100
	mt := &MyTask{}
	respChan := make(chan interface{}, workerCount)
	freeWorkerChan := initializeWorkers(workerCount, respChan, mt)
	input, finish := scheduler(freeWorkerChan, exitChan, respChan, workerCount)

	go func() {
		for i := 0; i < 100000; i++ {
			resp := <-respChan
			fmt.Println(resp)
		}
	}()

	for i := 0; i < 100000; i++ {
		input <- i
	}

	finish <- 1
	<-exitChan
}
