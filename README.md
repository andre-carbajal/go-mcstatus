# go-mcstatus

`go-mcstatus` is a Go library and command-line tool (CLI) for querying the status of Minecraft servers, supporting both **Java Edition** and **Bedrock Edition**.

## 🛠 Supported Protocols
- **Java Edition**: Server List Ping (SLP) via TCP.
- **Bedrock Edition**: Unconnected Ping (Query) via UDP.

## 🚀 CLI Usage

### Installation
If you have Go installed, you can install the binary directly:
```bash
go install github.com/andre-carbajal/go-mcstatus/cmd@latest
```

### Basic Commands
Query a **Java Edition** server (default port 25565):
```bash
mcstatus mc.hypixel.net
```

Query a **Bedrock Edition** server (default port 19132):
```bash
mcstatus --bedrock play.nethergames.org
```

Get **latency only** (ping):
```bash
mcstatus --ping mc.hypixel.net
```

---

## 📦 Usage as a Go Library

You can import `go-mcstatus` into your own project to integrate server status queries.

### Installation
```bash
go get github.com/andre-carbajal/go-mcstatus
```

### Usage Example (Java Edition)
```go
package main

import (
	"fmt"
	"log"
	"github.com/andre-carbajal/go-mcstatus"
)

func main() {
	// Create a new Java server instance
	server, err := mcstatus.NewJavaServer("mc.hypixel.net")
	if err != nil {
		log.Fatal(err)
	}

	// Get the full status
	status, err := server.Status()
	if err != nil {
		log.Fatal(err)
	}

	// Access the data (requires type assertion for specific fields)
	if resp, ok := status.(*mcstatus.JavaStatusResponse); ok {
		fmt.Printf("Version: %s\n", resp.Version.Name)
		fmt.Printf("Players: %d/%d\n", resp.Players.Online, resp.Players.Max)
		fmt.Printf("Latency: %dms\n", resp.GetLatency())
	}
}
```

### Usage Example (Bedrock Edition)
```go
package main

import (
	"fmt"
	"log"
	"github.com/andre-carbajal/go-mcstatus"
)

func main() {
	server, _ := mcstatus.NewBedrockServer("play.nethergames.org")
	
	status, err := server.Status()
	if err != nil {
		log.Fatal(err)
	}

	if resp, ok := status.(*mcstatus.BedrockStatusResponse); ok {
		fmt.Printf("MOTD: %s\n", resp.MOTD)
		fmt.Printf("Map: %s\n", resp.MapName)
	}
}
```

---

## 🛠 Project Structure
- `cmd/`: Contains the executable binary source code.
- 📁 Root (`/`): Contains the library logic (package `mcstatus`).

## 📄 License
This project is licensed under the MIT License.
