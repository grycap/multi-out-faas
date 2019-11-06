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

package function

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	//"github.com/grycap/multi-out-faas/clients"
	//"github.com/grycap/multi-out-faas/config"
	//"github.com/grycap/multi-out-faas/events"
	"handler/function/clients"
	"handler/function/config"
	"handler/function/events"
)

// Handle a serverless request
func Handle(req []byte) string {

	// Get the config file name from "CONFIG_FILE" environment variable
	configFileName, ok := os.LookupEnv("CONFIG_FILE")
	if !ok {
		configFileName = "config"
	}
	configFile, err := os.Open("/var/openfaas/secrets/" + configFileName)
	if err != nil {
		log.Println("Error opening config file")
		return ""
	}
	defer configFile.Close()

	config, err := config.ReadConfig(configFile)
	if err != nil {
		log.Println(err.Error())
		return ""
	}

	// Process event
	event, err := events.ReadEvent(string(req))
	if err != nil {
		log.Println(err.Error())
		return ""
	}
	log.Println("Received " + event.EventSource + " event from file '" + event.ObjectKey + "'")

	// Check prefixes and suffixes
	providersToUpload := make(map[string]string)
	var prefixOk, suffixOk bool
	for _, output := range config.Outputs {
		prefixOk = false
		suffixOk = false
		// Prefixes
		if len(output.Prefix) == 0 {
			prefixOk = true
		} else {
			for _, prefix := range output.Prefix {
				if strings.HasPrefix(event.ObjectKey, prefix) {
					prefixOk = true
					break
				}
			}
		}
		if prefixOk {
			// Suffixes
			if len(output.Suffix) == 0 {
				suffixOk = true
			} else {
				for _, suffix := range output.Suffix {
					if strings.HasSuffix(event.ObjectKey, suffix) {
						suffixOk = true
						break
					}
				}
			}
		}
		if prefixOk && suffixOk {
			providersToUpload[output.StorageProviderName] = output.Path
		}
	}

	// If file does not match with any prefix or suffix terminate the function
	if len(providersToUpload) == 0 {
		log.Println("The file '" + event.ObjectKey + "' does not match the specification of any output")
		return ""
	}

	// Manage download
	// Create temporary folder to store the downloaded file
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		log.Println("Error creating file")
		return ""
	}
	defer os.RemoveAll(dir)

	// Get clients for the event source storage providers
	providerClients := make(map[string]clients.StorageClient)
	var fileName string
	for name, provider := range config.StorageProviders {
		if provider.Type == event.EventSource {
			providerClients[name] = clients.GetClient(&provider)
			fileName, err = providerClients[name].Download(dir, event.Path)
			if err != nil {
				log.Println(err.Error())
				continue
			}
			log.Println("File '" + event.ObjectKey + "' successfully downloaded from storage provider '" + provider.Name + "'")
			break
		}
	}
	if fileName == "" {
		log.Println("The file '" + event.ObjectKey + "' cannot be downloaded from any storage provider")
		return ""
	}

	// Manage upload
	for provName, provPath := range providersToUpload {
		client, ok := providerClients[provName]
		uploadPath := provPath + "/" + filepath.Base(fileName)
		// Get the client for specified output
		if !ok {
			for _, provider := range config.StorageProviders {
				if provName == provider.Name {
					client = clients.GetClient(&provider)
					providerClients[provider.Name] = client
				}
			}
		}
		// Upload the file
		err = client.Upload(fileName, uploadPath)
		if err != nil {
			log.Println("Error uploading file '" + fileName + "' to storage provider '" + provName + "'")
		} else {
			log.Println("File '" + fileName + "' successfully uploaded to storage provider '" + provName + "'")
		}
	}

	return ""
}
