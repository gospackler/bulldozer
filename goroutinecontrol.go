package main

import (
	"fmt"
)

type Task interface {
	Run(int)
}

type MyTask struct {
	Id int
}

func (mt *MyTask) Run(data int) {
	fmt.Println("Data is ", data)
}

type WorkerChannel struct {
	Receiver chan int
	Finisher chan int
	ChanTask Task
}

func NewWorkerChannel(goRoutineCount int, task Task) *WorkerChannel {
	return &WorkerChannel{
		Receiver: make(chan int, goRoutineCount),
		Finisher: make(chan int),
		ChanTask: task,
	}
}

// Send the channel to resp to.
// id is there just for making the logs better of.
func worker(myChan *WorkerChannel, freeChan chan *WorkerChannel, id int) {
	go func() {
		for {
			select {
			case data := <-myChan.Receiver:
				// Processing task
				myChan.ChanTask.Run(data)
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

func initializeWorkers(workerCount int, task Task) chan *WorkerChannel {
	freeWorkerChan := make(chan *WorkerChannel, workerCount)
	func() {
		for i := 0; i < workerCount; i++ {
			workerChan := NewWorkerChannel(workerCount, task)
			worker(workerChan, freeWorkerChan, i)
			// Everyone is free right now. Ask for some work please !
			freeWorkerChan <- workerChan
		}
	}()
	return freeWorkerChan
}

// Scheduler returns the pipe send data on.
// @args - the workers that are free.
func scheduler(freeWorkerChan chan *WorkerChannel, exitChan chan int, workerCount int) (pipe chan int, finish chan int) {
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
	freeWorkerChan := initializeWorkers(workerCount, mt)
	pipe, finish := scheduler(freeWorkerChan, exitChan, workerCount)

	for i := 0; i < 100000; i++ {
		//		time.Sleep(1 * time.Millisecond)
		pipe <- i
	}

	for i := 0; i < 100000; i++ {
		pipe <- i
	}

	finish <- 1
	//	finish <- 1
	<-exitChan
}
