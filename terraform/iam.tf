locals {
  service_accounts = [
    "item-api",
    "item-bff"
  ]

  service_apis = [
    "functions.googleapis.com",
    "firestore.googleapis.com",
    "iam.googleapis.com",
  ]
}

resource "google_service_account" "this" {
  for_each = toset(local.service_accounts)

  project      = var.project_id
  account_id   = "cf-${each.value}"
  display_name = "[Function] ${title(each.value)}"
}

resource "google_project_iam_member" "api_firebase" {
  project = var.project_id
  member  = "serviceAccount:${google_service_account.this["item-api"].email}"
  role    = "role/datastore.user"
}

resource "google_firestore_database" "this" {
  project     = var.project_id
  location_id = var.region
  type        = "FIRESTORE_NATIVE"
  name        = "(default)"
}
