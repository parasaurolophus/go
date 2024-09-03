// Copyright 2024 Kirk Rader

package utilities

// Process items in a set of data concurrently. Specifically, start n+1
// goroutines and wait for them all to complete after invoking the given
// generator function. The generator function must send input values in a
// round-robin fashion to the set of transformer channels it is passed. The
// transformer goroutines will send the result of invoking the given transform
// function to the consumer goroutine, which invokes the given consume
// function:
//
//	                   +-----------+
//	              +-->>| transform |----+
//	              |    +-----------+    |
//	              |          .          |
//	+----------+  |          .          |    +---------+
//	| generate |--+     concurrent      |-->>| consume |
//	+----------+  |     goroutines      |    +---------+
//	              |          .          |
//	              |          .          |
//	              |    +-----------+    |
//	              +-->>| transform |----+
//	                   +-----------+
//
// Note that this function will hang if any of the generate, transform or
// consume functions do not return. If your transform function invokes some SDK
// function or API that can hang, consider the use of WithTimeLimit to allow
// the batch to run to completion even if some operations would otherwise block
// it (but then be aware of the consequences of resulting resource leaks).
//
// See CloseAndWait, CloseAllAndWait, StartWorker, StartWorkers
func ProcessBatch[Input any, Output any](

	numTransformers int,
	transformersBufferSize int,
	consumerBufferSize int,
	generate func(transformers []chan<- Input),
	transform func(Input) Output,
	consume func(Output),

) {

	// start a goroutine that will apply the consume function to each value
	// sent to its channel
	consumer, awaitConsumer := StartWorker(consumerBufferSize, consume)
	defer CloseAndWait(consumer, awaitConsumer)

	// wrap the transform function in a closure that will send a given
	// transformed input to the consumer channel
	produce := func(request Input) {
		consumer <- transform(request)
	}

	// start n goroutines each of which will call the produce closure for each
	// value sent to its channel
	transformers, awaitTransformers := StartWorkers(numTransformers, transformersBufferSize, produce)
	defer CloseAllAndWait(transformers, awaitTransformers)

	// generate must send values of type Input to the channels it is
	// passed,then return
	generate(transformers)
}
