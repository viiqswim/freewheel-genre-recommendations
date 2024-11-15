# Genre Recommendation System

## Table of Contents

- [Genre Recommendation System](#genre-recommendation-system)
  - [Table of Contents](#table-of-contents)
  - [Overview](#overview)
  - [Features](#features)
  - [Project Structure](#project-structure)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Configuration](#configuration)
  - [Running the Application](#running-the-application)
  - [Testing the Services](#testing-the-services)
  - [Deployment](#deployment)
  - [Contributing](#contributing)
  - [License](#license)

## Overview

The **Genre Recommendation System** is a Go-based application designed to process asset information, interact with a Data Science (DS) service to predict genres, generate CSV recommendations, and store the results in AWS S3. Leveraging AWS Lambda functions and Localstack for local development, this system ensures scalability, maintainability, and efficient genre categorization of assets.

## Features

- **Event-Driven Architecture**: Utilizes AWS S3 events to trigger processing workflows.
- **AWS Lambda Functions**: Two main Lambdas handle asset processing and CSV generation/upload.
- **Mock Data Science Service**: Simulates genre predictions for testing purposes.
- **Local Development with Localstack**: Emulates AWS services locally for seamless development.
- **Modular Codebase**: Organized internal packages promote reusability and maintainability.
- **Comprehensive Testing**: Includes unit and integration tests to ensure functionality.

## Project Structure

```

genre_recommendation/
├── cmd/
│   ├── process_and_send/
│   │   └── main.go
│   ├── generate_and_upload_csv/
│   │   └── main.go
├── internal/
│   ├── aws/
│   │   └── s3.go
│   ├── ds/
│   │   └── ds_client.go
│   ├── csv/
│   │   └── generator.go
│   └── config/
│       └── config.go
├── mocks/
│   └── ds_mock.go
├── .env
├── go.mod
└── go.sum

```

- **`cmd/`**: Contains the entry points for Lambda functions.
  - **`process_and_send/`**: Lambda to process assets and send data to the DS service.
  - **`generate_and_upload_csv/`**: Lambda to generate and upload CSV recommendations.
- **`internal/`**: Houses internal packages.
  - **`aws/`**: AWS S3 interactions.
  - **`ds/`**: Communication with the DS service.
  - **`csv/`**: CSV generation logic.
  - **`config/`**: Configuration management.
- **`mocks/`**: Contains mock services for testing.
  - **`ds_mock.go`**: Mock DS service server.
- **`.env`**: Environment variables.
- **`go.mod` & `go.sum`**: Go module dependencies.

## Prerequisites

Before setting up the project, ensure you have the following installed:

- **Go**: Version 1.18 or later.
- **AWS CLI**: Configured for Localstack.
- **Localstack**: To emulate AWS services locally.
- **Docker**: (Optional) For containerizing services.
- **Git**: For version control.

## Installation

1. **Clone the Repository**

```bash
git clone <https://github.com/yourusername/genre_recommendation.git>
cd genre_recommendation
```

1. **Initialize Go Modules**

```bash
go mod tidy
```

## Configuration

1. **Environment Variables**

    Create a `.env` file in the root directory with the following content:

```bash
AWS_REGION=us-east-1
S3_BUCKET=genre-recommendations
ASSET_INFO_KEY=assets/asset_info.json
RECOMMENDATIONS_KEY=recommendations/
DS_SERVICE_URL=http://localhost:9090/predict
DS_PORT=9090
```

**Field Descriptions:**

- **`AWS_REGION`**: AWS region (arbitrary when using Localstack).
- **`S3_BUCKET`**: Name of the S3 bucket used.
- **`ASSET_INFO_KEY`**: S3 key path for the asset information JSON file.
- **`RECOMMENDATIONS_KEY`**: S3 key prefix where the recommendations CSV will be stored.
- **`DS_SERVICE_URL`**: URL of the DS service (mock service URL for local testing).
- **`DS_PORT`**: Port on which the mock DS service runs.

2. **Load Environment Variables**

Ensure the `.env` file is loaded when running the application. The `internal/config` package handles this automatically.

## Running the Application

To run the application locally, follow these steps:

1. **Start Localstack**

    Localstack emulates AWS services locally. Ensure it's running before interacting with S3.

    ```bash
    localstack start
    ```

    **Note**: By default, Localstack services are accessible at `http://localhost:4566`.

2. **Create the S3 Bucket**

    If not already created, use AWS CLI to create the S3 bucket in Localstack:

    ```bash
    aws --endpoint-url=http://localhost:4566 s3 mb s3://genre-recommendations --region us-east-1
    ```

3. **Run the Mock DS Service**

    The mock DS service simulates genre predictions.

    ```bash
    go run mocks/ds_mock.go
    ```

    **Expected Output:**

    ```
    Mock DS service running on port 9090
    ```

4. **Upload Sample Asset Information to S3**

    Create a sample `asset_info.json`:

    ```json
    [
      { "id": "1", "title": "Asset One" },
      { "id": "2", "title": "Asset Two" },
      { "id": "3", "title": "Asset Three" },
      { "id": "4", "title": "Asset Four" }
    ]
    ```

    Save this JSON in a file named `sample_asset_info.json`.

    Upload it to S3:

    ```bash
    aws --endpoint-url=http://localhost:4566 s3 cp sample_asset_info.json s3://genre-recommendations/assets/asset_info.json
    ```

5. **Run Lambda Functions Locally**

    **a. Process Assets & Send to DS**

    ```bash
    go run cmd/process_and_send/main.go
    ```

    **b. Generate & Upload CSV**

    ```bash
    go run cmd/generate_and_upload_csv/main.go
    ```

    **Alternatively**, if triggers are set up, uploading to S3 will automatically invoke the Lambdas.

6. **Verify the Recommendations CSV**

    Download the generated `recommendations.csv` from S3:

    ```bash
    aws --endpoint-url=http://localhost:4566 s3 cp s3://genre-recommendations/recommendations/recommendations.csv .
    ```

    **Check the Content:**

    Open `recommendations.csv` to ensure it contains:

    ```
    ID,Title,Genres
    1,Asset One,Comedy|Drama
    2,Asset Two,Action|Thriller
    3,Asset Three,Sci-Fi|Adventure
    4,Asset Four,Documentary|History
    ```

## Testing the Services

To ensure that each component works correctly, follow the steps outlined in the [STEPS.md](https://www.notion.so/STEPS.md) file below.

## Deployment

Once satisfied with local testing, deploy the Lambdas to AWS using Bingo.

1. **Install Bingo**

    ```bash
    go install github.com/lithammer/bingo@latest
    ```

2. **Configure `bingo.yaml`**

    Update the `bingo.yaml` file to include the two main Lambdas:

    ```yaml
    # bingo.yaml
    functions:
      process_and_send:
        handler: ./cmd/process_and_send/main
        runtime: go1.x
        events:
          - s3:
              bucket: genre-recommendations
              events: s3:ObjectCreated:*
              filter:
                prefix: assets/
                suffix: .json
    
      generate_and_upload_csv:
        handler: ./cmd/generate_and_upload_csv/main
        runtime: go1.x
        events:
          - s3:
              bucket: genre-recommendations
              events: s3:ObjectCreated:*
              filter:
                prefix: aggregated_data.json
                suffix: .json
    ```

3. **Deploy with Bingo**

    ```bash
    bingo deploy
    ```

    **Note**: Ensure your AWS credentials are configured correctly. If using Localstack, adjust deployment settings accordingly.

## Contributing

Contributions are welcome! Please fork the repository and create a pull request with your enhancements.

## License

This project is licensed under the [MIT License](https://www.notion.so/LICENSE).
