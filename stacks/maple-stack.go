package stacks

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsecr"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsecs"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type CloudformationStackProps struct {
	awscdk.StackProps
}

func MapleStackFormation(scope constructs.Construct, id string, props *CloudformationStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// vpc for ecs cluster
	vpc := awsec2.NewVpc(stack, jsii.String("maple-vpc"), &awsec2.VpcProps{
		MaxAzs:      jsii.Number(2), // Use 2 availability zones
		NatGateways: jsii.Number(1), // Use 1 NAT gateway to save costs
	})

	// ecr repo to store docker images
	ecrRepo := awsecr.NewRepository(stack, jsii.String("maple-ecr-repo"), &awsecr.RepositoryProps{
		RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
	})

	// ecs cluster, which is just the grouping of the ecs services
	ecsCluster := awsecs.NewCluster(stack, jsii.String("maple-ecs-cluster"), &awsecs.ClusterProps{
		Vpc: vpc,
	})

	// ecs task definition
	ecsTaskDef := awsecs.NewTaskDefinition(stack, jsii.String("maple-ecs-task-def"), &awsecs.TaskDefinitionProps{
		Compatibility: awsecs.Compatibility_FARGATE,
		Cpu:           jsii.String("256"),
		MemoryMiB:     jsii.String("512"),
	})

	// add container to task definition
	ecsTaskDef.AddContainer(jsii.String("maple-container"), &awsecs.ContainerDefinitionOptions{
		Image:          awsecs.ContainerImage_FromEcrRepository(ecrRepo, jsii.String("latest")),
		Essential:      jsii.Bool(true),
		MemoryLimitMiB: jsii.Number(512),
		Logging: awsecs.LogDrivers_AwsLogs(&awsecs.AwsLogDriverProps{
			StreamPrefix: jsii.String("maple"),
		}),
	})

	// ecs service, what actually runs the containers
	ecsService := awsecs.NewFargateService(stack, jsii.String("maple-ecs-service"), &awsecs.FargateServiceProps{
		Cluster:        ecsCluster,
		TaskDefinition: ecsTaskDef,
	})
	_ = ecsService

	return stack
}
