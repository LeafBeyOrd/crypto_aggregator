variable "project_id" {
  description = "The GCP project ID"
  type        = string
  default     = "noble-courier-433609-s7"
}

variable "region" {
  description = "The GCP region"
  type        = string
  default     = "us-central1"
}

variable "bucket_name" {
  description = "The name of the GCS bucket"
  type        = string
  default     = "crypto_raw_input"
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
