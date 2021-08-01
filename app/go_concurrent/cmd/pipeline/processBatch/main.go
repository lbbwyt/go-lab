package main

import (
	"context"
	"github.com/deliveryhero/pipeline"
	"github.com/deliveryhero/pipeline/example/processors"
	"log"
	"time"
)

func main() {
	// Create a context that times out after 5 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create a pipeline that emits 1-6 at a rate of one int per second
	p := pipeline.Delay(ctx, time.Second, pipeline.Emit(1, 2, 3, 4, 5, 6))

	// Use the BatchMultiplier to multiply 2 adjacent numbers together
	p = pipeline.ProcessBatch(ctx, 2, time.Minute, &processors.BatchMultiplier{}, p)

	// Finally, lets print the results and see what happened
	for result := range p {
		log.Printf("result: %d\n", result)
	}

	// Output
	//2021/07/01 10:11:36 result: 2
	//2021/07/01 10:11:38 result: 12
	//2021/07/01 10:11:40 error: could not multiply [5 6], context deadline exceeded
}
