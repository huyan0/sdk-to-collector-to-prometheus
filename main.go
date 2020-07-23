package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/sdk/metric/controller/push"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
)

func main() {
	// initiate OTLP exporter
	exporter, err := otlp.NewExporter(
		otlp.WithInsecure()) // Configure as needed.
	if err != nil {
		log.Fatalf("failed to create exporter: %v", err)
	}
	defer func() {
		err := exporter.Stop()
		if err != nil {
			log.Fatalf("failed to stop exporter: %v", err)
		}
	}()
	// initiate push controller
	pusher := push.New(simple.NewWithInexpensiveDistribution(), exporter, push.WithPeriod(1*time.Second))

	pusher.Start()

	meter := pusher.Provider().Meter("example")
	ctx := context.Background()

	// Create two instruments with Go SDK metric package
	counter0 := metric.Must(meter).NewInt64Counter(
		"a_counter_i",
		metric.WithDescription("Adds i every time, growth should be O(i^2)"),
	)
	counter1 := metric.Must(meter).NewInt64Counter(
		"a_counter_1",
		metric.WithDescription("Adds 1 every time, growth should be O(i)"),
	)

	recorder := metric.Must(meter).NewInt64ValueRecorder(
		"a_valuerecorder",
		metric.WithDescription("Records values"),
	)

	// Add initial values to the instruments
	counter1.Add(ctx, 0, kv.String("kind", "counter"))
	counter0.Add(ctx, 0, kv.String("kind", "counter"))
	recorder.Record(ctx, 0, kv.String("kind", "gauge"))
	var wg sync.WaitGroup
	wg.Add(1)
	value := 0
	// Repeatedly record values every 1 seconds
	go func() {
		for i := 1; i <= 100; i++ {
			time.Sleep(1 * time.Second)
			value += i
			fmt.Printf("%d. Recording %d in gauge and adding 1 to counter, counter0 should be %d, counter1 should %d \n", i, i, i, value)
			recorder.Record(ctx, int64(i), kv.String("kind", "gauge"))
			counter0.Add(ctx, int64(i), kv.String("kind", "counter"))
			counter1.Add(ctx, int64(1), kv.String("kind", "counter"))
		}
	}()
	wg.Wait()
}
