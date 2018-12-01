package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"
)

var (
	rate         = flag.Int("rate", 0, "rate limit of execution")
	inflight     = flag.Int("inflight", 1, "count of parallel execution")
	arguments    []string
	commandArgId = -1
	binary       string
)

func main() {
	loadFlags()

	wg := &sync.WaitGroup{}
	workerCh := make(chan string, *inflight)
	for i := 0; i < *inflight; i++ {
		wg.Add(1)
		go execWorker(workerCh, wg)
	}

	var limiter <-chan time.Time
	if *rate > 0 {
		limiter = time.Tick((time.Second / time.Duration(*rate)) * time.Nanosecond)
	}

	scanner := bufio.NewScanner(os.Stdin)
	// First read without rate limiting.
	if scanner.Scan() {
		workerCh <- scanner.Text()
	}
	for scanner.Scan() {
		if limiter != nil {
			<-limiter
		}
		workerCh <- scanner.Text()
	}
	close(workerCh)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	wg.Wait()
}

func execWorker(ch <-chan string, wg *sync.WaitGroup) {
	args := make([]string, len(arguments))
	copy(args, arguments)
	for arg := range ch {
		if commandArgId >= 0 {
			args[commandArgId] = arg
		}
		cmd := exec.Command(binary, args...)
		cmd.Env = os.Environ()

		if b, err := cmd.Output(); err != nil {
			log.Println(err)
		} else {
			fmt.Print(string(b))
		}
	}
	wg.Done()
}

func loadFlags() {
	flag.Parse()
	arguments = flag.Args()
	if len(arguments) < 1 {
		log.Fatalf("No arguments supply")
	}
	var err error
	binary, err = exec.LookPath(arguments[0])
	if err != nil {
		log.Fatal(err)
	}
	arguments = arguments[1:]
	for i, s := range arguments {
		if s == "{}" {
			commandArgId = i
			break
		}
	}
}
