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

func usage() {
	fmt.Println("usage: edge-tts [-h] [-t TEXT] [-f FILE] [-v VOICE] [-l] [--rate RATE] [--volume VOLUME] [--words-in-cue WORDS_IN_CUE] [--write-media WRITE_MEDIA] [--write-subtitles WRITE_SUBTITLES] [--proxy PROXY]\n")
	fmt.Println("Microsoft Edge TTS\n")
	fmt.Println("options:")
	pflag.PrintDefaults()
}

import (
	"context"
	"fmt"
	"github.com/ericsource/tts-go/src/azure_tts"
	//"github.com/spf13/pflag"
)

func main() {
	
}
```