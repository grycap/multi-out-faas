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

package events

import (
	"reflect"
	"testing"
)

func TestReadMinioEvent(t *testing.T) {
	minioEvent := `{
		"Key":"images/nature-wallpaper-229.jpg",
		"Records":[
			{
				"s3":{
					"object":{
						"key":"nature-wallpaper-229.jpg",
						"userMetadata":{
							"content-type":"image/jpeg"
						},
						"eTag":"dd20b7e4b74467ff16ce2d901c054419",
						"contentType":"image/jpeg",
						"sequencer":"153C9A7A7A3FB6AE",
						"versionId":"1",
						"size":1019645
					},
					"s3SchemaVersion":"1.0",
					"bucket":{
						"ownerIdentity":{
							"principalId":"minio"
						},
						"name":"images",
						"arn":"arn:aws:s3:::images"
					},
					"configurationId":"Config"
				},
				"requestParameters":{
					"sourceIPAddress":"10.244.0.0:34852"
				},
				"responseElements":{
					"x-amz-request-id":"153C9A7A7A3FB6AE",
					"x-minio-origin-endpoint":"http://10.244.1.3:9000"
				},
				"source":{
					"userAgent":"",
					"host":"",
					"port":""
				},
				"eventVersion":"2.0",
				"eventName":"s3:ObjectCreated:Put",
				"awsRegion":"",
				"eventTime":"2018-06-29T10:23:44Z",
				"eventSource":"minio:s3",
				"userIdentity":{
					"principalId":"minio"
				}
			}
		],
		"EventName":"s3:ObjectCreated:Put"
	}`

	expected := Event{
		Path:        "images/nature-wallpaper-229.jpg",
		ObjectKey:   "nature-wallpaper-229.jpg",
		EventTime:   "2018-06-29T10:23:44Z",
		EventSource: "minio",
	}

	if event, err := ReadEvent(minioEvent); err != nil || !reflect.DeepEqual(*event, expected) {
		t.Error("Error loading minio event")
	}
}

func TestReadS3Event(t *testing.T) {
	s3Event := `{
		"Records":[
			{
				"awsRegion":"us-east-1",
				"eventName":"ObjectCreated:Put",
				"eventSource":"aws:s3",
				"eventTime":"2019-02-23T11:40:46.473Z",
				"eventVersion":"2.1",
				"requestParameters":{
					"sourceIPAddress":"84.123.4.23"
				},
				"responseElements":{
					"x-amz-id-2":"XXXXX",
					"x-amz-request-id":"XXXXX"
				},
				"s3":{
					"bucket":{
						"arn":"arn:aws:s3:::scar-darknet-bucket",
						"name":"scar-darknet-bucket",
						"ownerIdentity":{
							"principalId":"XXXXX"
						}
					},
					"configurationId":"XXXXX",
					"object":{
						"eTag":"XXXXX",
						"key":"scar-darknet-s3/input/dog.jpg",
						"sequencer":"XXXXX",
						"size":999
					},
					"s3SchemaVersion":"1.0"
				},
				"userIdentity":{
					"principalId":"AWS:XXXXX"
				}
			}
		]
	}`

	expected := Event{
		Path:        "scar-darknet-bucket/scar-darknet-s3/input/dog.jpg",
		ObjectKey:   "scar-darknet-s3/input/dog.jpg",
		EventTime:   "2019-02-23T11:40:46.473Z",
		EventSource: "s3",
	}

	if event, err := ReadEvent(s3Event); err != nil || !reflect.DeepEqual(*event, expected) {
		t.Error("Error loading S3 event")
	}
}

func TestReadOnedataEvent(t *testing.T) {
	onedataEvent := `{
		"Key":"/my-onedata-space/files/file.txt",
		"Records":[
			{
				"objectKey":"file.txt",
				"objectId":"0000034500046EE9C6775...",
				"eventTime":"2019-02-07T09:51:04.347823",
				"eventSource":"OneTrigger"
			}
		]
	}`

	expected := Event{
		Path:        "/my-onedata-space/files/file.txt",
		ObjectKey:   "file.txt",
		EventTime:   "2019-02-07T09:51:04.347823",
		EventSource: "onedata",
	}

	if event, err := ReadEvent(onedataEvent); err != nil || !reflect.DeepEqual(*event, expected) {
		t.Error("Error loading Onedata event")
	}
}

func TestReadInvalidEvents(t *testing.T) {
	tests := []string{
		"",
		"{}",
		"[]",
		`{
			"Records":[
				"eventSource": ""
			]
		}`,
		`{
			"Records":[
				"s3": [""]
			]
		}`,
	}

	for _, test := range tests {
		_, err := ReadEvent(test)
		if err == nil {
			t.Error("Error reading invalid events")
		}
	}
}
