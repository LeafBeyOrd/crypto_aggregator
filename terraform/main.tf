provider "google" {
  project = "your-project-id"
  region  = "your-region"
}

# Cloud Run service
resource "google_cloud_run_service" "default" {
  name     = "crypto-processor"
  location = var.region

  template {
    spec {
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
      }
    }
  }

  autogenerate_revision_name = true
}

# IAM Binding for Cloud Run service account
resource "google_project_iam_member" "run_invoker" {
  project = var.project_id
  role    = "roles/run.invoker"
  member  = "serviceAccount:${google_service_account.default.email}"
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

# Cloud Scheduler Job to trigger Cloud Run
resource "google_cloud_scheduler_job" "crypto_job" {
  name             = "daily-crypto-processor"
  schedule         = "0 0 * * *"  # Runs at midnight every day
  time_zone        = "Etc/UTC"
  http_target {
    http_method = "POST"
    uri         = google_cloud_run_service.default.status[0].url
    body        = <<EOF
{
  "PROCESS_DATE": "${var.date}"
}
EOF
    headers = {
      "Content-Type" = "application/json"
    }
  }
}
