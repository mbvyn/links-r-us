package pipeline

import (
	"context"
	"golang.org/x/xerrors"
	"sync"
)

type fifo struct {
	proc Processor
}

// FIFO returns a StageRunner that processes incoming payloads in a first-in
// first-out fashion. Each input is passed to the specified processor and its
// output is emitted to the next stage.
func FIFO(proc Processor) StageRunner {
	return fifo{proc: proc}
}

// Run implements StageRunner
func (r fifo) Run(ctx context.Context, params StageParams) {
	for {
		select {
		case <-ctx.Done():
			return // Asked to cleanly shut down
		case payloadIn, ok := <-params.Input():
			if !ok {
				return // No more data available
			}

			payloadOut, err := r.proc.Process(ctx, payloadIn)
			if err != nil {
				wrappedErr := xerrors.Errorf("pipeline stage %d: %w", params.StageIndex(), err)
				maybeEmitError(wrappedErr, params.Error())
				return
			}

			// If the processor did not output a payload for the
			// next stage there is nothing we need to do.
			if payloadOut == nil {
				payloadIn.MarkAsProcessed()
				continue
			}

			// Output processed data
			select {
			case params.Output() <- payloadOut:
			case <-ctx.Done():
				return // Asked to cleanly shut down
			}
		}

	}
}

type fixedWorkerPool struct {
	fifos []StageRunner
}

// FixedWorkerPool returns a StageRunner that spins up a pool containing
// numWorkers to process incoming payloads in parallel and emit their outputs
// to the next stage.
func FixedWorkerPool(proc Processor, numWorkers int) StageRunner {
	if numWorkers <= 0 {
		panic("FixedWorkerPool: numWorkers must be > 0")
	}

	fifos := make([]StageRunner, numWorkers)
	for i := 0; i < numWorkers; i++ {
		fifos[i] = FIFO(proc)
	}

	return &fixedWorkerPool{fifos: fifos}
}

// Run implements StageRunner.
func (p *fixedWorkerPool) Run(ctx context.Context, params StageParams) {
	var wg sync.WaitGroup

	// Spin up each worker in the pool and wait for them to exit
	for i := 0; i < len(p.fifos); i++ {
		wg.Add(1)
		go func(fifoIndex int) {
			p.fifos[fifoIndex].Run(ctx, params)
			wg.Done()
		}(i)
	}
	
	wg.Wait()
}

// maybeEmitError attempts to queue err to a buffered error channel. If the
// channel is full, the error is dropped.
func maybeEmitError(err error, errCh chan<- error) {
	select {
	case errCh <- err: // error emitted.
	default: // error channel is full of other errors.
	}
}
