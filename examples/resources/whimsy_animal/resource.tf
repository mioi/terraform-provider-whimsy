# Basic animal resource
resource "whimsy_animal" "example" {}

# Animal with triggers
resource "whimsy_animal" "server" {
  triggers = {
    instance_id = "i-1234567890abcdef0"
  }
}

output "animal_name" {
  value = whimsy_animal.example.name
}

output "server_animal" {
  value = whimsy_animal.server.name
}