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

type ReportDesc struct {
	Name string
	Desc string
}

var ReportDescs = []ReportDesc{
	ReportDesc{Name: "transformations_usage", Desc: "Transformation usage"},
	ReportDesc{Name: "transformations_limit", Desc: "Transformation limit"},
	ReportDesc{Name: "transformations_used_percent", Desc: "Transformation used percent"},
	ReportDesc{Name: "objects_usage", Desc: "Object usage"},
	ReportDesc{Name: "objects_limit", Desc: "Object limit"},
	ReportDesc{Name: "objects_used_percent", Desc: "Object used percent"},
	ReportDesc{Name: "bandwidth_usage", Desc: "Bandwidth usage"},
	ReportDesc{Name: "bandwidth_limit", Desc: "Bandwidth limit"},
	ReportDesc{Name: "bandwidth_used_percent", Desc: "Bandwidth used percent"},
	ReportDesc{Name: "storage_usage", Desc: "Storage usage"},
	ReportDesc{Name: "storage_limit", Desc: "Storage limit"},
	ReportDesc{Name: "storage_used_percent", Desc: "Storage used percent"},
	ReportDesc{Name: "requests", Desc: "Requests"},
	ReportDesc{Name: "resources", Desc: "Resources"},
	ReportDesc{Name: "derived_resources", Desc: "Derived resources"},
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
