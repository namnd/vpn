terraform {
  backend "s3" {
    bucket       = "namnd-vpn-infra"
    key          = "terraform.tfstate"
    region       = "ap-southeast-2"
    use_lockfile = true
  }
}

provider "aws" {
  region = "ap-southeast-2"
}

data "aws_iam_policy_document" "create_ec2_policy_document" {
  statement {
    actions = [
      "ec2:CreateTags",
      "ec2:DescribeImages",
      "ec2:DescribeInstances",
      "ec2:RunInstances",
      "ec2:StartInstances",
      "ec2:StopInstances",
      "ec2:TerminateInstances",
    ]

    resources = ["*"]
  }
}

resource "aws_iam_policy" "create_ec2_policy" {
  name        = "ec2-create-policy"
  description = "Allows creating EC2 instances"
  policy      = data.aws_iam_policy_document.create_ec2_policy_document.json
}

