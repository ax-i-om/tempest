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
	"os"

	"github.com/ax-i-om/tempest/internal/globals"
	"github.com/ax-i-om/tempest/internal/handlers"
	"github.com/ax-i-om/tempest/internal/worker"
	"github.com/spf13/cobra"
)

// jsonCmd represents the json command
var jsonCmd = &cobra.Command{
	Use:   "json <filename|filepath>",
	Short: "Launch Tempest and output results to the specified JSON file",
	Long: `
Launch Tempest and output results to the specified JSON file

In order to gracefully shut down Tempest, press "Ctrl + C" in 
the terminal **ONCE** and wait until the remaining goroutines 
finish executing (typically <60s) In order to forcefully shut 
down Tempest press "Ctrl + C" in the terminal **TWICE**
CAUTION: FORCEFULLY SHUTTING DOWN TEMPEST MAY RESULT IN ISSUES 
INCLUDING, BUT NOT LIMITED TO, DATA LOSS AND FILE CORRUPTION`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Usage()
			os.Exit(0)
		}
		var err error
		// Set output mode to json
		globals.Mode = "json"
		// Set filename to args[2], append .json if necessary
		globals.Filename = handlers.FixName(args[0], ".json")
		fmt.Println("Output Mode:", globals.Mode)
		fmt.Println("File Name:", globals.Filename)
		fmt.Println()
		// Set the globally declared jsonfile variable to filename, create one if it doesn't exist
		globals.Jsonfile, err = os.OpenFile(globals.Filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil { // Error when attempting to open/create JSON file, meaning issues could occur when trying to call write()
			// Close all files/flush all writers
			handlers.Wipe()
			fmt.Fprintf(os.Stderr, "%s\n", err)
			// Exit with error
			os.Exit(1)
		}
		worker.Launch()
	},
}

func init() {
	rootCmd.AddCommand(jsonCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// jsonCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// jsonCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
