# tts-go
A Go library for Azure's Cognitive Services tts API and Use Microsoft Edge's online tts Service

```bash
go run example/edge_tts.go --version

go run example/edge_tts.go --text="good boy"
go run example/edge_tts.go --proxy=socks5://127.0.0.1:7890 --text="good boy"
```

```bash
go run example/azure_tts.go --version

go run example/azure_tts.go --list-voices --locale=en-US --gender=Female

go run example/azure_tts.go --text="good boy"
```

```go
package main

import (
	"context"
	"fmt"
	"github.com/ericsource/tts-go/src/azure_tts"
	//"github.com/spf13/pflag"
)

func main() {
	
}
```