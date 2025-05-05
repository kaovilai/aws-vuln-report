package s3_headbucket_auth

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// from https://github.com/openshift/openshift-velero-plugin/pull/223/files#diff-4f17f1708744bd4d8cb7a4232212efa0e3bfde2b9c7b12e3a23dcc913b9fc2ec
func TestGetBucketRegion(t *testing.T) {
	tests := []struct {
		name    string
		bucket  string
		region  string
		wantErr bool
	}{
		{
			// This should work anonymously, this bucket is made public with Bucket policy
			// {
			// 	"Version": "2012-10-17",
			// 	"Statement": [
			// 		{
			// 			"Sid": "publicList",
			// 			"Effect": "Allow",
			// 			"Principal": "*",
			// 			"Action": "s3:ListBucket",
			// 			"Resource": "arn:aws:s3:::openshift-velero-plugin-s3-auto-region-test-1"
			// 		}
			// 	]
			// }
			// ❯ aws s3api head-bucket --bucket openshift-velero-plugin-s3-auto-region-test-1 --no-sign-request 
			// {
			//     "BucketRegion": "us-east-1",
			//     "AccessPointAlias": false
			// }
			name:    "openshift-velero-plugin-s3-auto-region-test-1",
			bucket:  "openshift-velero-plugin-s3-auto-region-test-1",
			region:  "us-east-1",
			wantErr: false,
		},
		{
			// This should require creds
			// ❯ aws s3api head-bucket --bucket openshift-velero-plugin-s3-auto-region-test-2 --no-sign-request

			// An error occurred (403) when calling the HeadBucket operation: Forbidden 
			name:    "openshift-velero-plugin-s3-auto-region-test-2",
			bucket:  "openshift-velero-plugin-s3-auto-region-test-2",
			region:  "us-west-1",
			wantErr: false,
			// TODO: path/param to creds on ci
		},
		{
			name:    "openshift-velero-plugin-s3-auto-region-test-3",
			bucket:  "openshift-velero-plugin-s3-auto-region-test-3",
			region:  "eu-central-1",
			wantErr: false,
		},
		{
			name:    "openshift-velero-plugin-s3-auto-region-test-4",
			bucket:  "openshift-velero-plugin-s3-auto-region-test-4",
			region:  "sa-east-1",
			wantErr: false,
		},
		{
			name:    "velero-6109f5e9711c8c58131acdd2f490f451", // oadp prow aws bucket name
			bucket:  "velero-6109f5e9711c8c58131acdd2f490f451",
			region:  "us-east-1",
			wantErr: false,
			// TODO: add creds usage here.
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetBucketRegion(tt.bucket)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBucketRegion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.region {
				t.Errorf("GetBucketRegion() = %v, want %v", got, tt.region)
			}
		})
	}
}


// GetBucketRegion returns the AWS region that a bucket is in, or an error
// if the region cannot be determined.
// copied from https://github.com/openshift/openshift-velero-plugin/pull/223/files#diff-da482ef606b3938b09ae46990a60eb0ad49ebfb4885eb1af327d90f215bf58b1
// modified to aws-sdk-go-v2
func GetBucketRegion(bucket string) (string, error) {
	var region string
	// GetBucketRegion will attempt to get the region for a bucket using the client's configured region to determine
	// which AWS partition to perform the query on.
	// Client therefore needs to be configured with region.
	// In local dev environments, you might have ~/.aws/config that could be loaded and set with default region.
	// In cluster/CI environment, ~/.aws/config may not be configured, so set hinting region server explicitly.
	// Also set to use anonymous credentials. If the bucket is private, this function would not work unless we modify it to take credentials.
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion("us-east-1"), // This is not default region being used, this is to specify a region hinting server that we will use to get region from.
	)
	if err != nil {
		return "", err
	}
	region, err = manager.GetBucketRegion(context.Background(), s3.NewFromConfig(cfg), bucket, func(o *s3.Options) {
	    // TODO: get creds from bsl 
		o.Credentials = credentials.NewStaticCredentialsProvider("anon-credentials", "anon-secret", "") // this works with private buckets.. why? supposed to require cred with s3:ListBucket https://docs.aws.amazon.com/AmazonS3/latest/API/API_HeadBucket.html
	})
	if region != "" {
		return region, nil
	}
	return "", errors.New("unable to determine bucket's region: " + err.Error())
}
