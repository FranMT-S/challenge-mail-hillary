# Prerequisites

1. install terraform
2. install aws cli and configure the credentials, the role must have the permissions to create the cluster and the task definitions, must be can create awslogs-group and pull the images from ecr


# Deploy Client and Api



1. create repository

``` bash
aws ecr create-repository --repository-name challengehillaryclintonemails-api --region us-east-1

aws ecr create-repository --repository-name challengehillaryclintonemails-app --region us-east-1
```

2. create the images with the next command and the root of the project

``` bash
docker compose up 
```

3. tag the images

replace the values with the uri of the repository, the structure is `aws_account_id.dkr.ecr.region.amazonaws.com/repository_name:tag`

``` bash
docker tag challengehillaryclintonemails-api:latest aws_account_id.dkr.ecr.region.amazonaws.com/challengehillaryclintonemails-api:latest

docker tag challengehillaryclintonemails-app:latest aws_account_id.dkr.ecr.region.amazonaws.com/challengehillaryclintonemails-app:latest
```

4. push the images

get permissions to push the images with the next command, replace the values region and aws_account_id with the values of your environment

``` bash
aws ecr get-login-password --region region | docker login --username AWS --password-stdin aws_account_id.dkr.ecr.region.amazonaws.com
```

after get the permissions push the images
``` bash
docker push aws_account_id.dkr.ecr.region.amazonaws.com/challengehillaryclintonemails-api:latest

docker push aws_account_id.dkr.ecr.region.amazonaws.com/challengehillaryclintonemails-app:latest
```

You can read more about [here](https://docs.aws.amazon.com/AmazonECR/latest/userguide/docker-push-ecr-image.html)

## Terraform Client and Api

1. navigate to the folder main

``` bash
cd main
```

2. copy the template.tfvars in a new file terraform.tfvars and fill the values


2. run terraform

``` bash
terraform init
terraform plan
terraform apply
```

# Cockroachdb

If you want, you can use the terraform module to create a cluster in cockroachdb

## Prerequisites

We need to get a api key from cockroachdb

You need go to [Acces Management](https://cockroachlabs.cloud/access), click the tab Service Account, then create the a new service account, then we need to make sure to have the roles to create and admin cluster in the actions options.

You can read the oficial documentation [here](https://www.cockroachlabs.com/docs/stable/cockroachcloud-get-started.html).

## Terraform

1. navigate to the folder cockroachdb

``` bash
cd cockroachdb
```

2. copy the template.tfvars in a new file terraform.tfvars and fill the values

``` bash
cp template.tfvars tfvars
terraform init
terraform apply
```