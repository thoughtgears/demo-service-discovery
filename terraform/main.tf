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
  service_accounts = [
    "item-api",
    "item-bff",
  ]

  service_apis = [
    "cloudfunctions.googleapis.com",
  ]
}

resource "google_project_service" "this" {
  for_each = toset(local.service_apis)

  project = var.project_id
  service = each.value
}

resource "google_service_account" "this" {
  for_each = toset(local.service_accounts)

  project      = var.project_id
  account_id   = "cf-${each.value}"
  display_name = "[Function] ${title(each.value)}"
}

resource "google_project_iam_member" "frontend_run_invoker" {
  project = var.project_id
  member  = "serviceAccount:${google_service_account.this["item-bff"].email}"
  role    = "roles/run.invoker"
}

resource "google_project_iam_member" "frontend_function_invoker" {
  project = var.project_id
  member  = "serviceAccount:${google_service_account.this["item-bff"].email}"
  role    = "roles/cloudfunctions.invoker"
}

