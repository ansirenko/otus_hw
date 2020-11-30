package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func worker(doneCh <-chan interface{}, workCh <-chan Task, errCh chan error) {
	for task := range workCh {
		select {
		case <-doneCh:
			return
		default:
			err := task()
			if err != nil {
				errCh <- err
			}
		}
	}
}

// Run starts tasks in N goroutines and stops its work when receiving M errors from tasks.
func Run(tasks []Task, n int, m int) error {
	if n <= 0 {
		return nil // Started without workers.
	}

	errCh := make(chan error, len(tasks))
	doneCh := make(chan interface{})

	t := make(chan Task, n)

	wg := sync.WaitGroup{}
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			worker(doneCh, t, errCh)
			wg.Done()
		}()
	}

	defer func() {
		wg.Wait()
	}()

	for _, tsk := range tasks {
		t <- tsk

		if len(errCh) >= m && m > 0 {
			close(doneCh)
			return ErrErrorsLimitExceeded
		}
	}
	close(t)

	return nil
}
