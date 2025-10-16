# S3 HeadBucket Authentication Tests

[![S3 HeadBucket Auth Tests](https://github.com/kaovilai/aws-vuln-report/actions/workflows/s3-headbucket-auth-tests.yml/badge.svg)](https://github.com/kaovilai/aws-vuln-report/actions/workflows/s3-headbucket-auth-tests.yml)
[![S3 HeadBucket Auth Matrix Tests](https://github.com/kaovilai/aws-vuln-report/actions/workflows/s3-headbucket-auth-matrix.yml/badge.svg)](https://github.com/kaovilai/aws-vuln-report/actions/workflows/s3-headbucket-auth-matrix.yml)

This directory contains tests for the `GetBucketRegion` function, which determines the AWS region of an S3 bucket. The tests include both public and private buckets.

## Expected Behavior Confirmed by AWS Security

**Important Note**: The behavior observed in these tests has been confirmed by AWS Security as **expected and not a security concern**.

While the AWS documentation for the HeadBucket API states that `s3:ListBucket` permission is required, the actual behavior allows bucket region lookups without valid credentials in certain scenarios. This was confirmed by AWS Security (Engagement ID: CACenGS4Mha_KeJ=e3jBSLD6rPZ2iNtfuJUv9QJViaCOt7GVNDg) after investigation.

Key observations:
- The `aws-sdk-go-v2` `GetBucketRegion` function can query bucket regions on private buckets without valid credentials
- The AWS CLI with `--no-sign-request` flag behaves differently and fails as documented
- This is **expected behavior** according to AWS Security, not a vulnerability

This repository serves to document this behavior and provide test coverage for applications that rely on it (such as OADP/Velero plugins).

## GitHub Actions Workflows

Two GitHub Actions workflows have been set up to automatically run these tests:

1. **Basic Workflow**: Runs on push to the main branch, on pull requests to the main branch, and can be triggered manually.
2. **Matrix Testing Workflow**: Runs tests across multiple Go versions and operating systems, with additional features like caching and test result reporting.

### Basic Workflow Configuration

The basic workflow is defined in `.github/workflows/s3-headbucket-auth-tests.yml` and does the following:

1. Checks out the code
2. Sets up Go 1.24.0
3. Configures AWS credentials
4. Runs the tests in this directory

### Required Secrets

To run the tests for private buckets, you need to set up the following GitHub secrets:

- `AWS_ACCESS_KEY_ID`: Your AWS access key ID
- `AWS_SECRET_ACCESS_KEY`: Your AWS secret access key

These credentials should have permissions to access the S3 buckets used in the tests.

### Setting Up GitHub Secrets

1. Go to your GitHub repository
2. Click on "Settings"
3. Click on "Secrets and variables" in the left sidebar
4. Click on "Actions"
5. Click on "New repository secret"
6. Add the secrets mentioned above

### Matrix Testing Workflow Configuration

The matrix testing workflow is defined in `.github/workflows/s3-headbucket-auth-matrix.yml` and does the following:

1. Runs tests on multiple Go versions (1.22, 1.23, 1.24)
2. Runs tests on multiple operating systems (Ubuntu, macOS, Windows)
3. Caches Go dependencies to speed up the workflow
4. Outputs test results in JSON format
5. Uploads test results as artifacts
6. Generates a test summary
7. Runs weekly on Sundays at midnight UTC in addition to push and PR triggers

## Running Tests Locally

To run the tests locally:

```bash
cd s3-headbucket-auth
go test -v ./...
```

Make sure you have AWS credentials configured locally if you want to test with private buckets.

## Test Coverage

To run tests with coverage reporting:

```bash
cd s3-headbucket-auth
go test -v ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

This will generate a coverage report in HTML format that you can open in your browser.
