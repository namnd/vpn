package routes

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/base64"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"text/template"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/gin-gonic/gin"
	"github.com/namnd/vpn/models"
)

//go:embed user_data.sh
var userDataTemplate string

type UserDataTemplateData struct {
	HostName         string
	TailscaleAuthKey string
}

type NodeFormInput struct {
	CountryName  string `form:"country"`
	InstanceType string `form:"instance_type"`
}

func CreateNode(c *gin.Context) {
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

	tmpl, err := template.New("userData").Parse(userDataTemplate)
	if err != nil {
		slog.Error("failed to parse user_data.sh", "error", err)
	}
	var userData bytes.Buffer
	hostName := fmt.Sprintf("%s-%s", formInput.CountryName, randomSuffix(4))
	if err := tmpl.Execute(&userData, UserDataTemplateData{
		HostName:         hostName,
		TailscaleAuthKey: os.Getenv("TAILSCALE_AUTH_KEY"),
	}); err != nil {
		slog.Error("failed to execute userData template", "error", err)
	}

	imageID, err := getLatestAMI(client)
	if err != nil {
		slog.Error("failed to get latest AMI", "error", err)
	}

	input := &ec2.RunInstancesInput{
		MaxCount:     aws.Int32(1),
		MinCount:     aws.Int32(1),
		ImageId:      aws.String(imageID),
		InstanceType: types.InstanceType(formInput.InstanceType),
		UserData:     aws.String(base64.StdEncoding.EncodeToString(userData.Bytes())),
		TagSpecifications: []types.TagSpecification{
			{
				ResourceType: types.ResourceTypeInstance,
				Tags: []types.Tag{
					{
						Key:   aws.String("Name"),
						Value: aws.String(hostName),
					},
				},
			},
		},
	}

	result, err := client.RunInstances(context.TODO(), input)
	if err != nil {
		slog.Error("failed to start instance", "error", err)
	}

	instanceID := *result.Instances[0].InstanceId
	slog.Info("created instance successfully", "region", region, "id", instanceID)
}

func randomSuffix(length int) string {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rng.Intn(len(charset))]
	}
	return string(result)
}

func getLatestAMI(client *ec2.Client) (string, error) {
	input := &ec2.DescribeImagesInput{
		Owners: []string{"amazon"},
		Filters: []types.Filter{
			{
				Name:   aws.String("name"),
				Values: []string{"al2023-ami-2023*"},
			},
			{
				Name:   aws.String("architecture"),
				Values: []string{"arm64"},
			},
		},
	}

	result, err := client.DescribeImages(context.TODO(), input)
	if err != nil {
		return "", fmt.Errorf("failed to get AMI: %v", err)
	}

	if len(result.Images) == 0 {
		return "", fmt.Errorf("no matching AMI found")
	}

	latestAMI := result.Images[0]
	for _, img := range result.Images[1:] {
		if img.CreationDate != nil && latestAMI.CreationDate != nil {
			if *img.CreationDate > *latestAMI.CreationDate {
				latestAMI = img
			}
		}
	}
	return *latestAMI.ImageId, nil
}
