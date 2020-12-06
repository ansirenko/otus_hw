package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"fmt"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type worker struct {
	sync.Mutex
	errCount int
}

func (w *worker) work(doneCh <-chan interface{}, workCh <-chan Task) {
	for {
		select {
		case <-doneCh:
			return
		case task, ok := <-workCh:
			if !ok {
				return
			}
			if err := task(); err != nil {
				w.Lock()
				w.errCount++
				w.Unlock()
			}
		}
	}
}

// Run starts tasks in N goroutines and stops its work when receiving M errors from tasks.
func Run(tasks []Task, n int, m int) error {
	if n <= 0 {
		return fmt.Errorf("workers count must be greater than 0") // Started without workers.
	}

	doneCh := make(chan interface{})
	taskCh := make(chan Task)

	taskRunner := worker{}

	wg := sync.WaitGroup{}
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			taskRunner.work(doneCh, taskCh)
			wg.Done()
		}()
	}
	defer wg.Wait()

	for _, task := range tasks {
		taskCh <- task
		if m <= 0 {
			continue
		}
		taskRunner.Lock()
		currentErrCounters := taskRunner.errCount
		taskRunner.Unlock()
		if currentErrCounters >= m {
			close(doneCh)
			close(taskCh)
			return ErrErrorsLimitExceeded
		}
	}
	close(taskCh)

	return nil
}
