
# Crypto Aggregator

This project is designed to aggregate cryptocurrency transaction data, read from a CSV file stored in Google Cloud Storage (GCS), and write the aggregated results to a BigQuery table. The system fetches cryptocurrency conversion rates from the CoinGecko API to convert transaction values into USD.

## Requirements

- Go 1.16+
- Docker
- Terraform
- Google Cloud SDK
- A Google Cloud project with billing enabled

## Setup and Deployment

### 1. Local Development

#### Prerequisites

- Ensure you have Go installed on your local machine.
- Install dependencies by running:

  ```bash
  go mod tidy
  ```
This will generate go.sum file 


###  Deploying to Google Cloud

#### Step 1: Configure Google Cloud SDK

Authenticate and set your Google Cloud project:

```bash
gcloud auth login
gcloud config set project noble-courier-433609-s7
```

Then enable necessary APIs:

```bash
gcloud services enable containerregistry.googleapis.com
gcloud services enable run.googleapis.com
gcloud services enable bigquery.googleapis.com
gcloud services enable cloudscheduler.googleapis.com
```

#### step 2: enable container registry and push docker image of go code
1. enable Google Container Registry API
2. docker build -t gcr.io/noble-courier-433609-s7/crypto-processor:latest .
3. gcloud auth configure-docker
4. docker push gcr.io/noble-courier-433609-s7/crypto-processor:latest


#### Step 3: Terraform Setup

Navigate to the `terraform/` directory:

```bash
cd terraform
```

Initialize Terraform:

```bash
terraform init
```

```bash
gcloud auth application-default login
```

Review the Terraform plan:

```bash
terraform plan
```

Apply the Terraform configuration:

```bash
terraform apply
```

This will:

- Set up a Google Cloud Storage (GCS) bucket to store the input CSV files.
- Create a BigQuery dataset and table.
- Deploy the application to Cloud Run jobs.
- Set up a Cloud Scheduler job to trigger the Cloud Run service daily.

#### Step 3: Deploy the Go Application

Build and push the Docker image to Google Container Registry (GCR):

```bash
gcloud builds submit --tag gcr.io/YOUR_PROJECT_ID/crypto-aggregator
```

Update the Cloud Run service to use the newly built image:

```bash
gcloud run deploy crypto-aggregator   --image gcr.io/YOUR_PROJECT_ID/crypto-aggregator   --platform managed   --region YOUR_REGION
```

### 3. Environment Variables

The application uses the following environment variables:

- `BUCKET_NAME`: Name of the GCS bucket containing the input CSV files.
- `PROJECT_ID`: Google Cloud project ID.
- `BQ_DATASET`: BigQuery dataset name.
- `BQ_TABLE`: BigQuery table name.
- `PROCESS_DATE`: Date of the data to process (format: YYYY-MM-DD).

### 4. Google Cloud Scheduler

The Cloud Scheduler job is configured to run the Cloud Run service daily. It passes the current date as `PROCESS_DATE` to the Cloud Run service, which then reads the appropriate CSV file from GCS and writes the aggregated results to BigQuery.
