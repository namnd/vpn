package models

type Node struct {
	ID           string
	Name         string
	Status       string
	InstanceType string
}

type Country struct {
	Name  string
	Nodes []Node
}

var CountriesInOrder = []string{"New Zealand", "Singapore", "Japan", "India", "US"}

var CountryRegion = map[string]string{
	"New Zealand": "ap-southeast-6",
	"Singapore":   "ap-southeast-1",
	"Japan":       "ap-northeast-1",
	"India":       "ap-south-1",
	"US":          "us-east-1",
}
