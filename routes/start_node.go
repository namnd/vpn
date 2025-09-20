package routes

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/gin-gonic/gin"
	"github.com/namnd/vpn/models"
)

func StartNode(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		slog.Error("invalid instanceID")
	}

	var formInput NodeFormInput
	err := c.ShouldBind(&formInput)
	if err != nil {
		slog.Error("invalid input", "err", err)
	}

	region := models.CountryRegion[formInput.CountryName]
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		slog.Error("failed to load AWS config", "error", err)
	}

	client := ec2.NewFromConfig(cfg)

	result, err := client.StartInstances(context.TODO(), &ec2.StartInstancesInput{
		InstanceIds: []string{id},
	})
	if err != nil {
		slog.Error("failed to start instance", "error", err)
	}

	fmt.Println(result)

}
