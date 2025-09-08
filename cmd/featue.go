/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present Wenhui Shen <swh@admpub.com>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published
   by the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package cmd

import (
	"fmt"
	"strings"

	"github.com/coscms/webcore/library/license"
	"github.com/spf13/cobra"
)

// Nging feature信息

var featureCmd = &cobra.Command{
	Use:   "feature",
	Short: "",
	RunE:  featureRunE,
}

func featureRunE(cmd *cobra.Command, args []string) error {
	if license.SkipLicenseCheck {
		fmt.Println(`<ALL>`)
		return nil
	}
	fmt.Println(strings.Join(license.FeatureList(), `,`))
	return nil
}

func init() {
	rootCmd.AddCommand(featureCmd)
}
