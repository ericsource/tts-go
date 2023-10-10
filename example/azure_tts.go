package main

import (
	"context"
	"fmt"
	"github.com/ericsource/tts-go/src/azure_tts"
	"github.com/spf13/pflag"
	"os"
)

func main() {
	listVoices := pflag.BoolP("list-voices", "l", false, "lists available voices and exits")
	locale := pflag.StringP("locale", "", "en-US", "locale for voice lists ex: zh-CN, en-US")
	gender := pflag.StringP("gender", "", "", "gender for voice lists ex: Male, Female")
	text := pflag.StringP("text", "t", "", "what TTS will say")
	////file := pflag.StringP("file", "f", "", "same as --text but read from file")
	//voice := pflag.StringP("voice", "v", "zh-CN-XiaoxiaoNeural", "voice for TTS")
	//volume := pflag.String("volume", "+0%", "set TTS volume")
	//rate := pflag.String("rate", "+0%", "set TTS rate")
	writeMedia := pflag.String("write-media", "example/media.mp3", "send media output to file instead of stdout")
	//proxy := pflag.String("proxy", "", "use a proxy for TTS and voice list")
	version := pflag.BoolP("version", "", false, "show version")
	pflag.Usage = func() {
		fmt.Println("usage: edge-tts [-h] [-t TEXT] [-f FILE] [-v VOICE] [-l] [--rate RATE] [--volume VOLUME] [--words-in-cue WORDS_IN_CUE] [--write-media WRITE_MEDIA] [--write-subtitles WRITE_SUBTITLES] [--proxy PROXY]\n")
		fmt.Println("Microsoft Edge TTS\n")
		fmt.Println("options:")
		pflag.PrintDefaults()
	}
	pflag.Parse()

	if *version {
		major, minor, patch := azure_tts.ParseVersion()
		fmt.Printf("edge-tts version %v.%v.%v\n", major, minor, patch)
		return
	}

	if *listVoices {
		tts, err := azure_tts.New("xxxx", azure_tts.RegionEastUS)
		if err != nil {
			fmt.Printf("Failed to new azure tts: %v\n", err)
			return
		}
		tts.FetchVoiceList()
		voives := tts.Find(*gender, *locale)
		// Name: Microsoft Server Speech Text to Speech Voice (ar-SA, ZariyahNeural)
		// ShortName: ar-SA-ZariyahNeural
		// Gender: Female
		// Locale: ar-SA
		for _, voive := range voives {
			fmt.Printf("%s:%v\n", "Name", voive.Name)
			fmt.Printf("%s:%v\n", "ShortName", voive.ShortName)
			fmt.Printf("%s:%v\n", "Gender", voive.Gender)
			fmt.Printf("%s:%v\n\n", "Locale", voive.Locale)
		}
		fmt.Printf("TOTAL:%v\n", len(voives))
		return
	}

	if len(*text) > 0 {
		//fmt.Println(*text)
		tts, err := azure_tts.New("xxxx", azure_tts.RegionEastUS)
		if err != nil {
			fmt.Printf("Failed to new azure tts: %v\n", err)
			return
		}
		file, err := os.OpenFile(*writeMedia, os.O_APPEND|os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
		if err != nil {
			fmt.Printf("Failed to open file: %v\n", err)
			return
		}
		defer file.Close()

		// Digitize a text string using the enUS locale, female voice and specify the
		// audio format of a 16Khz, 32kbit mp3 file.
		ctx := context.Background()
		b, err := tts.SynthesizeWithContext(
			ctx,
			*text,
			"en-US",
			"Female",
			"en-US-NancyNeural",
			azure_tts.Audio16khz32kbitrateMonoMp3)

		if err != nil {
			fmt.Printf("unable to synthesize, received: %v\n", err)
			return
		}

		n, err := file.Write(b)
		if err != nil {
			fmt.Println("file write err:", err)
			return
		}
		fmt.Println("media len:", n)
		return
	}
}
