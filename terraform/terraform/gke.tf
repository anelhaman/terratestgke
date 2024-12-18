provider "google" {
  project = var.google_project
}

resource "google_service_account" "default" {
  account_id   = "service-account-id"
  display_name = "Service Account"
  project = var.google_project
}

resource "google_container_cluster" "primary" {

  name     = var.gke-name
  location = var.google_region
  project = var.google_project  

  # # We can't create a cluster with no node pool defined, but we want to only use
  # # separately managed node pools. So we create the smallest possible default
  # # node pool and immediately delete it.
  remove_default_node_pool = true
  initial_node_count = 1
  
}

resource "google_container_node_pool" "primary_preemptible_nodes" {
  name       = "my-node-pool"
  location   = var.google_region
  cluster    = google_container_cluster.primary.name
  project    = var.google_project
  node_count = 1

  node_config {
    preemptible  = true
    machine_type = "e2-medium"
    disk_size_gb = 20

    # Google recommends custom service accounts that have cloud-platform scope and permissions granted via IAM Roles.
    service_account = google_service_account.default.email
    oauth_scopes    = [
      "https://www.googleapis.com/auth/cloud-platform"
    ]
  }
}

