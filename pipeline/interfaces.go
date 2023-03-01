package pipeline

import "context"

// Payload is implemented by values that can be sent through a pipeline.
type Payload interface {
	// Clone returns a new Payload that is a deep-copy of the original.
	Clone() Payload

	// MarkAsProcessed is invoked by the pipeline when the Payload either
	// reaches the pipeline sink or it gets discarded by one of the
	// pipeline stages.
	MarkAsProcessed()
}

// Processor is implemented by types that can process Payloads as part of a
// pipeline stage.
type Processor interface {
	// Process operates on the input payload and returns back a new payload
	// to be forwarded to the next pipeline stage. Processors may also opt
	// to prevent the payload from reaching the rest of the pipeline by
	// returning a nil payload value instead.
	Process(context.Context, Payload) (Payload, error)
}

// ProcessorFunc is an adapter to allow the use of plain functions as Processor
// instances. If f is a function with the appropriate signature, ProcessorFunc(f)
// is a Processor that calls f.
type ProcessorFunc func(context.Context, Payload) (Payload, error)
