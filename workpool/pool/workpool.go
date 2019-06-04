package pool

import (
	"errors"
	"fmt"
	"log"
	"runtime"
	"sync"
	"sync/atomic"
)

var (
	Error_Capacity = errors.New("Thread pool already full")
)

type PoolWorker interface {
	DoWork(workRoutine int)
}

type poolWork struct {
	work          PoolWorker
	resultChannel chan error
}

type WorkPool struct {
	shutdownQueueChannel chan string
	shutdownWorkChannel chan struct{}
	shutdownWaitGroup sync.WaitGroup
	queueChannel chan poolWork
	workChannel chan PoolWorker
	queueWork int32
	activesRoutines int32
	queueCapacity int32
}

// init is called when the system is inited.
func init() {
	log.SetPrefix("TRACE: ")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func New(numberOfRoutines int,queueCapacity int32) *WorkPool{
	pool := WorkPool{
		shutdownQueueChannel:make(chan string),
		shutdownWorkChannel:make(chan struct{}),
		queueChannel:make(chan poolWork),
		workChannel:make(chan PoolWorker,queueCapacity),
		queueWork:0,
		activesRoutines:0,
		queueCapacity:queueCapacity,
	}

	pool.shutdownWaitGroup.Add(numberOfRoutines)

	return &pool

}

// writeStdout is used to write a system message directly to stdout.
func writeStdout(goRoutine string, functionName string, message string) {
	log.Printf("%s : %s : %s\n", goRoutine, functionName, message)
}

// writeStdoutf is used to write a formatted system message directly stdout.
func writeStdoutf(goRoutine string, functionName string, format string, a ...interface{}) {
	writeStdout(goRoutine, functionName, fmt.Sprintf(format, a...))
}

// CatchPanic is used to catch any Panic and log exceptions to Stdout. It will also write the stack trace.
func catchPanic(err *error, goRoutine string, functionName string) {
	if r := recover(); r != nil {
		// Capture the stack trace
		buf := make([]byte, 10000)
		runtime.Stack(buf, false)

		writeStdoutf(goRoutine, functionName, "PANIC Defered [%v] : Stack Trace : %v", r, string(buf))

		if err != nil {
			*err = fmt.Errorf("%v", r)
		}
	}
}

// PostWork will post work into the WorkPool. This call will block until the Queue routine reports back
// success or failure that the work is in queue.
func (w *WorkPool) PostWork(goRoutine string,work PoolWorker) (err error){

	defer catchPanic(&err,goRoutine,"PostWork")

	poolWork := poolWork{work:work,resultChannel:make(chan error)}

	defer close(poolWork.resultChannel)

	w.queueChannel <- poolWork
	err = <- poolWork.resultChannel

	return err
}

// QueuedWork will return the number of work items in queue.
func (w *WorkPool) QueuedWork() int32 {
	return atomic.AddInt32(&w.queueWork, 0)
}

// ActiveRoutines will return the number of routines performing work.
func (w *WorkPool) ActiveRoutines() int32 {
	return atomic.AddInt32(&w.activesRoutines, 0)
}