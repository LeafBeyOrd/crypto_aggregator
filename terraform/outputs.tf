output "cloud_run_url" {
  value = google_cloud_run_service.service.status[0].url
}

output "gcs_bucket_name" {
  value = google_storage_bucket.input_bucket.name
}

output "bigquery_table_id" {
  value = "${google_bigquery_dataset.dataset.dataset_id}.${google_bigquery_table.table.table_id}"
}
