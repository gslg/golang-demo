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
	shutdownWorkChannel  chan struct{}
	shutdownWaitGroup    sync.WaitGroup
	queueChannel         chan poolWork
	workChannel          chan PoolWorker
	queueWork            int32
	activesRoutines      int32
	queueCapacity        int32
}

// init is called when the system is inited.
func init() {
	log.SetPrefix("TRACE: ")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

// New creates a new WorkPool.
func New(numberOfRoutines int, queueCapacity int32) *WorkPool {
	pool := WorkPool{
		shutdownQueueChannel: make(chan string),
		shutdownWorkChannel:  make(chan struct{}),
		queueChannel:         make(chan poolWork),
		workChannel:          make(chan PoolWorker, queueCapacity),
		queueWork:            0,
		activesRoutines:      0,
		queueCapacity:        queueCapacity,
	}

	// Add the total number of routines to the wait group
	pool.shutdownWaitGroup.Add(numberOfRoutines)

	// Launch the work routines to process work
	for workRoutine := 0; workRoutine < numberOfRoutines; workRoutine++ {
		go pool.workRoutine(workRoutine)
	}
	// Start the queue routine to capture and provide work
	go pool.queueRoutine()

	return &pool

}

func (w *WorkPool) Shutdown(goRoutine string) (err error) {
	defer catchPanic(&err, goRoutine, "Shutdown")

	writeStdout(goRoutine, "Shutdown", "Started")
	writeStdout(goRoutine, "Shutdown", "Queue Routine")

	w.shutdownQueueChannel <- "Down"
	<-w.shutdownQueueChannel

	close(w.queueChannel)
	close(w.shutdownQueueChannel)

	writeStdout(goRoutine, "Shutdown", "Shutting Down Work Routines")

	close(w.shutdownWorkChannel)
	w.shutdownWaitGroup.Wait()

	close(w.workChannel)
	writeStdout(goRoutine, "Shutdown", "Completed")

	return err
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
func (w *WorkPool) PostWork(goRoutine string, work PoolWorker) (err error) {

	defer catchPanic(&err, goRoutine, "PostWork")

	poolWork := poolWork{work: work, resultChannel: make(chan error)}

	defer close(poolWork.resultChannel)

	w.queueChannel <- poolWork
	err = <-poolWork.resultChannel

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

func (w *WorkPool) workRoutine(workRoutine int) {
	for {
		select {
		case <-w.shutdownWorkChannel:
			writeStdout(fmt.Sprintf("WorkRoutine %d", workRoutine), "workRoutine", "Going Down")
			w.shutdownWaitGroup.Done()
			return
		case poolWorker := <-w.workChannel:
			w.safelyDoWork(workRoutine, poolWorker)
			break
		}
	}
}

func (w *WorkPool) safelyDoWork(workRoutine int, poolWorker PoolWorker) {
	defer catchPanic(nil, "WorkRoutine", "SafelyDoWork")
	defer atomic.AddInt32(&w.activesRoutines, -1)

	//更新状态
	atomic.AddInt32(&w.queueWork, -1)
	atomic.AddInt32(&w.activesRoutines, 1)

	//执行任务
	poolWorker.DoWork(workRoutine)
}

func (w *WorkPool) queueRoutine() {
	for {
		select {
		case <-w.shutdownQueueChannel:
			writeStdout("Queue", "queueRoutine", "Going Down")
			w.shutdownQueueChannel <- "Down"
			return
		case queueItem := <-w.queueChannel:
			//如果工作队列已经满了
			if atomic.AddInt32(&w.queueWork, 0) == w.queueCapacity {
				queueItem.resultChannel <- Error_Capacity
				continue
			}

			atomic.AddInt32(&w.queueWork, 1)

			w.workChannel <- queueItem.work
			//告诉调用者任务已经入队了
			queueItem.resultChannel <- nil

			break
		}
	}
}
