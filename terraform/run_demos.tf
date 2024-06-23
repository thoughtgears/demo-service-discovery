locals {
  git_sha = "6ecd23a"
}

resource "google_cloud_run_v2_service" "item_api" {
  project  = var.project_id
  name     = "item-api"
  location = var.region
  ingress  = "INGRESS_TRAFFIC_ALL"

  template {
    containers {
      image = "${var.region}-docker.pkg.dev/${var.project_id}/demos/item-api:${local.git_sha}"
      resources {
        cpu_idle = true
        limits = {
          cpu    = "1000m"
          memory = "128Mi"
        }
      }
    }
    service_account                  = google_service_account.this["item-api"].email
    max_instance_request_concurrency = 20
    scaling {
      max_instance_count = 1
      min_instance_count = 0
    }
  }
}

resource "google_cloud_run_v2_service" "store_bff" {
  project  = var.project_id
  name     = "store-bff"
  location = var.region
  ingress  = "INGRESS_TRAFFIC_ALL"

  template {
    containers {
      image = "${var.region}-docker.pkg.dev/${var.project_id}/demos/store-bff:${local.git_sha}"
      resources {
        cpu_idle = true
        limits = {
          cpu    = "1000m"
          memory = "128Mi"
        }
      }
      env {
        name  = "DISCOVERY_URL"
        value = google_cloud_run_v2_service.service_discovery.uri
      }
      env {
        name  = "ENVIRONMENT"
        value = "dev"
      }
    }
    service_account                  = google_service_account.this["store-bff"].email
    max_instance_request_concurrency = 20
    scaling {
      max_instance_count = 1
      min_instance_count = 0
    }
  }
}
