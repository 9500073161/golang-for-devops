To connect AWS cloud.

1.open aws account
2.create IAM user
3.Download credentials
4.configure credentails
5.AWS CLI and Environment varriable (Visual Studio) >>

vcode>> RUN >>> >>Add configuration 


"env":{
                "AWS_ACCESS_KEY_ID":"...",
                "AWS_SECRET_ACCESS_KEY": "..",
                "AWS_DEFAULT_REGION": "eu-north-1"

            }

6.check your identity 

aws sts get-caller-identity

7. AWS SDK

 1.Load config
 2.Create SSH keyPair
 3.Find Ubuntu AMI
 4.Launch Ec2
 5.Output Instance ID

 8. I have created Ubuntu VM on my region, if you want to change your instance details. PLease go to below location.

 if instanceID, err = createEC2(ctx, "eu-north-1", "aws-demo-key2", "ubuntu/images/hvm-ssd-gp3/ubuntu-noble-24.04-amd64-server-20240701.1", types.InstanceTypeT3Micro); err != nil {
		log.Fatalf("createEC2 error: %s", err)
	}






