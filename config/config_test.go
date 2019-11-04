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

package config

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestReadConfig(t *testing.T) {
	testConfig := []string{
		`{
			"output": [
				{
				"storage_name": "s3-bucket",
				"path": "scar-ffmpeg/scar-batch-ffmpeg-split/video-output",
				"suffix": [
					"avi",
					"txt"
				],
				"prefix": [
					"startswith"
				]
				},
				{
				"storage_name": "minio-bucket",
				"path": "scar-ffmpeg/scar-batch-ffmpeg-split",
				"suffix": [
					"wav"
				]
				}
			],
			"storages": {
				"s3": [
				{
					"name": "s3-bucket"
				}
				],
				"minio": [
				{
					"name": "minio-bucket",
					"auth": {
					"access_key": "muser",
					"secret_key": "mpass",
					"endpoint": "http://myminio.example"
					}
				}
				]
			}
		}`,
	}

	expected := []Config{
		Config{
			StorageProviders: map[string]StorageProvider{
				"s3-bucket": StorageProvider{
					Name: "s3-bucket",
					Type: "s3",
				},
				"minio-bucket": StorageProvider{
					Name: "minio-bucket",
					Type: "minio",
					Auth: Auth{
						AccessKey: "muser",
						SecretKey: "mpass",
						Endpoint:  "http://myminio.example",
					},
				},
			},
			Outputs: []Output{
				Output{
					StorageProviderName: "s3-bucket",
					Path:                "scar-ffmpeg/scar-batch-ffmpeg-split/video-output",
					Suffix:              []string{"avi", "txt"},
					Prefix:              []string{"startswith"},
				},
				Output{
					StorageProviderName: "minio-bucket",
					Path:                "scar-ffmpeg/scar-batch-ffmpeg-split",
					Suffix:              []string{"wav"},
				},
			},
		},
	}

	for i, config := range testConfig {
		if c, err := ReadConfig(strings.NewReader(config)); err != nil || !reflect.DeepEqual(*c, expected[i]) {
			fmt.Printf("%+v\n", *c)
			fmt.Println("")
			fmt.Printf("%+v\n", expected[i])
			t.Error("Error reading configuration file")
		}
	}

}
