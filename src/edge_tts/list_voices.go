package edge_tts

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Voice struct {
	Name      string `json:"Name"`
	ShortName string `json:"ShortName"`
	Gender    string `json:"Gender"`
	Locale    string `json:"Locale"`
}

type ListVoices struct {
	Proxy  string  `json:"proxy"`
	Voices []Voice `json:"voices"`
}

func NewListVoices(proxy string) (*ListVoices, error) {
	listVoices := &ListVoices{
		Proxy:  proxy,
		Voices: []Voice{},
	}

	// 创建一个新的 HTTP 客户端
	client := &http.Client{}

	// 创建一个新的 HTTP 请求
	req, err := http.NewRequest("GET", VOICE_LIST, nil)
	if err != nil {
		fmt.Println("创建请求失败:", err)
		return nil, err
	}

	// 设置请求头
	req.Header.Set("Authority", "speech.platform.bing.com")
	req.Header.Set("Sec-CH-UA", `" Not;A Brand";v="99", "Microsoft Edge";v="91", "Chromium";v="91"`)
	req.Header.Set("Sec-CH-UA-Mobile", "?0")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.77 Safari/537.36 Edg/91.0.864.41")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	// 如果需要使用代理，设置代理地址
	if proxy != "" {
		proxyURL, err := url.Parse(proxy)
		if err != nil {
			fmt.Println("解析代理URL失败:", err)
			return nil, err
		}
		client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyURL)}
	}

	// 发送请求并获取响应
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("请求失败:", err)
		return nil, err
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取内容失败:", err)
		return nil, err
	}

	json.Unmarshal(body, &listVoices.Voices)

	return listVoices, nil
}

// find Gender="Male/Female", Language="es"
func (l *ListVoices) Find(gender, language string) []Voice {
	if l == nil {
		return []Voice{}
	}

	gender = strings.Title(gender)
	voices := []Voice{}
	// 遍历原始字符串切片，进行过滤和提取
	for _, voice := range l.Voices {
		if len(gender) > 0 {
			if gender != voice.Gender {
				continue
			}
		}
		if len(language) > 0 {
			if !strings.HasPrefix(voice.Locale, language+"-") && voice.Locale != language {
				continue
			}
		}
		voices = append(voices, voice)
	}
	return voices
}
