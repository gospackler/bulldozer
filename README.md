# bulldozer
Bulldozer does shit load of work for you.

## Design
### Task
* Init() -> channel for run, updatesChan
* Run(  )
     * Listens on the channel and if data comes in, works on it.
     * If error occurs, response is sent on the update channel.
     * Once it is done result is sent to update channel if any.
* Finish() 
    * closes both the channels. 

### Bulldozer 
* Create a channel which is buffered channel of TaskRun Channels.
* Start a global Channel which listens for DataInput.
* Accepts a set of Tasks.
* Calls the Init() of all the tasks.
    * Start goroutines listening on updates.
    * Push each channel to the buffered channel list.
* Start goroutines that listens for updates.
* Have a function that is global which which accepts the input.
* Get one of the channels from the free list.
* Push data to it.

### Test

Create n workers. The run function appends worker id to the input and returns it.# bulldozer
