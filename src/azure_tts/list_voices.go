package azure_tts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Voice struct {
	Name            string `json:"Name"`
	ShortName       string `json:"ShortName"`
	Gender          string `json:"Gender"`
	Locale          string `json:"Locale"`
	SampleRateHertz string `json:"SampleRateHertz"`
	VoiceType       string `json:"VoiceType"`
}

func (az *AzureCSTextToSpeech) FetchVoiceList() ([]Voice, error) {

	request, err := http.NewRequest(http.MethodGet, az.voiceServiceListURL, nil)
	if err != nil {
		return []Voice{}, err
	}

	//request.Header.Set("Authorization", "Bearer "+az.accessToken)
	request.Header.Set("Ocp-Apim-Subscription-Key", az.SubscriptionKey)
	client := &http.Client{Timeout: 20 * time.Second}
	response, err := client.Do(request)
	if err != nil {
		return []Voice{}, err
	}
	defer response.Body.Close()

	switch response.StatusCode {
	case http.StatusOK:
		var voices []Voice

		var buf bytes.Buffer
		_, err = io.Copy(&buf, response.Body)
		if err != nil {
			//fmt.Println("Failed to copy response body:", err)
			return []Voice{}, fmt.Errorf("Failed to copy response body: %v", err)
		}
		json.Unmarshal(buf.Bytes(), &voices)

		az.Voices = voices

		return voices, nil
	case http.StatusBadRequest:
		return []Voice{}, fmt.Errorf("%d - A required parameter is missing, empty, or null. Or, the value passed to either a required or optional parameter is invalid. A common issue is a header that is too long", response.StatusCode)
	case http.StatusUnauthorized:
		return []Voice{}, fmt.Errorf("%d - The request is not authorized. Check to make sure your subscription key or token is valid and in the correct region", response.StatusCode)
	case http.StatusTooManyRequests:
		return []Voice{}, fmt.Errorf("%d - You have exceeded the quota or rate of requests allowed for your subscription", response.StatusCode)
	case http.StatusBadGateway:
		return []Voice{}, fmt.Errorf("%d - Network or server-side issue. May also indicate invalid headers", response.StatusCode)
	}
	return []Voice{}, fmt.Errorf("%d - unexpected response code from voice list API", response.StatusCode)
}

// find Gender="Male/Female", Language="es"
func (az *AzureCSTextToSpeech) Find(gender, language string) []Voice {
	if az == nil {
		return []Voice{}
	}

	gender = strings.Title(gender)
	voices := []Voice{}
	// 遍历原始字符串切片，进行过滤和提取
	for _, voice := range az.Voices {
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
