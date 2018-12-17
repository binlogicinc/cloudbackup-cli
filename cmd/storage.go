// Copyright © 2018  Fermin Silva <fermin@binlogic.net>
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

package cmd

import (
	"fmt"
	"github.com/binlogicinc/cloudbackup-cli/api"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var storageCmd = &cobra.Command{
	Use:   "storage",
	Short: "Create, update, remove and get information for backup storages in Binlogic CloudBackup",
}

var storageNew = &cobra.Command{
	Use:     "new",
	Short:   "Add new backup storage to Binlogic CloudBackup",
	PreRunE: checkRequiredFlags,
	RunE: func(cmd *cobra.Command, args []string) error {
		name := getStringFlag(cmd, "name")
		path := getStringFlag(cmd, "path")
		bucket := getStringFlag(cmd, "bucket")
		accessKey := getStringFlag(cmd, "storage-access-key")
		secretKey := getStringFlag(cmd, "storage-secret-key")
		regionEndpoint := getStringFlag(cmd, "region-endpoint")

		stType := getStringFlag(cmd, "storage-type")

		storageType, err := api.ParseStorageType(stType)

		if err != nil {
			return err
		}

		if err := validateStorageParams(storageType, path, bucket, accessKey,
			secretKey, regionEndpoint); err != nil {

			return err
		}

		storage, err := getAPIClient().CreateStorage(name, storageType, path, bucket,
			regionEndpoint, accessKey, secretKey)

		if err != nil {
			return err
		}

		printVerbose("Storage created successfully")

		if getBoolFlag(cmd, "json") {
			fmt.Println(storage.JSONString())
		} else {
			fmt.Println(storage)
		}

		return nil
	},
}

var storageUpdate = &cobra.Command{
	Use:     "update",
	Short:   "Updates a backup storage in Binlogic CloudBackup",
	PreRunE: checkRequiredFlags,
	RunE: func(cmd *cobra.Command, args []string) error {
		storageID := getIntFlag(cmd, "storage-id")

		if storageID == 0 {
			return fmt.Errorf("Storage ID cannot be zero")
		}

		storage, err := getAPIClient().GetStorage(storageID)

		if err != nil {
			return err
		}

		if name := getStringFlag(cmd, "name"); name != "" {
			storage.Name = name
		}

		// if the user passed the storage type
		if stType := getStringFlag(cmd, "storage-type"); stType != "" {
			storageType, err := api.ParseStorageType(stType)

			if err != nil {
				return err
			}

			if storageType != storage.StorageType {
				return fmt.Errorf("Cant change storage type from %s to %s", storage.StorageType, storageType)
			}
		}

		if storage.StorageType == api.STORAGE_LOCAL && getStringFlag(cmd, "path") != "" {
			storage.LocalPath = getStringFlag(cmd, "path")
		} else {
			if bucket := getStringFlag(cmd, "bucket"); bucket != "" {
				storage.Bucket = bucket
			}

			if accessKey := getStringFlag(cmd, "storage-access-key"); accessKey != "" {
				storage.AccessKey = accessKey
			}

			if secretKey := getStringFlag(cmd, "storage-secret-key"); secretKey != "" {
				storage.SecretKey = secretKey
			}

			if regionEndpoint := getStringFlag(cmd, "region-endpoint"); regionEndpoint != "" {
				storage.RegionEndpoint = regionEndpoint
			}
		}

		if err := getAPIClient().UpdateStorage(storage); err != nil {
			return err
		}

		if getBoolFlag(cmd, "json") {
			fmt.Println(storage.JSONString())
		} else {
			fmt.Println(storage)
		}

		return nil
	},
}

var storageDelete = &cobra.Command{
	Use:     "delete",
	Short:   "Delete a backup storage in Binlogic CloudBackup",
	PreRunE: checkRequiredFlags,
	RunE: func(cmd *cobra.Command, args []string) error {
		storageID := getIntFlag(cmd, "storage-id")

		if storageID == 0 {
			return fmt.Errorf("Storage ID cannot be zero")
		}

		if err := getAPIClient().DeleteStorage(storageID); err != nil {
			fmt.Fprint(os.Stderr, err, "\n")
		} else {
			printVerbose("Storage deleted successfully")
		}

		return nil
	},
}

var storageInfo = &cobra.Command{
	Use:     "info",
	Short:   "Get information for a backup storage in Binlogic CloudBackup",
	PreRunE: checkRequiredFlags,
	RunE: func(cmd *cobra.Command, args []string) error {
		storageID := getIntFlag(cmd, "storage-id")

		if storageID == 0 {
			return fmt.Errorf("Storage ID cannot be zero")
		}

		if storage, err := getAPIClient().GetStorage(storageID); err != nil {
			fmt.Fprint(os.Stderr, err, "\n")
		} else {
			if getBoolFlag(cmd, "json") {
				fmt.Println(storage.JSONString())
			} else {
				fmt.Println(storage)
			}
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(storageCmd)
	storageCmd.AddCommand(storageNew)
	storageCmd.AddCommand(storageUpdate)
	storageCmd.AddCommand(storageDelete)
	storageCmd.AddCommand(storageInfo)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	storageInfo.Flags().Int("storage-id", 0, "Storage ID")
	storageInfo.MarkFlagRequired("storage-id")
	storageInfo.Flags().Bool("json", false, "Output info in JSON format")

	storageDelete.Flags().Int("storage-id", 0, "Storage ID")
	storageDelete.MarkFlagRequired("storage-id")

	addCreateStorageFlags(storageNew)
	addCreateStorageFlags(storageUpdate)

	storageUpdate.Flags().Int("storage-id", 0, "Storage ID")
	storageUpdate.MarkFlagRequired("storage-id")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func validateStorageParams(st api.StorageType, path, bucket, accessKey, secretKey,
	regionEndpoint string) error {

	switch st {
	case api.STORAGE_LOCAL:
		if strings.TrimSpace(path) == "" {
			return fmt.Errorf("Local path cannot be empty")
		}

	case api.STORAGE_S3, api.STORAGE_GOOGLE, api.STORAGE_DIGITALOCEAN, api.STORAGE_ALIBABA:
		if strings.TrimSpace(bucket) == "" {
			return fmt.Errorf("Storage bucket cannot be empty")
		}

		if strings.TrimSpace(accessKey) == "" {
			return fmt.Errorf("Storage access key cannot be empty")
		}

		if strings.TrimSpace(secretKey) == "" {
			return fmt.Errorf("Storage secret key cannot be empty")
		}

		if strings.TrimSpace(regionEndpoint) == "" {
			return fmt.Errorf("Storage region endpoint cannot be empty")
		}
	}

	return nil
}

func addCreateStorageFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("json", false, "Output info in JSON format")

	cmd.Flags().String("name", "", "The storage name to show in the control panel")
	cmd.MarkFlagRequired("name")

	cmd.Flags().String("storage-type", "", "The storage type: "+
		"'local', 's3', 'google', 'digitalocean' or 'alibaba' ")

	cmd.MarkFlagRequired("storage-type")

	cmd.Flags().String("path", "", "The local storage full path to store the backups (for ex: '/data/backups')")

	cmd.Flags().String("bucket", "", "The cloud bucket to store the backups into (does not apply to local storage)")
	cmd.Flags().String("storage-access-key", "", "The access key for the cloud storage (does not apply to local storage)")
	cmd.Flags().String("storage-secret-key", "", "The secret key for the cloud storage (does not apply to local storage)")

	cmd.Flags().String("region-endpoint", "", "The cloud storage region endpoint, without https, as reported by your "+
		"provider (for ex: 's3.ap-south-1.amazonaws.com')")
}
