// Copyright © 2017  Fermin Silva <fermin@binlogic.net>
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type Storage struct {
	ID             int         `json:"id"`
	Name           string      `json:"name"`
	StorageType    StorageType `json:"storageType"`
	LocalPath      string      `json:"localPath"`
	Bucket         string      `json:"bucket"`
	AccessKey      string      `json:"storage-access-key"`
	SecretKey      string      `json:"storage-secret-key"`
	RegionEndpoint string      `json:"region-endpoint"`
}

type StorageType int

const (
	STORAGE_LOCAL        StorageType = 1
	STORAGE_S3           StorageType = 2
	STORAGE_GOOGLE       StorageType = 5
	STORAGE_DIGITALOCEAN StorageType = 6
	STORAGE_ALIBABA      StorageType = 7
)

/*
TODO enable this to make nicely printed json for the storage type, but WARNING: this causes marshalling errors when calling the API (as it goes as string instead of int)
*/

// func (s StorageType) MarshalJSON() ([]byte, error) {
// 	return json.Marshal(s.String())
// }

func (s StorageType) String() string {
	switch s {
	case STORAGE_LOCAL:
		return "Local Storage"
	case STORAGE_S3:
		return "AWS S3"
	case STORAGE_GOOGLE:
		return "Google Cloud Storage"
	case STORAGE_DIGITALOCEAN:
		return "DigitalOcean Spaces"
	case STORAGE_ALIBABA:
		return "Alibaba Object Storage"
	}

	return "Unknown"
}

func ParseStorageType(s string) (StorageType, error) {
	switch strings.TrimSpace(strings.ToLower(s)) {
	case "local":
		return STORAGE_LOCAL, nil

	case "s3":
		return STORAGE_S3, nil

	case "google":
		return STORAGE_GOOGLE, nil

	case "digitalocean":
		return STORAGE_DIGITALOCEAN, nil

	case "alibaba":
		return STORAGE_ALIBABA, nil

	}

	return 0, fmt.Errorf("Storage type %s not recognized", s)
}

func (s Storage) String() string {
	switch s.StorageType {
	case STORAGE_LOCAL:
		return fmt.Sprintf("ID: %d\nName: %s\nStorage Type: %s\nPath: %s",
			s.ID, s.Name, s.StorageType, s.LocalPath)

	case STORAGE_S3, STORAGE_GOOGLE, STORAGE_DIGITALOCEAN, STORAGE_ALIBABA:
		secret := "****"

		if len(s.SecretKey) > 4 {
			secret += s.SecretKey[len(s.SecretKey)-4:] //last 4 characters
		}

		return fmt.Sprintf("ID: %d\nName: %s\nStorage Type: %s\nBucket: %s\nRegion Endpoint: %s\nAccess Key: %s\nSecret Key:%s",
			s.ID, s.Name, s.StorageType, s.Bucket, s.RegionEndpoint, s.AccessKey, secret)
	}

	return fmt.Sprintf("ID: %d\nName: %s\nStorage Type: %s", s.ID, s.Name, s.StorageType)
}

func (s Storage) JSONString() string {
	bs, _ := json.Marshal(s)
	j := string(bs)

	if len(s.SecretKey) > 4 {
		secret := "****" + s.SecretKey[len(s.SecretKey)-4:] //last 4 characters

		j = strings.Replace(j, s.SecretKey, secret, -1)
	}

	return j
}

func (c *Client) GetStorage(id int) (storage Storage, err error) {
	resp, err := c.httpClient.SignedGet(c.host+"/storages/"+strconv.Itoa(id), defaultHeaders)

	if err != nil {
		return
	}

	defer resp.Body.Close()

	var body []byte
	body, err = ioutil.ReadAll(resp.Body)

	if resp.StatusCode/100 != 2 {
		if err != nil {
			return
		}

		_, err = c.httpClient.isResponseOk(body)

		if err != nil {
			return
		}

		return storage, fmt.Errorf("Storage API returned HTTP %d but there is no error "+
			"in response '%s' (this should not happen!)", resp.StatusCode, string(body))
	}

	err = json.Unmarshal(body, &storage)

	return
}

func (c *Client) DeleteStorage(id int) error {
	resp, err := c.httpClient.SignedDelete(c.host+"/storages/"+strconv.Itoa(id), defaultHeaders)

	if err != nil {
		return err
	}

	body, val, err := c.httpClient.parseResponseJSON(resp)

	if err != nil {
		return err
	}

	_, err = c.httpClient.isJSONResponseOk(body, val)

	if err != nil {
		return err
	}

	return nil
}

func (c *Client) CreateStorage(name string, storageType StorageType,
	localPath, bucket, region, accessKey, secretKey string) (storage Storage, err error) {

	storage = Storage{ID: 0, Name: name, StorageType: storageType, LocalPath: localPath,
		Bucket: bucket, RegionEndpoint: region, AccessKey: accessKey, SecretKey: secretKey,
	}

	val, err := c.httpClient.postJSON(c.host+"/storages", storage)

	if err != nil {
		err = wrap("while doing client post", err)
		return
	}

	if id, ok := val["id"]; ok {
		if intID, ok2 := id.(float64); ok2 { //json marshalling converts ints to floats
			storage.ID = int(intID)
		}
	}

	if storage.ID <= 0 {
		err = fmt.Errorf("Missing ID from storage response %v", val)
	}

	return
}

func (c *Client) UpdateStorage(s Storage) error {
	if s.ID <= 0 {
		return fmt.Errorf("Invalid ID %d for storage", s.ID)
	}

	if s.Name == "" {
		return fmt.Errorf("Storage name cannot be empty")
	}

	_, err := c.httpClient.postJSON(c.host+"/storages/"+strconv.Itoa(s.ID), s)

	if err != nil {
		return wrap("while doing client post", err)
	}

	return nil
}
