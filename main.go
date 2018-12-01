package main

import (
	"flag"
	"log"
	"os/exec"
	"strings"
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
