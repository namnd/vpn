package routes

import (
	"context"
	"fmt"
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
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		slog.Error("failed to load AWS config", "error", err)
	}

	input := &ec2.DescribeInstancesInput{
		MaxResults: aws.Int32(5),
	}

	var countries []models.Country

	for name, region := range models.CountryRegion {
		regionCfg := cfg.Copy()
		regionCfg.Region = region
		regionClient := ec2.NewFromConfig(regionCfg)

		result, err := regionClient.DescribeInstances(context.TODO(), input)

		if err != nil {
			slog.Error("failed to list instances", "error", err)
		}

		var nodes []models.Node
		for _, v := range result.Reservations {
			for _, i := range v.Instances {
				fmt.Println(aws.ToString(i.InstanceId))
				nodes = append(nodes, models.Node{
					Name:   aws.ToString(i.InstanceId),
					Status: string(i.State.Name),
				})
			}
		}

		countries = append(countries, models.Country{
			Name:  name,
			Flag:  models.CountryFlags[name],
			Nodes: nodes,
		})

	}

	p := gintemplrenderer.New(
		c.Request.Context(),
		http.StatusOK,
		ui.Home(countries),
	)

	c.Render(http.StatusOK, p)
}
