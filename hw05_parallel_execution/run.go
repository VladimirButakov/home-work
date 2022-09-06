package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var (
	ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
	ErrNoWorkers           = errors.New("no workers found")
	ErrMaxErrorsIsNotValid = errors.New("max errors number is too small")
)

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	waitGroup := sync.WaitGroup{}
	tasksChanel := make(chan Task)
	var errorsCount int64
	var err error

	if n <= 0 {
		return ErrNoWorkers
	}

	if m <= 0 {
		return ErrMaxErrorsIsNotValid
	}

	waitGroup.Add(n)

	for i := 0; i < n; i++ {
		go func() {
			for v := range tasksChanel {
				if v == nil {
					continue
				}

				err := v()
				if err != nil {
					atomic.AddInt64(&errorsCount, 1)
				}
			}

			waitGroup.Done()
		}()
	}

	for _, task := range tasks {
		if atomic.LoadInt64(&errorsCount) >= int64(m) {
			err = ErrErrorsLimitExceeded

			break
		}

		tasksChanel <- task
	}

	close(tasksChanel)

	waitGroup.Wait()

	return err
}
