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
  - [Running and Testing the Application](#running-and-testing-the-application)

## Overview

The **Genre Recommendation System** is a Go-based application designed to process asset information, interact with a Data Science (DS) service to predict genres, generate CSV recommendations, and store the results in AWS S3. Leveraging AWS Lambda functions and Localstack for local development, this system ensures scalability, maintainability, and efficient genre categorization of assets.

## Features

- **Event-Driven Architecture**: Utilizes AWS S3 events to trigger processing workflows.
- **AWS Lambda Functions**: Two main Lambdas handle asset processing and CSV generation/upload.
- **Mock Data Science Service**: Simulates genre predictions for testing purposes.
- **Local Development with Localstack**: Emulates AWS services locally for seamless development.

## Project Structure

```text
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
- **Git**: For version control.

## Installation

1. **Clone the Repository**

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

## Running and Testing the Application

To run each component and ensure they work correctly, follow the steps outlined in the [STEPS.md](https://www.notion.so/STEPS.md) file.
