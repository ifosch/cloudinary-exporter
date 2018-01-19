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

func GetRequest() (req *http.Request, err error) {
	key, secret, cloudName, err := getCredentials()
	if err != nil {
		return nil, err
	}

	req, err = http.NewRequest(
		"GET",
		fmt.Sprintf(
			"https://%s:%s@api.cloudinary.com/v1_1/%s/usage",
			key,
			secret,
			cloudName,
		),
		nil,
	)
	return req, err
}

func GetUsageReport(req *http.Request) (usageReport *UsageReport, err error) {
	client := http.Client{}
	rs, err := client.Do(req)
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

func DerivedResources(usageReport UsageReport) float64 {
	return float64(usageReport.DerivedResources)
}

func Resources(usageReport UsageReport) float64 {
	return float64(usageReport.Resources)
}

func Requests(usageReport UsageReport) float64 {
	return float64(usageReport.Requests)
}

func StorageUsage(usageReport UsageReport) float64 {
	return float64(usageReport.Storage.Usage)
}

func StorageLimit(usageReport UsageReport) float64 {
	return float64(usageReport.Storage.Limit)
}

func StorageUsedPercent(usageReport UsageReport) float64 {
	return float64(usageReport.Storage.UsedPercent)
}

func BandwidthUsage(usageReport UsageReport) float64 {
	return float64(usageReport.Bandwidth.Usage)
}

func BandwidthLimit(usageReport UsageReport) float64 {
	return float64(usageReport.Bandwidth.Limit)
}

func BandwidthUsedPercent(usageReport UsageReport) float64 {
	return float64(usageReport.Bandwidth.UsedPercent)
}

func ObjectsUsage(usageReport UsageReport) float64 {
	return float64(usageReport.Objects.Usage)
}

func ObjectsLimit(usageReport UsageReport) float64 {
	return float64(usageReport.Objects.Limit)
}

func ObjectsUsedPercent(usageReport UsageReport) float64 {
	return float64(usageReport.Objects.UsedPercent)
}

func TransformationUsage(usageReport UsageReport) float64 {
	return float64(usageReport.Transformations.Usage)
}

func TransformationLimit(usageReport UsageReport) float64 {
	return float64(usageReport.Transformations.Limit)
}

func TransformationUsedPercent(usageReport UsageReport) float64 {
	return float64(usageReport.Transformations.UsedPercent)
}
