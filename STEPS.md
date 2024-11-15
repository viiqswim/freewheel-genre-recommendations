# Genre Recommendation System - Step-by-Step Guide

This guide provides a detailed, step-by-step process to set up, run, and verify the **Genre Recommendation System** locally. Follow each step to ensure the application functions correctly.

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
- **Docker**: (Optional) For containerizing services.
- **Git**: For version control.

## 2. Setup Localstack

Localstack allows you to emulate AWS services locally. Follow these steps to set it up:

1. **Install Localstack**

   If you haven't installed Localstack, you can do so via `pip`:

   ```bash
   pip install localstack

```

**Alternatively, using Docker:**

```bash
docker pull localstack/localstack

```

1. **Start Localstack**

    **Using Docker:**

    ```bash
    docker run -d -p 4566:4566 -p 4571:4571 --name localstack localstack/localstack
    
    ```

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

    ```
    make_bucket: genre-recommendations
    
    ```

2. **Verify Bucket Creation**

    ```bash
    aws --endpoint-url=http://localhost:4566 s3 ls
    
    ```

    **Expected Output:**

    ```
    2023-10-04 15:20:30 genre-recommendations
    
    ```

## 4. Run the Mock DS Service

The mock DS service simulates genre predictions. Ensure it's running before processing assets.

2. **Run the Mock DS Server**

    ```bash
    go run mocks/ds_mock.go
    
    ```

    **Expected Output:**

    ```
    Mock DS service running on port 9090
    
    ```

    **Note**: If you encounter a port conflict, modify the `DS_PORT` in your `.env` file and restart the server on a different port.

## 5. Upload Sample Asset Information

Upload a sample `asset_info.json` to trigger the Lambda functions.

1. **Create `sample_asset_info.json`**

    Create a file named `sample_asset_info.json` with the following content:

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

    ```
    upload: ./sample_asset_info.json to s3://genre-recommendations/assets/asset_info.json
    
    ```

3. **Verify file was uploaded**

    ```bash
    aws --endpoint-url=http://localhost:4566 s3 ls s3://genre-recommendations/assets/
    ```

    **Expected Output:**

    ```bash
    2024-11-14 18:07:46        162 asset_info.json
    ```

## 6. Run Lambda Functions Locally

Execute the Lambda functions to process assets and generate recommendations.

### a. Process Assets & Send to DS

2. **Run the Lambda Function**

    ```bash
    go run cmd/process_and_send/main.go
    
    ```

    **Expected Output:**

    ```text
    Uploading object to S3 from generated CSV...
    Processed assets and sent to DS successfully.
    
    ```

    **Note**: The function uploads `aggregated_data.json` to S3.

3. **Verify the file was uploaded**

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

    ```
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

    ```
    download: s3://genre-recommendations/recommendations/recommendations.csv to ./recommendations.csv
    
    ```

2. **Check the CSV Content**

    Open `recommendations.csv` with a text editor or use `cat`:

    ```bash
    cat recommendations.csv
    
    ```

    **Expected Content:**

    ```
    ID,Title,Genres
    1,Asset One,Comedy|Drama
    2,Asset Two,Action|Thriller
    3,Asset Three,Sci-Fi|Adventure
    4,Asset Four,Documentary|History
    
    ```

## 8. Testing with Different Asset IDs

To validate the mock DS service's dynamic responses, test with various asset IDs.

1. **Modify `sample_asset_info.json`**

    Change asset IDs to observe different genre predictions.

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
    go run cmd/process_and_sendmain.go
    
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

    ```
    ID,Title,Genres
    5,Asset Five,Documentary|History
    6,Asset Six,Documentary|History
    
    ```

    **Note**: Assets with IDs not `1`, `2`, or `3` default to `["Documentary", "History"]`.
