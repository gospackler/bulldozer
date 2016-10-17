# bulldozer
Bulldozer does loads of work on the background. It can be used for working on embarassingly parallel problems. 

```
type Task interface {
	Run(interface{}) interface{}
}
```

It can run any `task type` that follows the interface definition mentioned above.

## Design
### Task
As mentioned above.

### Bulldozer 
* Initialize workers
	>> Goroutines that are ready to work are initialized.
	>> Initialized with a worker which has the embarassingly parallel function.
* Scheduler
	>> The scheduler handles the initialized go-routines and divides work among them.
	>> The exit signal is received to signal the end of execution.
* respChan
	>> This is the channel where the response from the work comes out.
* Finish signal.
	>> This can be used to finish the execution of the workers.

### Test
