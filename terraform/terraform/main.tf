terraform {
  required_version = "~> 0.14.8"

  required_providers {
    google = {
      version = "~> 3.58.0"
    }

    kubectl = {
      source  = "hashicorp/kubernetes"
      version = ">= 2.4"
    }

    local = {
      version = "~> 2.1.0"
    }

    random = {
      version = "~> 3.1.0"
    }
  }
}
