# Basic color resource
resource "whimsy_color" "example" {}

# Color with triggers
resource "whimsy_color" "server" {
  triggers = {
    environment = "production"
  }
}

output "color_name" {
  value = whimsy_color.example.name
}

output "server_color" {
  value = whimsy_color.server.name
}