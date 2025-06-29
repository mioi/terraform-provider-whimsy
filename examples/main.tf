terraform {
  required_providers {
    whimsy = {
      source  = "github.com/mioi/whimsy"
      version = "~> 1.0"
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

# Generate random animal name for server
resource "whimsy_animal" "server" {}

# Generate combined name for database using whimsy_name
resource "whimsy_name" "database" {
  parts     = ["color", "plant"]
  delimiter = "-"
}

# Example usage with AWS resources
resource "aws_instance" "web" {
  ami           = "ami-0abcdef1234567890"
  instance_type = "t2.micro"

  tags = {
    Name = "traefik-${resource.whimsy_animal.server.name}"
  }
}

resource "aws_db_instance" "database" {
  identifier = "data-${resource.whimsy_name.database.name}"
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
