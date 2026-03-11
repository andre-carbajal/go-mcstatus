package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"gomcstatus"
)

func main() {
	bedrock := flag.Bool("bedrock", false, "Query a Bedrock server")
	ping := flag.Bool("ping", false, "Only return latency")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Usage: gomcstatus [--bedrock] [--ping] <address>")
		os.Exit(1)
	}

	address := args[0]

	var server gomcstat.Server
	var err error

	if *bedrock {
		server, err = gomcstat.NewBedrockServer(address)
	} else {
		server, err = gomcstat.NewJavaServer(address)
	}

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if *ping {
		latency, err := server.Ping()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Latency: %dms\n", latency)
		return
	}

	status, err := server.Status()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	b, _ := json.MarshalIndent(status, "", "  ")
	fmt.Println(string(b))
}
