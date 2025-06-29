# Terraform Provider: Whimsy

[![Tests](https://github.com/mioi/terraform-provider-whimsy/actions/workflows/test.yml/badge.svg)](https://github.com/mioi/terraform-provider-whimsy/actions/workflows/test.yml)
[![Release](https://img.shields.io/github/release/mioi/terraform-provider-whimsy.svg)](https://github.com/mioi/terraform-provider-whimsy/releases)
[![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A Terraform provider that generates random yet memorable names with combinations of plants, animals, and colors. Perfect for creating human-friendly names for infrastructure resources while maintaining the "pets vs. cattle" principle.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.21

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
      source  = "github.com/mioi/whimsy"
      version = "~> 1.0"
    }
  }
}

provider "whimsy" {}
```

## Using the Provider

### Basic Usage

Generate composite memorable names for your infrastructure:

```hcl
# Generate individual names for server and database
resource "whimsy_animal" "server" {}

resource "whimsy_name" "database" {
  parts     = ["color", "plant"]
  delimiter = "-"
}

# Use in your resources
resource "aws_instance" "web" {
  ami           = "ami-0abcdef1234567890"
  instance_type = "t2.micro"
  
  tags = {
    Name = "traefik-${resource.whimsy_animal.server.name}"  # e.g., "traefik-fox"
  }
}

resource "aws_rds_instance" "database" {
  identifier = "data-${resource.whimsy_name.database.name}"  # e.g., "data-red-oak"
  engine     = "mysql"
  # ... other configuration
}
```

### Using Triggers for Regeneration

Use the `triggers` attribute to regenerate names when specific resources or values change:

```hcl
# Regenerate name when the EC2 instance changes
resource "whimsy_plant" "server_name" {
  triggers = {
    instance_id = aws_instance.web.id
  }
}

# Regenerate name when multiple resources change
resource "whimsy_animal" "database_name" {
  triggers = {
    db_instance = aws_rds_instance.main.id
    environment = var.environment
    app_version = var.app_version
  }
}

# Force regeneration by changing trigger values
resource "whimsy_color" "env_color" {
  triggers = {
    force_new = "2024-01-01"  # Change this to force new name
  }
}
```

The `triggers` attribute accepts a map of string values. When any value in the map changes between Terraform runs, the resource will generate a new random name. This is particularly useful when you want names to change along with infrastructure recreations.

## Resources

This provider includes four resources:

- `whimsy_plant` - Generates random plant names (200+ options, max 6 chars)
- `whimsy_animal` - Generates random animal names (200+ options, max 6 chars)  
- `whimsy_color` - Generates random color names (200+ options, max 6 chars)
- `whimsy_name` - Combines multiple parts with configurable delimiter and order

### whimsy_name Resource

The `whimsy_name` resource allows you to generate combined names from multiple categories:

**Arguments:**
- `parts` - List of name parts to combine: `["plant", "animal", "color"]` (default: `["color", "animal"]`)
- `delimiter` - String to separate parts (default: `"-"`)
- `random` - Boolean to randomize part order (default: `false`)
- `triggers` - Map of values that trigger regeneration when changed

**Example:**
```hcl
resource "whimsy_name" "server" {
  parts     = ["color", "plant"]
  delimiter = "-"
}
# Generates names like: "blue-oak", "red-elm", "gold-ivy"
```

All names are lowercase, contain only English letters (a-z), and are designed to be memorable and pronounceable for infrastructure naming. Resources persist their generated names in Terraform state and only regenerate when triggers change.

## Architecture

The provider uses a DRY (Don't Repeat Yourself) design with a single generic resource implementation that handles all three name types. This approach:

- Eliminates code duplication
- Ensures consistent behavior across all resources
- Simplifies maintenance and testing
- Makes adding new name categories easy

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go build .` This will build the provider binary in the current directory.

To run tests, use:

```shell
go test ./...
```

*Note:* This provider only generates names and doesn't interact with external services, so tests are safe to run and complete quickly.