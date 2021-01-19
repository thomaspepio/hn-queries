package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/thomaspepio/hn-queries/endpoint"
	"github.com/thomaspepio/hn-queries/parser"

	"github.com/thomaspepio/hn-queries/index"
)

func main() {
	index := ingestHnLogs()
	startEndpoints(index)
}

func ingestHnLogs() *index.Index {
	now := time.Now()
	os.Stdout.WriteString(now.UTC().String() + " - Start indexing...\n")
	index := index.EmptyIndex()

	file, err := os.Open("./hn_logs.tsv")
	// file, err := os.Open("./test_data")
	// file, err := os.Open("./mini_test_data")
	if err != nil {
		panic(err.Error())
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	nbIndexed := 0
	for scanner.Scan() {
		line := scanner.Text()
		parsedQuery, parseError := parser.ParseHNQuery(line)

		if parseError != nil {
			os.Stdout.WriteString(parseError.Error() + "\n")
		} else {
			index.Add(parsedQuery)
			nbIndexed++
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	now = time.Now()
	os.Stdout.WriteString(now.UTC().String() + " - Indexing : OK\n")
	os.Stdout.WriteString(now.UTC().String() + " - " + strconv.Itoa(nbIndexed) + " log lines indexed\n")
	return index
}

func startEndpoints(index *index.Index) {
	if index == nil {
		panic("Could not start server : index is nil")
	}

	router := endpoint.Router(index)
	router.Run(":8080")
}
