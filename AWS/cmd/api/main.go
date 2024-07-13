package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func main() {
	var (
		instanceID string
		err        error
	)

	ctx := context.Background()

	if instanceID, err = createEC2(ctx, "eu-north-1", "aws-demo-key0", "ubuntu/images/hvm-ssd-gp3/ubuntu-noble-24.04-amd64-server-20240701.1", types.InstanceTypeT3Micro); err != nil {
		log.Fatalf("createEC2 error: %s", err)
	}

	fmt.Printf("Instance id: %s\n", instanceID)
}

func createEC2(ctx context.Context, region, keyName, imageName string, instanceType types.InstanceType) (string, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return "", fmt.Errorf("unable to load SDK config, %s", err)
	}

	ec2Client := ec2.NewFromConfig(cfg)

	_, err = ec2Client.CreateKeyPair(ctx, &ec2.CreateKeyPairInput{
		KeyName: aws.String(keyName),
	})

	if err != nil {
		return "", fmt.Errorf("CreateKeyPair error: %s", err)
	}

	vpcOutput, err := ec2Client.CreateVpc(ctx, &ec2.CreateVpcInput{
		CidrBlock: aws.String("10.0.0.0/16"),
	})
	if err != nil {
		return "", fmt.Errorf("CreateVpc error: %s", err)
	}

	subnetOutput, err := ec2Client.CreateSubnet(ctx, &ec2.CreateSubnetInput{
		VpcId:     vpcOutput.Vpc.VpcId,
		CidrBlock: aws.String("10.0.1.0/24"),
	})
	if err != nil {
		return "", fmt.Errorf("CreateSubnet error: %s", err)
	}

	imageOutput, err := ec2Client.DescribeImages(ctx, &ec2.DescribeImagesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("name"),
				Values: []string{imageName},
			},
			{
				Name:   aws.String("virtualization-type"),
				Values: []string{"hvm"},
			},
		},
		Owners: []string{"099720109477"}, //from OS owner
	})

	if err != nil {
		return "", fmt.Errorf("DescribeImages error: %s", err)
	}
	if len(imageOutput.Images) == 0 {
		return "", fmt.Errorf("imageOutput.Images is of 0 length")
	}

	// Check if instance type is supported in the region
	instanceTypesOutput, err := ec2Client.DescribeInstanceTypeOfferings(ctx, &ec2.DescribeInstanceTypeOfferingsInput{
		LocationType: types.LocationTypeRegion,
		Filters: []types.Filter{
			{
				Name:   aws.String("instance-type"),
				Values: []string{string(instanceType)},
			},
			{
				Name:   aws.String("location"),
				Values: []string{region},
			},
		},
	})
	if err != nil {
		return "", fmt.Errorf("DescribeInstanceTypeOfferings error: %s", err)
	}
	if len(instanceTypesOutput.InstanceTypeOfferings) == 0 {
		return "", fmt.Errorf("instance type %s is not available in region %s", instanceType, region)
	}

	runResult, err := ec2Client.RunInstances(ctx, &ec2.RunInstancesInput{
		ImageId:      imageOutput.Images[0].ImageId,
		KeyName:      aws.String(keyName),
		InstanceType: instanceType,
		MinCount:     aws.Int32(1),
		MaxCount:     aws.Int32(1),
		SubnetId:     subnetOutput.Subnet.SubnetId,
	})

	if err != nil {
		return "", fmt.Errorf("RunInstances error: %s", err)
	}

	if len(runResult.Instances) == 0 {
		return "", fmt.Errorf("runResult.Instances is of 0 length")
	}

	return *runResult.Instances[0].InstanceId, nil
}
