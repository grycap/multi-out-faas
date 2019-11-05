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
	"encoding/json"
	"errors"
	"net/url"
)

// Event struct used to load events
type Event struct {
	Path        string `json:"Key"`
	ObjectKey   string `json:"objectKey"`
	EventTime   string `json:"eventTime"`
	EventSource string `json:"eventSource"`
}

var errInvalidEvent = errors.New("Invalid event")

// ReadEvent function to process raw events
func ReadEvent(rawEvent string) (*Event, error) {
	var eventMap map[string]interface{}

	err := json.Unmarshal([]byte(rawEvent), &eventMap)
	if err != nil {
		return nil, errInvalidEvent
	}

	records, ok := eventMap["Records"].([]interface{})
	if !ok {
		return nil, errInvalidEvent
	}

	record0, ok := records[0].(map[string]interface{})
	if !ok {
		return nil, errInvalidEvent
	}

	// OneTrigger events
	if record0["eventSource"] == "OneTrigger" {
		event := &Event{
			Path:        eventMap["Key"].(string),
			ObjectKey:   record0["objectKey"].(string),
			EventTime:   record0["eventTime"].(string),
			EventSource: "onedata",
		}
		return event, nil
	}

	// MinIO and S3 events
	var source string
	if record0["eventSource"] == "aws:s3" {
		source = "s3"
	} else if record0["eventSource"] == "minio:s3" {
		source = "minio"
	} else {
		// Return error if "eventSource" has unsopported provider
		return nil, errInvalidEvent
	}

	key, ok := record0["s3"].(map[string]interface{})["object"].(map[string]interface{})["key"].(string)
	if !ok {
		return nil, errInvalidEvent
	}
	// Decode url encoded key
	key, err = url.QueryUnescape(key)
	if err != nil {
		return nil, errInvalidEvent
	}

	bucket, ok := record0["s3"].(map[string]interface{})["bucket"].(map[string]interface{})["name"].(string)
	if !ok {
		return nil, errInvalidEvent
	}

	event := &Event{
		Path:        bucket + "/" + key,
		ObjectKey:   key,
		EventTime:   record0["eventTime"].(string),
		EventSource: source,
	}

	return event, nil
}
