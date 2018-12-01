package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

var (
	rate         *int = flag.Int("rate", 0, "rate limit of execution")
	inflight     *int = flag.Int("inflight", 1, "count of parallel execution")
	arguments    []string
	commandArgId int = -1
	binnary      string
)

func main() {
	loadFlags()

	log.Printf("rate: %d\n", *rate)
	log.Printf("inflignt: %d\n", *inflight)
	log.Printf("command: %s %s\n", binnary, strings.Join(arguments, " "))

	wg := &sync.WaitGroup{}
	workerCh := make(chan string, *inflight)
	for i := 0; i < *inflight; i++ {
		wg.Add(1)
		go execWorker(workerCh, wg)
	}

	scanner := bufio.NewScanner(os.Stdin)
	// First read without rate limiting.
	if scanner.Scan() {
		workerCh <- scanner.Text()
	}
	limiter := time.Tick(time.Duration(1000000000 / *rate) * time.Nanosecond)
	for scanner.Scan() {
		<-limiter
		workerCh <- scanner.Text()
	}
	close(workerCh)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	wg.Wait()
}

func execWorker(ch chan string, wg *sync.WaitGroup) {
	args := make([]string, len(arguments))
	copy(args, arguments)
	for arg := range ch {
		if commandArgId >= 0 {
			args[commandArgId] = arg
		}
		cmd := exec.Command(binnary, args...)
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
	binnary, err = exec.LookPath(arguments[0])
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
