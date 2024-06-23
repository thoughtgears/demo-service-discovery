terraform {
  required_version = ">= 1.8.2"

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.34.0"
    }
  }
}

locals {
  service_apis = [
    "cloudfunctions.googleapis.com",
    "dns.googleapis.com",
    "compute.googleapis.com",
    "vpcaccess.googleapis.com",
  ]
}

resource "google_project_service" "this" {
  for_each = toset(local.service_apis)

  project = var.project_id
  service = each.value
}


