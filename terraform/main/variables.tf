# AWS variables
variable "aws_region" { 
  type = string 
  default = "us-east-1" 
}

variable "vpc_cidr"   { 
  type = string
  default = "10.0.0.0/16" 
}



# ECR variables
variable "ecr_api_image" {
  type    = string
}

variable "ecr_app_image" {
  type    = string
}


# Database variables
variable "db_host" {
  type    = string
}

variable "db_name" {
  type    = string
}

variable "db_user" {
  type    = string
}

variable "db_password" {
  type    = string
}

variable "db_port" {
  type    = string
}

variable "db_ssl" {
  type    = string
  default = "true"
}



