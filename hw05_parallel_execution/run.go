package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"fmt"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func worker(doneCh <-chan interface{}, workCh <-chan Task, errCh chan error){
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

func interrupter(errCh <-chan error, doneCh chan interface{}, terminate chan error, M int) {
	for {
		select {
		case err, ok := <-errCh:
			if !ok {
				return
			}
			if err != nil {
				fmt.Println("We got error: %w", err)
				M--
				if M <= 0 {
					fmt.Println("TERMINATE!!")
					terminate <- fmt.Errorf("lalalala")
					return
				}
			}
		case <-doneCh:
			return
		}

	}
}

// Run starts tasks in N goroutines and stops its work when receiving M errors from tasks
func Run(tasks []Task, N int, M int) error {
	errCh := make(chan error, len(tasks))
	doneCh := make(chan interface{})
	//terminate := make(chan error)

	//go interrupter(errCh, doneCh, terminate, M)

	t := make(chan Task, N)

	wg := sync.WaitGroup{}
	for i := 0; i < N; i++ {
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

		if len(errCh) >= M {
			close(doneCh)
			return ErrErrorsLimitExceeded
		}
	}
	close(t)

	return nil
}
