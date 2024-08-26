# Crypto Aggregator

## Overview

The Crypto Aggregator is a Go-based application designed to aggregate cryptocurrency transaction data stored in Google Cloud Storage (GCS) and store the results in a Google BigQuery table. The application is deployed using Google Cloud Run Jobs and can be scheduled to run automatically with Google Cloud Scheduler.

## Features

- **Google Cloud Integration**: The application is tightly integrated with Google Cloud services, including GCS for storage, BigQuery for data warehousing, Cloud Run for scalable compute, and Cloud Scheduler for job automation.

- **Automated Deployment and Scheduling**:
  - Terraform scripts are provided to automate the creation of all necessary Google Cloud resources, including Cloud Run Jobs, BigQuery datasets and tables, GCS buckets, and Cloud Scheduler jobs.

## Application Logic

### Environment Variables

The application relies on several environment variables for configuration:

- `BUCKET_NAME`: The name of the GCS bucket where the input CSV file is stored.
- `PROJECT_ID`: The Google Cloud project ID.
- `BQ_DATASET`: The BigQuery dataset name where the results will be stored.
- `BQ_TABLE`: The BigQuery table name where the results will be stored.

### Workflow

1. **Google Cloud Clients Initialization**:
   - The application initializes Google Cloud clients for Storage and BigQuery. These clients are used to read data from GCS and write data to BigQuery.

2. **Reading Input Data from GCS**:
   - The application reads a CSV file from GCS. The file is expected to be located at `gs://<BUCKET_NAME>/<PROCESS_DATE>/input.csv`.

3. **Processing the CSV Data**:
   - Each row in the CSV file contains transaction data, including information about the cryptocurrency used in the transaction.
   - The application extracts the `currencySymbol` and `currencyValueDecimal` from each row.
   - The `currencySymbol` is mapped to its corresponding CoinGecko ID using a pre-fetched list of supported coins.
   - The application fetches the conversion rate for the cryptocurrency to USD for the specific date.
   - It then converts the transaction amount to USD and aggregates it based on the date and project ID.
   - It uses a cache map to reduce the number of API calls to CoinGecko (saving cost and time by > 90%).

4. **Writing Aggregated Data to BigQuery**:
   - The aggregated data, including the total volume in USD and the number of transactions, is written to a BigQuery table.

## Deployment

### Prerequisites

- Google Cloud SDK installed and authenticated.
- GCP account with necessary APIs enabled.
- Docker installed.
- Terraform installed.

### Step-by-Step Deployment

1. **Clone the Repository**

2. **Update Variables**:
   - Modify `variables.tf` to set your `project_id` and `bucket_name` as these values need to be globally unique.

3. **Build and Push the Docker Image**:
    ```sh
    docker build -t gcr.io/<your-project-id>/crypto-processor:latest .
    docker push gcr.io/<your-project-id>/crypto-processor:latest
    ```

4. **Initialize and Apply Terraform**:
    ```sh
    cd terraform
    terraform init
    terraform apply
    ```

    Terraform will set up the following resources:
    - Google Cloud Run Job for executing the Go application.
    - GCS bucket for storing input CSV files.
    - BigQuery dataset and table for storing the results.
    - Cloud Scheduler job to trigger the Cloud Run Job periodically.

## Usage

After the deployment, the Cloud Run Job will be automatically triggered based on the schedule defined in Cloud Scheduler. The job will:

1. Fetch the input CSV file from the GCS bucket for the specific date.
2. Process the data, convert it to USD, and aggregate it.
3. Store the results in the specified BigQuery table.

## Further Improvements

Due to time constraints, this project can be enhanced in several ways:

1. **API Key Management**: Store the CoinGecko API key in Google Secret Manager for enhanced security.
2. **Data Visualization**: Create a dashboard for the BigQuery table using Looker Studio or another BI tool.
3. **CI/CD Pipeline**: Implement a CI/CD pipeline for automated testing, building, and deployment.