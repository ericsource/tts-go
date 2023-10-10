package azure_tts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// voiceListAPI is the source for supported voice list to region mapping
// See: https://docs.microsoft.com/en-us/azure/cognitive-services/speech-service/rest-text-to-speech#regions-and-endpoints
const voiceListAPI = "https://%s.tts.speech.microsoft.com/cognitiveservices/voices/list"

//go:generate enumer -type=voiceType -linecomment -json
type voiceType int

const (
	voiceStandard voiceType = iota // Standard
	voiceNeural                    // Neural
)

type regionVoiceListResponse struct {
	Name            string    `json:"Name"`
	ShortName       string    `json:"ShortName"`
	Gender          Gender    `json:"Gender"`
	Locale          Locale    `json:"Locale"`
	SampleRateHertz string    `json:"SampleRateHertz"`
	VoiceType       voiceType `json:"VoiceType"`
}

// supportedVoices represents the key used within the `localeToGender` map.
type supportedVoices struct {
	Gender Gender
	Locale Locale
}

type RegionVoiceMap map[supportedVoices]string

func (az *AzureCSTextToSpeech) buildVoiceToRegionMap() (RegionVoiceMap, error) {

	v, err := az.fetchVoiceList()
	if err != nil {
		return nil, err
	}

	m := make(map[supportedVoices]string)
	for _, x := range v {
		//if x.VoiceType == voiceStandard {
		//	m[supportedVoices{Gender: x.Gender, Locale: x.Locale}] = x.ShortName
		//}
		if x.VoiceType == voiceNeural {
			m[supportedVoices{Gender: x.Gender, Locale: x.Locale}] = x.ShortName
		}
	}

	return m, err
}

func (az *AzureCSTextToSpeech) fetchVoiceList() ([]regionVoiceListResponse, error) {

	request, err := http.NewRequest(http.MethodGet, az.voiceServiceListURL, nil)
	if err != nil {
		return nil, err
	}

	//request.Header.Set("Authorization", "Bearer "+az.accessToken)
	request.Header.Set("Ocp-Apim-Subscription-Key", az.SubscriptionKey)
	client := &http.Client{Timeout: 20 * time.Second}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	switch response.StatusCode {
	case http.StatusOK:
		var r []regionVoiceListResponse

		var buf bytes.Buffer
		_, err = io.Copy(&buf, response.Body)
		if err != nil {
			//fmt.Println("Failed to copy response body:", err)
			return nil, fmt.Errorf("Failed to copy response body: %v", err)
		}
		json.Unmarshal(buf.Bytes(), &r)

		return r, nil
	case http.StatusBadRequest:
		return nil, fmt.Errorf("%d - A required parameter is missing, empty, or null. Or, the value passed to either a required or optional parameter is invalid. A common issue is a header that is too long", response.StatusCode)
	case http.StatusUnauthorized:
		return nil, fmt.Errorf("%d - The request is not authorized. Check to make sure your subscription key or token is valid and in the correct region", response.StatusCode)
	case http.StatusTooManyRequests:
		return nil, fmt.Errorf("%d - You have exceeded the quota or rate of requests allowed for your subscription", response.StatusCode)
	case http.StatusBadGateway:
		return nil, fmt.Errorf("%d - Network or server-side issue. May also indicate invalid headers", response.StatusCode)
	}
	return nil, fmt.Errorf("%d - unexpected response code from voice list API", response.StatusCode)
}
