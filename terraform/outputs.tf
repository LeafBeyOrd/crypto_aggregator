output "cloud_run_job_name" {
  value = google_cloud_run_v2_job.crypto_processor_job.name
}

output "gcs_bucket_name" {
  value = google_storage_bucket.crypto_input_bucket.name
}

output "bigquery_table_id" {
  value = "${google_bigquery_dataset.crypto_dataset.dataset_id}.${google_bigquery_table.crypto_table.table_id}"
}
