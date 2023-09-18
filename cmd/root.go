/*
Tempest- Leveraging paste sites as a medium for discovery
Copyright © 2023 ax-i-om <addressaxiom@pm.me>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/
// package cmd ...
package cmd

import (
	"os"

	"github.com/ax-i-om/tempest/internal/handlers"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tempest",
	Short: "Leverage paste sites as a medium for discovery of objectionable/infringing materials",
	Long: `
Tempest is a simple, lightweight, and cross-platform solution 
designed to enable individuals to efficiently discover and 
extract active cloud storage/file sharing links from paste 
platforms such as Rentry.co. It was created to address the 
notable uptick in paste sites being used to distribute content 
that violates copyright and piracy statutes.

DISCLAIMER: It is the end user's responsibility to obey all 
applicable local, state, and federal laws/standards/regulations.
Developers assume no liability and are not responsible for any 
misuse or damage caused by this program. By using Tempest, you 
agree to the previous statements.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		handlers.LogErr(err, "exiting root command with status 1 (err)")
		os.Exit(1)
	}
	handlers.LogInfo("successfully exiting root command")
	os.Exit(0)
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.tempest.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "print debug information to the console")
}
