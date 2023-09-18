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
	"encoding/csv"
	"errors"
	"fmt"
	"os"

	"github.com/ax-i-om/tempest/internal/globals"
	"github.com/ax-i-om/tempest/internal/handlers"
	"github.com/ax-i-om/tempest/internal/worker"
	"github.com/spf13/cobra"
)

// csvCmd represents the csv command
var csvCmd = &cobra.Command{
	Use:   "csv <filename|filepath>",
	Short: "Launch Tempest and output results to the specified CSV file",
	Long: `
Launch Tempest and output results to the specified CSV file

In order to gracefully shut down Tempest, press "Ctrl + C" in 
the terminal **ONCE** and wait until the remaining goroutines 
finish executing (typically <60s) In order to forcefully shut 
down Tempest press "Ctrl + C" in the terminal **TWICE**
CAUTION: FORCEFULLY SHUTTING DOWN TEMPEST MAY RESULT IN ISSUES 
INCLUDING, BUT NOT LIMITED TO, DATA LOSS AND FILE CORRUPTION`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Usage()
		} else {
			launch := true
			var existed bool
			var err error
			// Set output mode to csv
			globals.Mode = "csv"
			// Set filename to args[2], append .csv if necessary
			globals.Filename = handlers.FixName(args[0], ".csv")
			fmt.Println("Output Mode:", globals.Mode)
			fmt.Println("File Name:", globals.Filename)
			fmt.Println()
			// Set the globally declared jsonfile variable to filename
			globals.Csvfile, err = os.OpenFile(globals.Filename, os.O_WRONLY|os.O_APPEND, 0600)
			// Set existed to true, if it didn't exist, this will be set to false
			existed = true
			if err != nil {
				// Check if the error occurred because the file doesn't exist
				if errors.Is(err, os.ErrNotExist) {
					// If file doesn't exist, create one
					globals.Csvfile, err = os.OpenFile(globals.Filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
					if err != nil { // Error when attempting to create CSV file, meaning issues could occur when trying to call write()
						fmt.Fprintf(os.Stderr, "%s\n", err)
						// Close all files/flush all writers
						handlers.Wipe()
						launch = false
					}
					// If error doesn't occur when trying to create file, this means one likely did not already exist or may have been overwritten; therefore,
					// set existed flag to false
					existed = false
				} else { // An error unrelated to a files existence/lack-thereof occurred, resulting in an inability to create/open csvfile
					// Close all files/flush all writers
					handlers.Wipe()
					fmt.Fprintf(os.Stderr, "%s\n", err)
					launch = false
				}
			}
			// Create a new *csv.Writer that writes to csvfile, assign to globally declared variable writer
			globals.Writer = csv.NewWriter(globals.Csvfile)
			if !existed { // Check if the specified csv file already existed by referencing the existed flag, if it did not exist:
				// Create/format headers string slice
				headers := []string{"source", "link", "title", "description", "service", "uploaded", "type", "size", "filecount", "thumbnail", "downloads", "views"}
				// Write headers
				err := globals.Writer.Write(headers)
				if err != nil { //
					// Close all files/flush all writers
					handlers.Wipe()
					fmt.Fprintf(os.Stderr, "%s\n", err)
					launch = false
				}
				// Flush writer to ensure contents were written
				globals.Writer.Flush()
			}
			if launch {
				worker.Launch()
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(csvCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// csvCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// csvCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
