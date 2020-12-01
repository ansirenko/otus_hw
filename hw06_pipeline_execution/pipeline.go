package hw06_pipeline_execution //nolint:golint,stylecheck

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if len(stages) == 1 {
		return in
	}
	out := in
	for _, stage := range stages {
		out = stage(terminator(done, out))
	}
	return out
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
				out <- i
			}
		}
	}()
	return out
}
