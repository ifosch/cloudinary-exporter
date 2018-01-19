package cloudinary

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type UsageInfo struct {
	Usage       int64   `json:"usage"`
	Limit       int64   `json:"limit"`
	UsedPercent float64 `json:"used_percent"`
}

type UsageReport struct {
	Plan             string    `json:"plan"`
	LastUpdate       string    `json:"last_updated"`
	Transformations  UsageInfo `json:"transformations"`
	Objects          UsageInfo `json:"objects"`
	Bandwidth        UsageInfo `json:"bandwidth"`
	Storage          UsageInfo `json:"storage"`
	Requests         int64     `json:"requests"`
	Resources        int64     `json:"resources"`
	DerivedResources int64     `json:"derived_resources"`
}

func getCredentials() (key, secret, cloudName string, err error) {
	key = os.Getenv("CLOUDINARY_KEY")
	secret = os.Getenv("CLOUDINARY_SECRET")
	cloudName = os.Getenv("CLOUDINARY_CLOUD_NAME")

	if key == "" || secret == "" || cloudName == "" {
		err = errors.New("No credentials defined")
	}
	return key, secret, cloudName, err
}

func GetUsageReport() (usageReport *UsageReport, err error) {
	key, secret, cloudName, err := getCredentials()
	if err != nil {
		return nil, err
	}

	rs, err := http.Get(
		fmt.Sprintf(
			"https://%s:%s@api.cloudinary.com/v1_1/%s/usage",
			key,
			secret,
			cloudName,
		),
	)
	if err != nil {
		return nil, err
	}
	defer rs.Body.Close()

	bodyBytes, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		return nil, err
	}

	usageReport = new(UsageReport)
	err = json.Unmarshal(bodyBytes, &usageReport)
	return usageReport, err
}
