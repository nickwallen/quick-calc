package main

import (
	"fmt"
	"github.com/bcicen/go-units"
	"github.com/nickwallen/qcalc"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

func main() {
	// first arg is the output path
	outputPath := "qcalc-gen.csv"
	if len(os.Args) > 1 {
		outputPath = os.Args[1]
	}

	// second arg is the number of samples
	numSamples := 100_000
	if len(os.Args) > 2 {
		parsedInt, err := strconv.ParseInt(os.Args[2], 10, 32)
		if err != nil {
			panic(err)
		}
		numSamples = int(parsedInt)
	}

	numProducers := 1000
	start := time.Now()
	samples := make(chan string)
	done := make(chan int)

	// the consumer writes the samples to a file
	go consumer(numSamples, outputPath, samples, done)
	numWritten := waitForProducers(numProducers, samples, done)

	fmt.Printf("\nWrote %d/%d samples to '%s' with %d producer(s) in %s.\n",
		numWritten, numSamples, outputPath, numProducers, time.Since(start))
}

func waitForProducers(numWorkers int, samples chan string, done chan int) (numSamples int) {
	var wg sync.WaitGroup
	defer wg.Wait() // wait for the producers to complete

	// launch the workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			producer(samples, done)
		}()
	}

	// wait for the consumer to tell us how many were written
	numSamples = <-done
	return numSamples
}

func producer(output chan<- string, done <-chan int) {
	for {
		select {
		case output <- generateSample():
		case <-done:
			return // channel closed; consumer has enough samples
		}
	}
}

func consumer(numSamples int, outputPath string, input <-chan string, done chan<- int) {
	// signal to the producers that the consumer is done
	defer close(done)

	// open the file
	file, err := os.Create(outputPath)
	defer file.Close()
	if err != nil {
		panic("Unable to open output file: " + outputPath)
	}

	// consume the samples and write to disk
	i := 0
	defer func() { done <- i }()
	for ; i < numSamples; i++ {
		sample := <-input
		_, err := fmt.Fprintf(file, "%s \n", sample)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
}

func generateSample() string {
	for {
		input := randInput()
		output, err := qcalc.Calculate(input)
		if err == nil {
			// we have a good sample!
			sample := fmt.Sprintf("\"%s\", \"%s\"", input, output)
			return sample
		}
		// an invalid expression was generated like 1 meter + 3 pounds
	}
}

func randInput() (input string) {
	// the input expression can take a few forms
	switch rand.Intn(3) {
	case 0:
		// samples of the form "NUM UNIT in UNIT"
		input = fmt.Sprintf("%f %s in %s",
			randAmount(), randUnit(), randUnit())
	case 1:
		// samples in the form "NUM UNIT OP NUM UNIT"
		input = fmt.Sprintf("%.2f %s %s %.2f %s ",
			randAmount(), randUnit(), randOp(), randAmount(), randUnit())
	case 2:
		// samples in the form "NUM UNIT OP NUM INIT IN UNIT"
		input = fmt.Sprintf("%.2f %s %s %.2f %s in %s",
			randAmount(), randUnit(), randOp(), randAmount(), randUnit(), randUnit())
	default:
		panic("unexpected condition!")
	}
	return input
}

func randAmount() float64 {
	return rand.Float64() * float64(rand.Intn(100))
}

func randUnit() string {
	allUnits := units.All()
	index := rand.Intn(len(allUnits))
	randUnit := allUnits[index]
	if randUnit.Symbol != "" {
		return randUnit.Symbol
	}
	return randUnit.PluralName()
}

func randOp() string {
	ops := []string{"+", "-"}
	index := rand.Intn(2)
	return ops[index]
}
