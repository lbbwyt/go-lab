/*
Copyright 2020 Deepak S<deepaks@outlook.in>
*/

// Package goworkers implements a simple, flexible and lightweight
// goroutine worker pool implementation.
package goworkers

import (
	"sync"
	"sync/atomic"
)

const (
	// The size of the buffered queue where jobs are queued up if no
	// workers are available to process the incoming jobs, unless specified
	defaultQSize = 128
	// A comfortable size for the buffered output channel such that chances
	// for a slow receiver to miss updates are minute
	outputChanSize = 100
)

// GoWorkers is a collection of worker goroutines.
//
// All workers will be killed after Stop() is called if their respective job finishes.
type GoWorkers struct {
	numWorkers uint32
	maxWorkers uint32
	numJobs    uint32
	workerQ    chan func()
	bufferedQ  chan func()
	jobQ       chan func()
	stopping   int32
	done       chan struct{}
	// ErrChan is a safe buffered output channel of size 100 on which error
	// returned by a job can be caught, if any. The channel will be closed
	// after Stop() returns. Valid only for SubmitCheckError() and SubmitCheckResult().
	// You must start listening to this channel before submitting jobs so that no
	// updates would be missed. This is comfortably sized at 100 so that chances
	// that a slow receiver missing updates would be minute.
	ErrChan chan error
	// ResultChan is a safe buffered output channel of size 100 on which error
	// and output returned by a job can be caught, if any. The channels will be
	// closed after Stop() returns. Valid only for SubmitCheckResult().
	// You must start listening to this channel before submitting jobs so that no
	// updates would be missed. This is comfortably sized at 100 so that chances
	// that a slow receiver missing updates would be minute.
	ResultChan chan interface{}
}

// Options configures the behaviour of worker pool.
//
// Workers specifies the number of workers that will be spawned.
// If unspecified or zero, workers will be spawned as per demand.
//
// QSize specifies the size of the queue that holds up incoming jobs.
// Minimum value is 128.
type Options struct {
	Workers uint32
	QSize   uint32
}

// New creates a new worker pool.
//
// Accepts optional Options{} argument.
func New(args ...Options) *GoWorkers {
	gw := &GoWorkers{
		workerQ: make(chan func()),
		// Do not remove jobQ. To stop receiving input once Stop() is called
		jobQ:       make(chan func()),
		ErrChan:    make(chan error, outputChanSize),
		ResultChan: make(chan interface{}, outputChanSize),
		done:       make(chan struct{}),
	}

	gw.bufferedQ = make(chan func(), defaultQSize)
	if len(args) == 1 {
		gw.maxWorkers = args[0].Workers
		if args[0].QSize > defaultQSize {
			gw.bufferedQ = make(chan func(), args[0].QSize)
		}
	}

	go gw.start()

	return gw
}

// JobNum returns number of active jobs
func (gw *GoWorkers) JobNum() uint32 {
	return atomic.LoadUint32(&gw.numJobs)
}

// WorkerNum returns number of active workers
func (gw *GoWorkers) WorkerNum() uint32 {
	return atomic.LoadUint32(&gw.numWorkers)
}

// Submit is a non-blocking call with arg of type `func()`
func (gw *GoWorkers) Submit(job func()) {
	if atomic.LoadInt32(&gw.stopping) == 1 {
		return
	}
	atomic.AddUint32(&gw.numJobs, uint32(1))
	gw.jobQ <- func() { job() }
}

// SubmitCheckError is a non-blocking call with arg of type `func() error`
//
// Use this if your job returns 'error'.
// Use ErrChan buffered channel to read error, if any.
func (gw *GoWorkers) SubmitCheckError(job func() error) {
	if atomic.LoadInt32(&gw.stopping) == 1 {
		return
	}
	atomic.AddUint32(&gw.numJobs, uint32(1))
	gw.jobQ <- func() {
		err := job()
		if err != nil {
			select {
			case gw.ErrChan <- err:
			default:
			}
		}
	}
}

// SubmitCheckResult is a non-blocking call with arg of type `func() (interface{}, error)`
//
// Use this if your job returns output and error.
// Use ErrChan buffered channel to read error, if any.
// Use ResultChan buffered channel to read output, if any.
// For a job, either of error or output would be sent if available.
func (gw *GoWorkers) SubmitCheckResult(job func() (interface{}, error)) {
	if atomic.LoadInt32(&gw.stopping) == 1 {
		return
	}
	atomic.AddUint32(&gw.numJobs, uint32(1))
	gw.jobQ <- func() {
		result, err := job()
		if err != nil {
			select {
			case gw.ErrChan <- err:
			default:
			}
		} else {
			select {
			case gw.ResultChan <- result:
			default:
			}
		}
	}
}

// Wait waits for the jobs to finish running.
//
// This is a blocking call and returns when all the active and queued jobs are finished.
// If 'wait' argument is set true, Wait() waits until the result and the error channels are emptied.
// Setting 'wait' argument to true ensures that you can read all the values from the result and
// the error channels before this function unblocks.
// Jobs cannot be submitted until this function returns. If any, will be discarded.
func (gw *GoWorkers) Wait(wait bool) {
	if !atomic.CompareAndSwapInt32(&gw.stopping, 0, 1) {
		return
	}

	for {
		if gw.JobNum() == 0 {
			break
		}
	}

	if wait {
		for {
			if len(gw.ResultChan)|len(gw.ErrChan) == 0 {
				break
			}
		}
	}

	atomic.StoreInt32(&gw.stopping, 0)
}

// Stop gracefully waits for the jobs to finish running and releases the associated resources.
//
// This is a blocking call and returns when all the active and queued jobs are finished.
// If wait is true, Stop() waits until the result and the error channels are emptied.
// Setting wait to true ensures that you can read all the values from the result and the
// error channels before your parent program exits.
func (gw *GoWorkers) Stop(wait bool) {
	if !atomic.CompareAndSwapInt32(&gw.stopping, 0, 1) {
		return
	}
	if gw.JobNum() != 0 {
		<-gw.done
	}

	if wait {
		for {
			if len(gw.ResultChan)|len(gw.ErrChan) == 0 {
				break
			}
		}
	}

	// close the input channel
	close(gw.jobQ)
}

var mx sync.Mutex

func (gw *GoWorkers) spawnWorker() {
	defer mx.Unlock()
	mx.Lock()
	if ((gw.maxWorkers == 0) || (gw.WorkerNum() < gw.maxWorkers)) && (gw.JobNum() > gw.WorkerNum()) {
		go gw.startWorker()
	}
}

func (gw *GoWorkers) start() {
	defer func() {
		close(gw.bufferedQ)
		close(gw.workerQ)
		close(gw.ErrChan)
		close(gw.ResultChan)
	}()

	// start a worker in advance
	go gw.startWorker()

	go func() {
		for {
			select {
			// keep processing the queued jobs
			case job, ok := <-gw.bufferedQ:
				if !ok {
					return
				}
				go func() {
					gw.spawnWorker()
					gw.workerQ <- job
				}()
			}
		}
	}()

	for {
		select {
		case job, ok := <-gw.jobQ:
			if !ok {
				return
			}
			select {
			// if possible, process the job without queueing
			case gw.workerQ <- job:
				go gw.spawnWorker()
			// queue it if no workers are available
			default:
				gw.bufferedQ <- job
			}
		}
	}
}

func (gw *GoWorkers) startWorker() {
	defer func() {
		atomic.AddUint32(&gw.numWorkers, ^uint32(0))
	}()

	atomic.AddUint32(&gw.numWorkers, 1)

	for job := range gw.workerQ {
		job()
		if (atomic.AddUint32(&gw.numJobs, ^uint32(0)) == 0) && (atomic.LoadInt32(&gw.stopping) == 1) {
			gw.done <- struct{}{}
		}
	}
}
