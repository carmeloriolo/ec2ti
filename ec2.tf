locals {
  instance_count = 10
}

resource "random_string" "random" {
  count         = local.instance_count
  length        = 8
  special       = false
}

resource "aws_instance" "ec2ti" {
  count         = local.instance_count
  ami           = "ami-0d57c0143330e1fa7"
  instance_type = "t2.micro"

  tags = {
    Name = format("%s", random_string.random[count.index].result) 
  }
}

provider "aws" {
  region                      = "eu-west-1"
  skip_credentials_validation = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true

  endpoints {
    ec2 = "http://localstack:4566"
  }
}
