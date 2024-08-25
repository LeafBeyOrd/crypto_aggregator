variable "project_id" {
  description = "The GCP project ID"
  type        = string
  default     = "project-id"
}

variable "region" {
  description = "The GCP region"
  type        = string
  default     = "us-central1"
}

variable "bucket_name" {
  description = "The name of the GCS bucket"
  type        = string
  default     = "bucket-name"
}

variable "bq_dataset_name" {
  description = "The name of the BigQuery dataset"
  type        = string
  default     = "crypto_dataset"
}

variable "bq_table_name" {
  description = "The name of the BigQuery table"
  type        = string
  default     = "crypto_transactions"
}

variable "date" {
  description = "Date passed to the job"
  type        = string
  default     = "2024-08-20"  # We can set a default or override it from the Cloud Scheduler
}
