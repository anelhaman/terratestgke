output "cluster_name" {
  value = google_container_cluster.primary.name
}
output "region" {
  value = var.google_region
}
output "project" {
  value = var.google_project
}