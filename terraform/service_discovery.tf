resource "google_artifact_registry_repository" "remote" {
  project       = var.project_id
  location      = var.region
  format        = "DOCKER"
  mode          = "REMOTE_REPOSITORY"
  repository_id = "github-remote"

  remote_repository_config {
    description                 = "custom docker remote with credentials"
    disable_upstream_validation = true
    docker_repository {
      custom_repository {
        uri = "https://ghcr.io"
      }
    }
  }
}

resource "google_cloud_run_v2_service" "service_discovery" {
  project  = var.project_id
  location = var.region
  name     = "service-discovery"
  ingress  = "INGRESS_TRAFFIC_ALL"

  template {
    containers {
      image = "${var.region}-docker.pkg.dev/${var.project_id}/${google_artifact_registry_repository.remote.name}/thoughtgears/service-discovery:7322fd3"
      env {
        name  = "GCP_PROJECT_ID"
        value = var.project_id
      }
    }

    scaling {
      max_instance_count = 1
      min_instance_count = 0
    }

    service_account                  = google_service_account.service_discovery.email
    max_instance_request_concurrency = 20
  }
}
