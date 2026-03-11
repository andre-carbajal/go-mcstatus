# gomcstat

`gomcstat` is a Go library and command-line tool (CLI) for querying the status of Minecraft servers, supporting both **Java Edition** and **Bedrock Edition**.

## 🛠 Supported Protocols
- **Java Edition**: Server List Ping (SLP) via TCP.
- **Bedrock Edition**: Unconnected Ping (Query) via UDP.

## 🚀 CLI Usage

### Installation
If you have Go installed, you can install the binary directly:
```bash
go install github.com/andre-carbajal/gomcstat/cmd/gomcstat@latest
```

### Basic Commands
Query a **Java Edition** server (default port 25565):
```bash
gomcstat mc.hypixel.net
```

Query a **Bedrock Edition** server (default port 19132):
```bash
gomcstat --bedrock play.nethergames.org
```

Get **latency only** (ping):
```bash
gomcstat --ping mc.hypixel.net
```

---

## 📦 Usage as a Go Library

You can import `gomcstat` into your own project to integrate server status queries.

### Installation
```bash
go get github.com/andre-carbajal/gomcstat
```

### Usage Example (Java Edition)
```go
package main

import (
	"fmt"
	"log"
	"gomcstat"
)

func main() {
	// Create a new Java server instance
	server, err := gomcstat.NewJavaServer("mc.hypixel.net")
	if err != nil {
		log.Fatal(err)
	}

	// Get the full status
	status, err := server.Status()
	if err != nil {
		log.Fatal(err)
	}

	// Access the data (requires type assertion for specific fields)
	if resp, ok := status.(*gomcstat.JavaStatusResponse); ok {
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
	"gomcstat"
)

func main() {
	server, _ := gomcstat.NewBedrockServer("play.nethergames.org")
	
	status, err := server.Status()
	if err != nil {
		log.Fatal(err)
	}

	if resp, ok := status.(*gomcstat.BedrockStatusResponse); ok {
		fmt.Printf("MOTD: %s\n", resp.MOTD)
		fmt.Printf("Map: %s\n", resp.MapName)
	}
}
```

---

## 🛠 Project Structure
- `cmd/`: Contains the executable binary source code.
- 📁 Root (`/`): Contains the library logic (package `gomcstat`).

## 📄 License
This project is licensed under the MIT License.
