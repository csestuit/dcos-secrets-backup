// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
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
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

// backupCmd represents the backup command
var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		validateCipher()
		if destfile == "secrets.tar" && sourcefile != "secrets.tar" {
			fmt.Println("You specified a source file in a backup command.  Did you mean to specify a destination file?")
			os.Exit(1)
		}

		cluster, err := NewCluster(hostname, username, password)
		if err != nil {
			fmt.Println("Unable to connect to cluster")
			os.Exit(1)
		}

		// secretList, err := cluster.Get("/secrets/v1/secret/default/?list=true")
		secretList, returnCode, err := cluster.Call("GET", "/secrets/v1/secret/default/?list=true", nil)
		if err != nil || returnCode != http.StatusOK {
			fmt.Println("Unable to obtain list of secrets")
			os.Exit(1)
		}

		var secrets struct {
			Array []string `json:array`
		}

		json.Unmarshal(secretList, &secrets)

		secretSlice := []Secret{}

		// Get all secrets, add them to the files array
		for _, secretID := range secrets.Array {
			fmt.Printf("Getting secret '%s'\n", secretID)
			// secretValue, err := cluster.Get("/secrets/v1/secret/default/" + secretPath)
			secretJSON, returnCode, err := cluster.Call("GET", "/secrets/v1/secret/default/"+secretID, nil)
			if err != nil || returnCode != http.StatusOK {
				fmt.Println("TODO: error handling here")
				panic(err)
			}

			e := encrypt(string(secretJSON), cipherkey)
			secretSlice = append(secretSlice, Secret{ID: secretID, EncryptedJSON: e})
		}

		fmt.Println("Writing to tar at " + destfile)
		writeTar(secretSlice, destfile)

		os.Exit(0)

	},
}

func init() {
	rootCmd.AddCommand(backupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// backupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// backupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
