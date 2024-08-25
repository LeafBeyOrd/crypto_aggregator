# IAM Bindings to allow the default Cloud Run Job service account to access BigQuery and GCS
resource "google_project_iam_member" "run_job_bq_access" {
  project = var.project_id
  role    = "roles/bigquery.dataEditor"
  member  = "serviceAccount:683426209835-compute@developer.gserviceaccount.com"
}

resource "google_project_iam_member" "run_job_gcs_access" {
  project = var.project_id
  role    = "roles/storage.objectAdmin"
  member  = "serviceAccount:683426209835-compute@developer.gserviceaccount.com"
}
