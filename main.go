package main

import (
	"flag"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// 首先，定义所有的环境变量和标志名
const (
	EndpointEnv        = "OSS_ENDPOINT"
	AccessKeyIDEnv     = "OSS_ACCESS_KEY_ID"
	AccessKeySecretEnv = "OSS_ACCESS_KEY_SECRET"
	BucketNameEnv      = "OSS_BUCKET_NAME"
	TargetPathEnv      = "OSS_TARGET_PATH"
	LocalPathEnv       = "OSS_LOCAL_PATH"

	EndpointFlag        = "endpoint"
	AccessKeyIDFlag     = "accessKeyID"
	AccessKeySecretFlag = "accessKeySecret"
	BucketNameFlag      = "bucketName"
	TargetPathFlag      = "targetPath"
	LocalPathFlag       = "localPath"
)

// 定义所有的 flag 变量
var (
	endpoint        = flag.String(EndpointFlag, getEnv(EndpointEnv, ""), "oss endpoint")
	accessKeyID     = flag.String(AccessKeyIDFlag, getEnv(AccessKeyIDEnv, ""), "oss accessKeyID")
	accessKeySecret = flag.String(AccessKeySecretFlag, getEnv(AccessKeySecretEnv, ""), "oss accessKeySecret")
	bucketName      = flag.String(BucketNameFlag, getEnv(BucketNameEnv, ""), "oss bucketName")
	targetPath      = flag.String(TargetPathFlag, getEnv(TargetPathEnv, ""), "oss storage directory")
	localPath       = flag.String(LocalPathFlag, getEnv(LocalPathEnv, ""), "the local directory or file to upload")
)

// getEnv 返回环境变量的值。如果环境变量未定义，它将返回指定的默认值。
func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func main() {
	flag.Parse()

	if *endpoint == "" || *accessKeyID == "" || *accessKeySecret == "" || *bucketName == "" || *targetPath == "" || *localPath == "" {
		log.Fatal("All flags must be set.")
	}

	client, err := oss.New(*endpoint, *accessKeyID, *accessKeySecret)
	if err != nil {
		log.Fatalf("Failed to create a new OSS client: %v", err)
	}

	// Filling in bucket name, e.g., examplebucket.
	bucket, err := client.Bucket(*bucketName)
	if err != nil {
		log.Fatalf("Failed to get bucket: %v", err)
	}

	var files []string

	info, err := os.Stat(*localPath)
	if err != nil {
		log.Fatalf("Failed to get information about the local directory or file: %v", err)
	}

	if !info.IsDir() {
		files = append(files, *localPath)
	} else {
		filepath.Walk(*localPath, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				log.Printf("Failed to access path %q: %v. Continue walking.", path, err)
				return nil
			}

			if !info.IsDir() {
				files = append(files, path)
			}
			return nil
		})
	}

	for _, filename := range files {
		fileInfo, err := os.Stat(filename)
		if err != nil {
			log.Printf("Failed to get information for file %q: %v. Continue uploading next file.", filename, err)
			continue
		}

		remotePath := strings.Join([]string{*targetPath, fileInfo.Name()}, "/")
		err = bucket.PutObjectFromFile(remotePath, filename)
		if err != nil {
			log.Printf("Failed to upload file %q: %v. Continue uploading next file.", filename, err)
			continue
		}

		fmt.Printf("File: %s -> Uploaded to: %s\n", filename, remotePath)
	}
}
