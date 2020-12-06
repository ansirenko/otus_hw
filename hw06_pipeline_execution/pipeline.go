package hw06_pipeline_execution //nolint:golint,stylecheck

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if len(stages) == 1 && stages[0] == nil {
		return in
	}
	for _, stage := range stages {
		in = stage(terminator(done, in))
	}
	return in
}

func terminator(done In, in In) Out {
	out := make(Bi)
	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				return
			case i, ok := <-in:
				if !ok {
					return
				}
				select {
				case <-done:
					return
				case out <- i:
				}
			}
		}
	}()
	return out
}
