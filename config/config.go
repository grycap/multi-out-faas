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
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
)

// Config struct used to load the configuration
type Config struct {
	StorageProviders map[string]StorageProvider
	Outputs          []Output
}

// StorageProvider struct used to load storage providers
type StorageProvider struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Auth Auth   `json:"auth"`
}

// Auth struct used to load storage provider authentication
type Auth struct {
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	Endpoint  string `json:"endpoint"`
	Token     string `json:"token"`
	Space     string `json:"space"`
}

// Output struct used to load output configurations
type Output struct {
	StorageProviderName string   `json:"storage_name"`
	Path                string   `json:"path"`
	Suffix              []string `json:"suffix"`
	Prefix              []string `json:"prefix"`
}

type storages struct {
	S3      []StorageProvider `json:"s3"`
	Minio   []StorageProvider `json:"minio"`
	Onedata []StorageProvider `json:"onedata"`
}

type rawConfig struct {
	Storages storages `json:"storages"`
	Outputs  []Output `json:"output"`
}

func convertStorages(s *storages) map[string]StorageProvider {
	storageProviders := make(map[string]StorageProvider)
	if s.S3 != nil {
		for _, s3Prov := range s.S3 {
			s3Prov.Type = "s3"
			storageProviders[s3Prov.Name] = s3Prov
		}
	}
	if s.Minio != nil {
		for _, minioProv := range s.Minio {
			minioProv.Type = "minio"
			storageProviders[minioProv.Name] = minioProv
		}
	}
	if s.Onedata != nil {
		for _, onedataProv := range s.Onedata {
			onedataProv.Type = "onedata"
			storageProviders[onedataProv.Name] = onedataProv
		}
	}
	return storageProviders
}

// ReadConfig function to read the user defined configuration
func ReadConfig(fileContent io.Reader) (*Config, error) {
	jsonConfig, err := ioutil.ReadAll(fileContent)
	if err != nil {
		return nil, errors.New("Error loading config file")
	}
	var c rawConfig
	err = json.Unmarshal(jsonConfig, &c)
	if err != nil {
		return nil, errors.New("Invalid config format")
	}
	config := &Config{
		StorageProviders: convertStorages(&c.Storages),
		Outputs:          c.Outputs,
	}
	return config, nil
}
