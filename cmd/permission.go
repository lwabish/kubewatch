/*
Copyright 2016 Skippbox, Ltd.

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
	"github.com/bitnami-labs/kubewatch/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// permissionConfigCmd represents the slack subcommand
var permissionConfigCmd = &cobra.Command{
	Use:   "permission",
	Short: "specific permission configuration",
	Long:  `specific permission configuration`,
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := config.New()
		if err != nil {
			logrus.Fatal(err)
		}

		scname, err := cmd.Flags().GetString("scname")
		if err == nil {
			if len(scname) > 0 {
				conf.Handler.Permission.ScName = scname
			}
		} else {
			logrus.Fatal(err)
		}
		chmod, err := cmd.Flags().GetString("chmod")
		if err == nil {
			if len(chmod) > 0 {
				conf.Handler.Permission.Chmod = chmod
			}
		} else {
			logrus.Fatal(err)
		}
		chown, err := cmd.Flags().GetString("chown")
		if err == nil {
			if len(chown) > 0 {
				conf.Handler.Permission.Chown = chown
			}
		}

		if err = conf.Write(); err != nil {
			logrus.Fatal(err)
		}
	},
}

func init() {
	permissionConfigCmd.Flags().StringP("scname", "s", "", "Specify storage class name")
	permissionConfigCmd.Flags().StringP("chmod", "c", "", "Specify permission")
	permissionConfigCmd.Flags().StringP("chown", "", "", "Specify owner")
}
