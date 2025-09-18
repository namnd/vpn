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

type CreateNodeInput struct {
	CountryName string `form:"country"`
}

func CreateNode(c *gin.Context) {
	var formInput CreateNodeInput
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

	input := &ec2.RunInstancesInput{
		MaxCount:     aws.Int32(1),
		MinCount:     aws.Int32(1),
		ImageId:      aws.String("ami-07dbf7fde6187421a"),
		InstanceType: types.InstanceType("t4g.small"),
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
	slog.Info("Instance created successfully", "id", instanceID)
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
