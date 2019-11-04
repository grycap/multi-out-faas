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
	"strings"

	//"github.com/grycap/multi-out-faas/config"
	"handler/function/config"
)

var errInvalidProvider = errors.New("Invalid provider")

// StorageClient interface for all storage clients
type StorageClient interface {
	Download(directory, path string) (fileName string, err error)
	Upload(file, path string) error
}

// GetClient factory function to get the appropiate storage client
func GetClient(provider *config.StorageProvider) StorageClient {
	switch providerType := strings.ToLower(provider.Type); providerType {
	case "minio":
		return getMinioClient(provider.Auth.Endpoint, provider.Auth.AccessKey, provider.Auth.SecretKey)
	default:
		return nil
	}
}
