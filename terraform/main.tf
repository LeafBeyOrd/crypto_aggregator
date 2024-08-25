provider "google" {
  project = var.project_id
  region  = var.region
}

# Create a GCS Bucket
resource "google_storage_bucket" "crypto_input_bucket" {
  name     = var.bucket_name
  location = var.region
}

# BigQuery dataset
resource "google_bigquery_dataset" "crypto_dataset" {
  dataset_id = var.bq_dataset_name
  project    = var.project_id
  location   = var.region
}

# BigQuery table
resource "google_bigquery_table" "crypto_table" {
  dataset_id = google_bigquery_dataset.crypto_dataset.dataset_id
  table_id   = var.bq_table_name
  project    = var.project_id

  schema = jsonencode([
    {
      "name": "date",
      "type": "DATE",
      "mode": "REQUIRED"
    },
    {
      "name": "project_id",
      "type": "STRING",
      "mode": "REQUIRED"
    },
    {
      "name": "number_of_transactions",
      "type": "INTEGER",
      "mode": "REQUIRED"
    },
    {
      "name": "total_volume_usd",
      "type": "FLOAT",
      "mode": "REQUIRED"
    }
  ])
}

# Cloud Run Job for processing crypto data
resource "google_cloud_run_v2_job" "crypto_processor_job" {
  name     = "crypto-processor-job"
  location = var.region

  template {
    template {
        containers {
        image = "gcr.io/${var.project_id}/crypto-processor"

        env {
            name  = "BUCKET_NAME"
            value = var.bucket_name
        }

        env {
            name  = "PROJECT_ID"
            value = var.project_id
        }

        env {
            name  = "BQ_DATASET"
            value = var.bq_dataset_name
        }

        env {
            name  = "BQ_TABLE"
            value = var.bq_table_name
        }

        env {
            name  = "PROCESS_DATE"
            value = var.process_date
        }
        }
    }
  }

  launch_stage = "GA"
}

