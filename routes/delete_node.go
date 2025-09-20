package routes

import (
	"context"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/gin-gonic/gin"
	"github.com/namnd/vpn/models"
)

func DeleteNode(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		slog.Error("invalid instanceID")
	}

	country := c.Query("country")

	region := models.CountryRegion[country]
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		slog.Error("failed to load AWS config", "error", err)
	}

	client := ec2.NewFromConfig(cfg)

	_, err = client.TerminateInstances(context.TODO(), &ec2.TerminateInstancesInput{
		InstanceIds: []string{id},
	})
	if err != nil {
		slog.Error("failed to delete instance", "error", err)
	}

	slog.Info("deleted instance successfully", "region", region, "id", id)
}
