/*
 * Copyright (C) GRyCAP - I3M - UPV
 *
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package clients

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// MinioClient struct to represent minio clients using aws-sdk-go/service/s3
type minioClient struct {
	s3Client *s3.S3
}

// Download method to get files from minio
func (mc *minioClient) Download(directory, path string) (fileName string, err error) {
	fileName = filepath.Base(path)
	pathSlice := strings.SplitN(strings.Trim(path, "/"), "/", 2)
	bucket := pathSlice[0]
	key := pathSlice[1]

	file, err := os.Create(directory + "/" + fileName)
	if err != nil {
		return "", errors.New("Error creating file")
	}
	defer file.Close()

	result, err := mc.s3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return "", errors.New("Error downloading file: " + err.Error())
	}
	defer result.Body.Close()

	_, err = io.Copy(file, result.Body)
	if err != nil {
		return "", errors.New("Error saving new file")
	}

	return file.Name(), nil
}

// Upload method to push files to minio
func (mc *minioClient) Upload(file, path string) error {
	pathSlice := strings.SplitN(strings.Trim(path, "/"), "/", 2)
	bucket := pathSlice[0]
	key := pathSlice[1]

	f, err := os.Open(file)
	if err != nil {
		return errors.New("Error opening file")
	}
	defer f.Close()

	_, err = mc.s3Client.PutObject(&s3.PutObjectInput{
		Body:   f,
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return errors.New("Error uploading file: " + err.Error())
	}

	return nil
}

func getMinioClient(endpoint, accessKey, secretKey string) StorageClient {
	s3config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(accessKey, secretKey, ""),
		Endpoint:         aws.String(endpoint),
		Region:           aws.String("us-east-1"),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	}
	newSession := session.New(s3config)
	s3Client := s3.New(newSession)

	return &minioClient{
		s3Client: s3Client,
	}
}
