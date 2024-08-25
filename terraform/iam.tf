# Service Account for Cloud Run Invoker
resource "google_service_account" "cloud_run_invoker_sa" {
  account_id   = "cloud-run-invoker-sa"
  display_name = "Cloud Run Invoker Service Account"
}

# Merge IAM Bindings for the service account
resource "google_project_iam_binding" "cloud_run_invoker_sa_binding" {
  project = var.project_id
  role    = "roles/run.invoker"
  members = [
    "serviceAccount:${google_service_account.cloud_run_invoker_sa.email}"
  ]
}

resource "google_project_iam_binding" "cloud_scheduler_admin_sa_binding" {
  project = var.project_id
  role    = "roles/cloudscheduler.admin"
  members = [
    "serviceAccount:${google_service_account.cloud_run_invoker_sa.email}"
  ]
}

data "google_compute_default_service_account" "default" {
  project = var.project_id
}

# Grant access to BigQuery
resource "google_project_iam_member" "run_job_bq_access" {
  project = var.project_id
  role    = "roles/bigquery.dataEditor"
  member  = "serviceAccount:${data.google_compute_default_service_account.default.email}"
}

# Grant access to GCS
resource "google_project_iam_member" "run_job_gcs_access" {
  project = var.project_id
  role    = "roles/storage.objectAdmin"
  member  = "serviceAccount:${data.google_compute_default_service_account.default.email}"
}
