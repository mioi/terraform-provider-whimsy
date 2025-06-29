# Basic plant resource
resource "whimsy_plant" "example" {}

# Plant with triggers
resource "whimsy_plant" "database" {
  triggers = {
    version = "1.0.0"
  }
}

output "plant_name" {
  value = whimsy_plant.example.name
}

output "database_plant" {
  value = whimsy_plant.database.name
}