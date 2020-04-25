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
	numSamples := 1_000_000
	if len(os.Args) > 2 {
		parsedInt, err := strconv.ParseInt(os.Args[2], 10, 32)
		if err != nil {
			panic(err)
		}
		numSamples = int(parsedInt)
	}

	start := time.Now()
	numWorkers := 10
	samples := make(chan string)
	done := make(chan bool)

	// the consumer writes the samples to a file
	go consumer(numSamples, outputPath, samples, done)

	// the producers generate the samples
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go producer(samples, done, &wg)
	}

	// wait for the generators to complete
	wg.Wait()
	elapsed := time.Since(start)
	fmt.Printf("\nWrote %d samples to '%s' with %d workers in %s.\n",
		numSamples, outputPath, numWorkers, elapsed)
}

func producer(output chan<- string, done <-chan bool, wg *sync.WaitGroup) {
	// signals to main when the worker is done
	defer wg.Done()
	for {
		sample := generateSample()
		select {
		case output <- sample:
		case <-done:
			return // our work here is done
		}
	}
}

func consumer(numSamples int, outputPath string, input <-chan string, done chan<- bool) {
	// open the file
	file, err := os.Create(outputPath)
	defer file.Close()
	if err != nil {
		panic("Unable to open output file: " + outputPath)
	}

	for i := 0; i < numSamples; i++ {
		sample := <-input
		fmt.Fprintf(file, "%s \n", sample)
	}

	// signal to the producers that nothing more is needed
	close(done)
}

func generateSample() string {
	input := randInput()
	result, err := qcalc.Calculate(input)

	var output string
	if err != nil && rand.Intn(100) < 20 {
		output = err.Error()
	} else {
		output = result
	}

	// write the sample to the output
	return fmt.Sprintf("\"%s\", \"%s\"", input, output)
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
