package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/hashicorp/consul/api"
	"os"
	"strings"
)

var (
	consulToken  string
	consulHost   string
	consulPrefix string
	dryRun       bool
	verbose      bool
)

func Usage() {
	fmt.Println("Usage: consul-dotenv-import <options> <env file>")
	flag.PrintDefaults()
}

func log(message string) {
	fmt.Println(message)
}

func errorLog(message string) {
	log(message)
}

func debugLog(message string) {
	if verbose {
		log(message)
	}
}

func main() {
	flag.Usage = Usage
	flag.StringVar(&consulToken, "token", "", "Consul token")
	flag.StringVar(&consulHost, "host", "http://127.0.0.1:8500", "Consul host")
	flag.StringVar(&consulPrefix, "prefix", "", "Consul KV prefix")
	flag.BoolVar(&dryRun, "dry", false, "Dry run (don't actual put values in Consul)")
	flag.BoolVar(&verbose, "verbose", false, "Verbose output")
	flag.Parse()

	if len(flag.Args()) != 1 {
		log("No env file provided")
		flag.Usage()
		os.Exit(2)
	}

	filePath := flag.Arg(0)

	envFile, err := os.Open(filePath)
	if err != nil {
		errorLog(fmt.Sprintf("could not open env file; %s", err.Error()))
		os.Exit(2)
	}
	defer envFile.Close()

	client, err := api.NewClient(&api.Config{Token: consulToken, Address: consulHost})
	if err != nil {
		errorLog(fmt.Sprintf("could not connect to Consul; %s", err.Error()))
		os.Exit(2)
	}

	sc := bufio.NewScanner(envFile)
	for sc.Scan() {
		line := sc.Text()
		if line != "" && strings.Index(line, "#") != 0 {
			lineParts := strings.Split(line, "=")
			if len(lineParts) == 2 {
				key := lineParts[0]
				value := lineParts[1]
				fullKey := strings.Trim(fmt.Sprintf("%s/%s", consulPrefix, key), "/")
				log(fmt.Sprintf("Putting %s = %s", fullKey, value))

				if !dryRun {
					p := &api.KVPair{Key: fullKey, Value: []byte(value)}
					_, err := client.KV().Put(p, nil)
					if err != nil {
						errorLog(fmt.Sprintf("could not save key in Consul; %s", err.Error()))
						os.Exit(2)
					}
				}
			} else {
				debugLog(fmt.Sprintf("could not parse line '%s', skipping", line))
			}
		}
	}
}
