terraform {
  required_providers {
    whimsy = {
      source = "mioi/whimsy"
    }
  }
}

provider "whimsy" {}

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

resource "aws_rds_instance" "database" {
  identifier = local.database_name
  engine     = "mysql"
  # ... other configuration
}

# Advanced trigger usage examples:

# Regenerate name when any of multiple resources change
resource "whimsy_color" "multi_trigger" {
  triggers = {
    vpc_id    = aws_vpc.main.id
    subnet_id = aws_subnet.main.id
    sg_id     = aws_security_group.web.id
  }
}

# Regenerate based on resource attributes
resource "whimsy_plant" "version_trigger" {
  triggers = {
    app_version = var.app_version
    timestamp   = timestamp()
  }
}

# Force regeneration by changing trigger values
resource "whimsy_animal" "force_regen" {
  triggers = {
    force_new = "2024-01-01" # Change this value to force regeneration
  }
}
