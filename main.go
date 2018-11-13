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
)

func Usage() {
	fmt.Println("Usage: consul-dotenv-import <options> <env file>")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = Usage
	flag.StringVar(&consulToken, "token", "", "Consul token")
	flag.StringVar(&consulHost, "host", "http://127.0.0.1:8500", "Consul host")
	flag.StringVar(&consulPrefix, "prefix", "", "Consul KV prefix")
	flag.BoolVar(&dryRun, "dry", false, "Dry run (don't actual put values in Consul)")
	flag.Parse()

	if len(flag.Args()) != 1 {
		fmt.Println("No env file provided")
		flag.Usage()
		os.Exit(2)
	}

	filePath := flag.Arg(0)

	envFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println("could not open env file; ", err.Error())
		os.Exit(2)
	}
	defer envFile.Close()

	client, err := api.NewClient(&api.Config{Token: consulToken, Address: consulHost})
	if err != nil {
		fmt.Println("could not connect to Consul; ", err.Error())
		os.Exit(2)
	}

	sc := bufio.NewScanner(envFile)
	for sc.Scan() {
		line := sc.Text()
		lineParts := strings.Split(line, "=")
		if len(lineParts) == 2 {
			key := lineParts[0]
			value := lineParts[1]
			fullKey := strings.Trim(fmt.Sprintf("%s/%s", consulPrefix, key), "/")
			fmt.Printf("Putting %s = %s\n", fullKey, value)

			if !dryRun {
				p := &api.KVPair{Key: fullKey, Value: []byte(value)}
				_, err := client.KV().Put(p, nil)
				if err != nil {
					fmt.Println("could not save key in Consul; ", err.Error())
					os.Exit(2)
				}
			}
		} else {
			fmt.Printf("error, could not parse line '%s'\n", line)
		}
	}
}
