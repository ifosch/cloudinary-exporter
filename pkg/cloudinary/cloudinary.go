package cloudinary

import (
	"errors"
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

func GetCredentials() (key, secret, cloud_name string, err error) {
	key = os.Getenv("CLOUDINARY_KEY")
	secret = os.Getenv("CLOUDINARY_SECRET")
	cloud_name = os.Getenv("CLOUDINARY_CLOUD_NAME")

	if key == "" || secret == "" || cloud_name == "" {
		err = errors.New("No credentials defined")
	}
	return key, secret, cloud_name, err
}
