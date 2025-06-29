terraform {
  required_providers {
    whimsy = {
      source = "mioi/whimsy"
    }
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "whimsy" {}

provider "aws" {
  region = "us-west-2"
}

# Generate random names for server
resource "whimsy_color" "server_color" {}
resource "whimsy_animal" "server_animal" {}

# Generate random names for database
resource "whimsy_color" "database_color" {}
resource "whimsy_plant" "database_plant" {}

# Local variables for constructed names
locals {
  server_name   = "traefik-${resource.whimsy_color.server_color.name}-${resource.whimsy_animal.server_animal.name}"
  database_name = "data-${resource.whimsy_color.database_color.name}-${resource.whimsy_plant.database_plant.name}"
}

# Example usage with AWS resources
resource "aws_instance" "web" {
  ami           = "ami-0abcdef1234567890"
  instance_type = "t2.micro"

  tags = {
    Name = local.server_name
  }
}

resource "aws_db_instance" "database" {
  identifier = local.database_name
  engine     = "mysql"

  instance_class    = "db.t3.micro"
  allocated_storage = 20
  db_name           = "whimsy"
  username          = "admin"
  password          = "changeme123"

  skip_final_snapshot = true
}

# Trigger example: regenerate when instance changes
resource "whimsy_color" "server_trigger" {
  triggers = {
    instance_id = aws_instance.web.id
  }
}
