# Terraform Provider: Whimsy

[![Tests](https://github.com/mioi/terraform-provider-whimsy/actions/workflows/test.yml/badge.svg)](https://github.com/mioi/terraform-provider-whimsy/actions/workflows/test.yml)
[![Release](https://img.shields.io/github/release/mioi/terraform-provider-whimsy.svg)](https://github.com/mioi/terraform-provider-whimsy/releases)
[![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Terraform](https://img.shields.io/badge/terraform-1.0+-5C4EE5.svg)](https://www.terraform.io)

A Terraform provider that generates random yet memorable names with combinations of plants, animals, and colors. Perfect for creating human-friendly names for infrastructure resources while maintaining the "pets vs. cattle" principle.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.19

## Building the Provider

1. Clone the repository
2. Enter the repository directory
3. Build the provider using the Go `install` command:

```shell
go install
```

## Adding the Provider to Your Terraform Configuration

```hcl
terraform {
  required_providers {
    whimsy = {
      source = "mioi/whimsy"
    }
  }
}

provider "whimsy"
```

## Using the Provider

### Basic Usage

Generate composite memorable names for your infrastructure:

```hcl
# Generate random names for server
data "whimsy_color" "server_color" {}
data "whimsy_animal" "server_animal" {}

# Generate random names for database
data "whimsy_color" "database_color" {}
data "whimsy_plant" "database_plant" {}

# Local variables for constructed names
locals {
  server_name   = "traefik-${data.whimsy_color.server_color.name}-${data.whimsy_animal.server_animal.name}"
  database_name = "data-${data.whimsy_color.database_color.name}-${data.whimsy_plant.database_plant.name}"
}

# Use in your resources
resource "aws_instance" "web" {
  ami           = "ami-0abcdef1234567890"
  instance_type = "t2.micro"
  
  tags = {
    Name = local.server_name  # e.g., "traefik-blue-fox"
  }
}

resource "aws_rds_instance" "database" {
  identifier = local.database_name  # e.g., "data-red-oak"
  engine     = "mysql"
  # ... other configuration
}
```

### Using Triggers for Regeneration

Use the `triggers` attribute to regenerate names when specific resources or values change:

```hcl
# Regenerate name when the EC2 instance changes
data "whimsy_plant" "server_name" {
  triggers = {
    instance_id = aws_instance.web.id
  }
}

# Regenerate name when multiple resources change
data "whimsy_animal" "database_name" {
  triggers = {
    db_instance = aws_rds_instance.main.id
    environment = var.environment
    app_version = var.app_version
  }
}

# Force regeneration by changing trigger values
data "whimsy_color" "env_color" {
  triggers = {
    force_new = "2024-01-01"  # Change this to force new name
  }
}
```

The `triggers` attribute accepts a map of string values. When any value in the map changes between Terraform runs, the data source will generate a new random name. This is particularly useful when you want names to change along with infrastructure recreations.

## Data Sources

This provider includes three data sources:

- `whimsy_plant` - Generates random plant names (200+ options, max 6 chars)
- `whimsy_animal` - Generates random animal names (200+ options, max 6 chars)  
- `whimsy_color` - Generates random color names (200+ options, max 6 chars)

All names are lowercase, contain only English letters (a-z), and are designed to be memorable and pronounceable for infrastructure naming.

## Architecture

The provider uses a DRY (Don't Repeat Yourself) design with a single generic data source implementation that handles all three name types. This approach:

- Eliminates code duplication
- Ensures consistent behavior across all data sources
- Simplifies maintenance and testing
- Makes adding new name categories easy

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* This provider only generates names and doesn't create actual cloud resources, so tests are safe to run.

```shell
make testacc
```