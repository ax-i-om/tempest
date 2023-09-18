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
	"fmt"

	"github.com/ax-i-om/tempest/internal/globals"
	"github.com/ax-i-om/tempest/internal/worker"
	"github.com/spf13/cobra"
)

// consoleCmd represents the console command
var consoleCmd = &cobra.Command{
	Use:   "console",
	Short: "Launch Tempest and output results to the terminal",
	Long: `
Launch Tempest and output results to the terminal

In order to gracefully shut down Tempest, press "Ctrl + C" in 
the terminal **ONCE** and wait until the remaining goroutines 
finish executing (typically <60s) In order to forcefully shut 
down Tempest press "Ctrl + C" in the terminal **TWICE**
CAUTION: FORCEFULLY SHUTTING DOWN TEMPEST MAY RESULT IN ISSUES 
INCLUDING, BUT NOT LIMITED TO, DATA LOSS AND FILE CORRUPTION`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		globals.DebugFlag, err = cmd.Flags().GetBool("debug")
		if err != nil {
			fmt.Println("Something went wrong when trying to set Debug mode, continuing without debug")
			globals.DebugFlag = false
		}
		// Set mode to console
		globals.Mode = "console"
		fmt.Println("Output Mode: Console")
		fmt.Println("")
		worker.Launch()
	},
}

func init() {
	rootCmd.AddCommand(consoleCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// consoleCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// consoleCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
