locals {
  service_accounts = [
    "item-api",
    "store-bff",
  ]
}

resource "google_service_account" "this" {
  for_each = toset(local.service_accounts)

  project      = var.project_id
  account_id   = "run-${each.value}"
  display_name = "[Run] ${title(each.value)}"
}

resource "google_project_iam_member" "store_bff_run_invoker" {
  project = var.project_id
  member  = "serviceAccount:${google_service_account.this["store-bff"].email}"
  role    = "roles/run.invoker"
}

resource "google_project_iam_member" "item_api_run_invoker" {
  project = var.project_id
  member  = "serviceAccount:${google_service_account.this["item-api"].email}"
  role    = "roles/run.invoker"
}

resource "google_project_iam_member" "item_api_datastore_user" {
  project = var.project_id
  member  = "serviceAccount:${google_service_account.this["item-api"].email}"
  role    = "roles/datastore.user"
}

resource "google_service_account" "service_discovery" {
  project      = var.project_id
  account_id   = "run-service-discovery"
  display_name = "[Run] Service Discovery"
}

resource "google_project_iam_member" "service_discovery" {
  project = var.project_id
  member  = "serviceAccount:${google_service_account.service_discovery.email}"
  role    = "roles/datastore.user"
}
