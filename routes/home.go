package routes

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/a-h/templ/examples/integration-gin/gintemplrenderer"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/gin-gonic/gin"
	"github.com/namnd/vpn/models"
	"github.com/namnd/vpn/ui"
)

func Home(c *gin.Context) {
	country := c.Param("country")
	if country == "" {
		country = models.CountriesInOrder[0]
	}

	region, found := models.CountryRegion[country]
	if !found {
		slog.Error("invalid country")
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		slog.Error("failed to load AWS config", "error", err)
	}

	client := ec2.NewFromConfig(cfg)
	input := &ec2.DescribeInstancesInput{
		MaxResults: aws.Int32(5),
	}

	result, err := client.DescribeInstances(context.TODO(), input)
	if err != nil {
		slog.Error("failed to list instances", "error", err)
	}

	var nodes []models.Node
	for _, v := range result.Reservations {
		for _, i := range v.Instances {
			var name string
			for _, v := range i.Tags {
				if *v.Key == "Name" {
					name = *v.Value
				}
			}
			nodes = append(nodes, models.Node{
				ID:     aws.ToString(i.InstanceId),
				Name:   name,
				Status: string(i.State.Name),
			})
		}
	}

	p := gintemplrenderer.New(
		c.Request.Context(),
		http.StatusOK,
		ui.Home(models.Country{
			Name:  country,
			Nodes: nodes,
		}),
	)

	c.Render(http.StatusOK, p)
}
