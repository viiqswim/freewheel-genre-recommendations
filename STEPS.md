# Genre Recommendation System - Step-by-Step Guide

This guide provides a detailed process to set up, run, and verify the **Genre Recommendation System** locally.

## Table of Contents

- [Genre Recommendation System - Step-by-Step Guide](#genre-recommendation-system---step-by-step-guide)
  - [Table of Contents](#table-of-contents)
  - [1. Prerequisites](#1-prerequisites)
  - [2. Setup Localstack](#2-setup-localstack)
  - [3. Create S3 Bucket](#3-create-s3-bucket)
  - [4. Run the Mock DS Service](#4-run-the-mock-ds-service)
  - [5. Upload Sample Asset Information](#5-upload-sample-asset-information)
  - [6. Run Lambda Functions Locally](#6-run-lambda-functions-locally)
    - [a. Process Assets \& Send to DS](#a-process-assets--send-to-ds)
    - [b. Generate \& Upload CSV](#b-generate--upload-csv)
  - [7. Verify CSV Generation](#7-verify-csv-generation)
  - [8. Testing with Different Asset IDs](#8-testing-with-different-asset-ids)

## 1. Prerequisites

Ensure you have the following installed:

- **Go**: Version 1.18 or later.
- **AWS CLI**: Configured for Localstack.
- **Localstack**: To emulate AWS services locally.
- **Git**: For version control.

## 2. Setup Localstack

Localstack allows you to emulate AWS services locally. Follow these steps to set it up:

1. **Install Localstack**

   ```bash
   pip install localstack
   ```

2. **Start Localstack**

   **Using Localstack CLI:**

   ```bash
   localstack start
   ```

## 3. Create S3 Bucket

Create the necessary S3 bucket in Localstack:

1. **Create the Bucket**

   ```bash
   aws --endpoint-url=http://localhost:4566 s3 mb s3://genre-recommendations --region us-east-1
   ```

   **Expected Output:**

   ```text
   make_bucket: genre-recommendations
   ```

2. **Verify Bucket Creation**

   ```bash
   aws --endpoint-url=http://localhost:4566 s3 ls
   ```

   **Expected Output:**

   ```text
   2023-10-04 15:20:30 genre-recommendations
   ```

## 4. Run the Mock DS Service

The mock DS service simulates genre predictions. Ensure it's running before processing assets.

1. **Run the Mock DS Server**

   ```bash
   go run mocks/ds_mock.go
   ```

   **Expected Output:**

   ```text
   Mock DS service running on port 9090
   ```

## 5. Upload Sample Asset Information

Upload a sample `asset_info.json` to trigger the Lambda functions.

1. **Create `sample_asset_info.json`**

   ```json
   [
     { "id": "1", "title": "Asset One" },
     { "id": "2", "title": "Asset Two" },
     { "id": "3", "title": "Asset Three" },
     { "id": "4", "title": "Asset Four" }
   ]
   ```

2. **Upload to S3**

   ```bash
   aws --endpoint-url=http://localhost:4566 s3 cp sample_asset_info.json s3://genre-recommendations/assets/asset_info.json
   ```

   **Expected Output:**

   ```text
   upload: ./sample_asset_info.json to s3://genre-recommendations/assets/asset_info.json
   ```

3. **Verify file was uploaded**

   ```bash
   aws --endpoint-url=http://localhost:4566 s3 ls s3://genre-recommendations/assets/
   ```

   **Expected Output:**

   ```text
   2024-11-14 18:07:46        162 asset_info.json
   ```

## 6. Run Lambda Functions Locally

Execute the Lambda functions to process assets and generate recommendations.

### a. Process Assets & Send to DS

1. **Run the Lambda Function**

   ```bash
   go run cmd/process_and_send/main.go
   ```

   **Expected Output:**

   ```text
   Uploading object to S3 from generated CSV...
   Processed assets and sent to DS successfully.
   ```

2. **Verify the file was uploaded**

   ```bash
   aws --endpoint-url=http://localhost:4566 s3 ls s3://genre-recommendations/
   ```

   **Expected Output:**

   ```text
                               PRE assets/
   2024-11-14 18:11:03         254 aggregated_data.json
   ```

### b. Generate & Upload CSV

1. **Run the Lambda Function**

   ```bash
   go run cmd/generate_and_upload_csv/main.go
   ```

   **Expected Output:**

   ```text
   Recommendations CSV generated and uploaded successfully.
   ```

2. **Verify the file was uploaded**

   ```bash
   aws --endpoint-url=http://localhost:4566 s3 ls s3://genre-recommendations/recommendations/
   ```

   **Expected Output:**

   ```text
   2024-11-14 18:13:08        133 recommendations.csv
   ```

## 7. Verify CSV Generation

Ensure that the `recommendations.csv` has been correctly uploaded to S3.

1. **Download the CSV from S3**

   ```bash
   aws --endpoint-url=http://localhost:4566 s3 cp s3://genre-recommendations/recommendations/recommendations.csv .
   ```

   **Expected Output:**

   ```text
   download: s3://genre-recommendations/recommendations/recommendations.csv to ./recommendations.csv
   ```

2. **Check the CSV Content**

   ```bash
   cat recommendations.csv
   ```

   **Expected Content:**

   ```text
   ID,Title,Genres
   1,Asset One,Comedy|Drama
   2,Asset Two,Action|Thriller
   3,Asset Three,Sci-Fi|Adventure
   4,Asset Four,Documentary|History
   ```

## 8. Testing with Different Asset IDs

To validate the mock DS service's dynamic responses, test with various asset IDs.

1. **Modify `sample_asset_info.json`**

   ```json
   [
     { "id": "5", "title": "Asset Five" },
     { "id": "6", "title": "Asset Six" }
   ]
   ```

2. **Upload the Modified File to S3**

   ```bash
   aws --endpoint-url=http://localhost:4566 s3 cp sample_asset_info.json s3://genre-recommendations/assets/asset_info.json
   ```

3. **Re-run Lambda Functions Locally**

   **a. Process Assets & Send to DS**

   ```bash
   go run cmd/process_and_send/main.go
   ```

   **b. Generate & Upload CSV**

   ```bash
   go run cmd/generate_and_upload_csv/main.go
   ```

4. **Verify Updated CSV**

   ```bash
   aws --endpoint-url=http://localhost:4566 s3 cp s3://genre-recommendations/recommendations/recommendations.csv .
   cat recommendations.csv
   ```

   **Expected Content:**

   ```text
   ID,Title,Genres
   5,Asset Five,Documentary|History
   6,Asset Six,Documentary|History
   ```
