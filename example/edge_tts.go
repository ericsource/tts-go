package main

import (
	"fmt"
	"github.com/ericsource/tts-go/src/edge_tts"
	"github.com/spf13/pflag"
	"os"
)

func usage() {
	fmt.Println("usage: edge-tts [-h] [-t TEXT] [-f FILE] [-v VOICE] [-l] [--rate RATE] [--volume VOLUME] [--words-in-cue WORDS_IN_CUE] [--write-media WRITE_MEDIA] [--write-subtitles WRITE_SUBTITLES] [--proxy PROXY]\n")
	fmt.Println("Microsoft Edge TTS\n")
	fmt.Println("options:")
	pflag.PrintDefaults()
}

func main() {
	listVoices := pflag.BoolP("list-voices", "l", false, "lists available voices and exits")
	locale := pflag.StringP("locale", "", "en-US", "locale for voice lists ex: zh-CN, en-US")
	gender := pflag.StringP("gender", "", "", "gender for voice lists ex: Male, Female")
	text := pflag.StringP("text", "t", "", "what TTS will say")
	//file := pflag.StringP("file", "f", "", "same as --text but read from file")
	voice := pflag.StringP("voice", "v", "zh-CN-XiaoxiaoNeural", "voice for TTS")
	volume := pflag.String("volume", "+0%", "set TTS volume")
	rate := pflag.String("rate", "+0%", "set TTS rate")
	writeMedia := pflag.String("write-media", "example/media.mp3", "send media output to file instead of stdout")
	proxy := pflag.String("proxy", "", "use a proxy for TTS and voice list")
	version := pflag.BoolP("version", "", false, "show version")
	pflag.Usage = usage
	pflag.Parse()

	if *version {
		major, minor, patch := edge_tts.ParseVersion()
		fmt.Printf("edge-tts version %v.%v.%v\n", major, minor, patch)
		return
	}

	if *listVoices {
		ListVoices, _ := edge_tts.NewListVoices(*proxy)
		voives := ListVoices.Find(*gender, *locale)
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
		tts := edge_tts.NewCommunicate().WithVoice(*voice).WithRate(*rate).WithVolume(*volume)
		file, err := os.OpenFile(*writeMedia, os.O_APPEND|os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
		if err != nil {
			fmt.Printf("Failed to open file: %v\n", err)
			return
		}
		tts.OpenWs()
		defer tts.Close()
		defer file.Close()

		texts := []*edge_tts.CommunicateTextTask{}
		texts = append(texts, &edge_tts.CommunicateTextTask{
			Option: tts.GetOption(),
			Text:   *text,
		})
		go tts.AllocateTask(texts)
		tts.CreatePool()
		for _, text := range texts {
			//fmt.Println(len(text.SpeechData))
			n, err := file.Write(text.SpeechData)
			if err != nil {
				fmt.Println("err:", err)
			}
			fmt.Println("media len:", n)
		}
	}
}
