package safeUrlChecker

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/driscollco-core/http-client"
	"github.com/driscollco-shortify/svc-shortcode/internal/conf"
	"strings"
)

type threatEntry struct {
	URL string `json:"url"`
}

type threatInfo struct {
	ThreatTypes      []string      `json:"threatTypes"`
	PlatformTypes    []string      `json:"platformTypes"`
	ThreatEntryTypes []string      `json:"threatEntryTypes"`
	ThreatEntries    []threatEntry `json:"threatEntries"`
}

type apiRequest struct {
	Client     interface{} `json:"client"`
	ThreatInfo threatInfo  `json:"threatInfo"`
}

type apiResponse struct {
	Matches []struct {
		ThreatType      string `json:"threatType"`
		PlatformType    string `json:"platformType"`
		ThreatEntryType string `json:"threatEntryType"`
		Threat          struct {
			URL string `json:"url"`
		} `json:"threat"`
	} `json:"matches"`
}

func IsUnsafe(url string) (bool, error) {
	reqData := apiRequest{
		Client: struct {
			ClientID      string `json:"clientId"`
			ClientVersion string `json:"clientVersion"`
		}{
			ClientID:      "shortify.pro", // Replace with your app's unique identifier
			ClientVersion: "1.0.0",        // Replace with your app version
		},
		ThreatInfo: threatInfo{
			ThreatTypes:      []string{"MALWARE", "SOCIAL_ENGINEERING", "UNWANTED_SOFTWARE", "POTENTIALLY_HARMFUL_APPLICATION"},
			PlatformTypes:    []string{"ANY_PLATFORM"},
			ThreatEntryTypes: []string{"URL", "EXECUTABLE", "IP_RANGE"},
			ThreatEntries:    []threatEntry{{URL: url}},
		},
	}

	req := http.NewRequest()
	req.Body(reqData)
	resp, err := req.Do(fmt.Sprintf("%s=%s", "https://safebrowsing.googleapis.com/v4/threatMatches:find?key", conf.Config.GCP.SafeSite.ApiKey), "POST")
	if err != nil {
		return false, err
	}

	if strings.Contains(string(resp.Content()), `"error":`) {
		return false, errors.New(string(resp.Content()))
	}

	var resData apiResponse
	if err = json.Unmarshal(resp.Content(), &resData); err != nil {
		return false, err
	}

	return len(resData.Matches) > 0, nil
}
