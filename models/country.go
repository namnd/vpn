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

var CountriesInOrder = []string{"Singapore", "Japan", "India"}

var CountryRegion = map[string]string{
	"Singapore": "ap-southeast-1",
	"Japan":     "ap-northeast-1",
	"India":     "ap-south-1",
}
