# Basic usage with defaults (color-animal with "-" delimiter)
resource "whimsy_name" "default" {}

# Custom parts and delimiter
resource "whimsy_name" "server" {
  parts     = ["color", "plant"]
  delimiter = "_"
}

# All three parts with custom delimiter
resource "whimsy_name" "service" {
  parts     = ["animal", "color", "plant"]
  delimiter = "."
}

# Random order with triggers
resource "whimsy_name" "random_order" {
  parts  = ["plant", "animal", "color"]
  random = true
  triggers = {
    version = "1.0.0"
  }
}

# Single part
resource "whimsy_name" "simple" {
  parts = ["color"]
}

# Use in other resources
locals {
  instance_name = "web-${whimsy_name.server.name}"
}

output "default_name" {
  value = whimsy_name.default.name
}

output "server_name" {
  value = whimsy_name.server.name
}

output "service_name" {
  value = whimsy_name.service.name
}

output "random_name" {
  value = whimsy_name.random_order.name
}

output "simple_name" {
  value = whimsy_name.simple.name
}

output "instance_name" {
  value = local.instance_name
}