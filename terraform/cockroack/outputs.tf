
output "connection_string" {
  value     = data.cockroach_connection_string.db.connection_string
  sensitive = true
}