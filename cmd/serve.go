/*
Copyright © 2022 Robert Jackiewicz <rob@jackiewicz.ca>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/rob0t7/domain-go/app"
	"github.com/rob0t7/domain-go/app/postgres"
	"github.com/rob0t7/domain-go/app/rest"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run REST web service",
	Long:  `Runs the REST web server. Responsed to REST requests to the Company Application`,
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := postgres.Open("postgres://postgres:postgres@localhost:5432/app")
		if err != nil {
			return err
		}
		repo := postgres.NewCompanyRepository(db)
		service := app.NewCompanyService(repo)
		server := rest.New(service)
		return server.ListenAndServe()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
