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
// Package worker contains high level functions for handling the execution
// of Tempest
package worker

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/ax-i-om/tempest/internal/globals"
	"github.com/ax-i-om/tempest/internal/handlers"
	"github.com/ax-i-om/tempest/pkg/bunkr"
	"github.com/ax-i-om/tempest/pkg/cyberdrop"
	"github.com/ax-i-om/tempest/pkg/dood"
	"github.com/ax-i-om/tempest/pkg/gofile"
	"github.com/ax-i-om/tempest/pkg/googledrive"
	"github.com/ax-i-om/tempest/pkg/mega"
	"github.com/ax-i-om/tempest/pkg/models"
	"github.com/ax-i-om/tempest/pkg/sendvid"
)

// Worker handles the randomly generated Rentry.co URL and processes the results
func worker(source string) error {
	// Performs a get request on the randomly generated Rentry.co URL.
	res, err := handlers.GetRes(source)
	if err != nil {
		if !strings.Contains(err.Error(), "exceeded") {
			handlers.LogErr(err, "worker failed to perform get request to "+source)
		}
		return err
	}
	// If a Status Code of 200 is returned, that means we randomly generated a valid Rentry.co link and can continue
	if res.StatusCode == 200 {
		// Prepare the contents of the response to be read
		body, err := io.ReadAll(res.Body)
		if err != nil {
			handlers.LogErr(err, "worker failed to read contents of res.body")
			return err
		}
		// Convert the slice of bytes to a string
		conv := string(body)

		// Create results slice
		var results []models.Entry = nil

		// Delegate the string to all specified modules

		// Mega Module
		vMega, err := mega.Delegate(conv, source)
		if err != nil {
			handlers.LogErr(err, "worker failed during delegation to Mega module")
			return err
		}
		results = append(results, vMega...)

		// Gofile Module
		vGofile, err := gofile.Delegate(conv, source)
		if err != nil {
			handlers.LogErr(err, "worker failed during delegation to Gofile module")
			return err
		}
		results = append(results, vGofile...)

		// Sendvid Module
		vSendvid, err := sendvid.Delegate(conv, source)
		if err != nil {
			handlers.LogErr(err, "worker failed during delegation to Senvid module")
			return err
		}
		results = append(results, vSendvid...)

		// Cyberdrop Module
		vCyberdrop, err := cyberdrop.Delegate(conv, source)
		if err != nil {
			handlers.LogErr(err, "worker failed during delegation to Cyberdrop module")
			return err
		}
		results = append(results, vCyberdrop...)

		// Bunkr Module
		vBunkr, err := bunkr.Delegate(conv, source)
		if err != nil {
			handlers.LogErr(err, "worker failed during delegation to Bunkr module")
			return err
		}
		results = append(results, vBunkr...)

		// Google Drive Module
		vGdrive, err := googledrive.Delegate(conv, source)
		if err != nil {
			handlers.LogErr(err, "worker failed during delegation to Google Drive module")
			return err
		}
		results = append(results, vGdrive...)

		// Dood Module
		vDood, err := dood.Delegate(conv, source)
		if err != nil {
			handlers.LogErr(err, "worker failed during delegation to Dood module")
			return err
		}
		results = append(results, vDood...)

		// Call the write function to write the results depending on specified output mode/filename
		handlers.Write(results)
	}

	// Close the response body, return an error/nil
	return res.Body.Close()
}

func run(cntx context.Context) error {
	limited := false
	for {
		select {
		case <-cntx.Done():
			// Iterate through rest of func main() and eventually exit (a few lines in this case)
			return nil
		default:
			time.Sleep(10 * time.Millisecond) // Sleeps for 1 millisecond (lol)

			// Waitgroup count ++
			globals.Wg.Add(1)

			go func() {
				// When func execution is complete, subtract 1 from waitgroup
				defer globals.Wg.Done()
				// Generate a random, 5 char long string [a-z0-9] and place it in the rentry.co string, pass as an argument to worker()
				err := worker("https://rentry.co/" + handlers.TrueRand(5, "abcdefghijklmnopqrstuvwxyz0123456789") + "/raw")
				if err != nil {
					// Call swapcheck on any errors to check if a swap is necessary
					limited = true
					if !strings.Contains(err.Error(), "exceeded") {
						handlers.LogErr(err, "worker failed to function properly, starting swapcheck")
					}
				}
			}()
		}
		if limited {
			fmt.Println("ERROR: Rate limited, consider switching connections")
			handlers.LogErr(errors.New("rate limited"), "Tempest blocked from performing further requests, consider switching proxy/vpn connection")
			return nil
		}
	}
}

// Launch begins the execution of Tempest's scraping
func Launch() {
	// Create new context with cancel
	cntx := context.Background()
	cntx, cancel := context.WithCancel(cntx)

	// Make new signal channel
	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, os.Interrupt)

	go func() {
		select {
		case <-sigChannel: // graceful
			cancel()
		case <-cntx.Done():
		}
		<-sigChannel // forceful
		os.Exit(2)
	}()

	// Launch run function with context
	err := run(cntx)

	// Closes the signal channel
	func() {
		signal.Stop(sigChannel)
		cancel()
	}()

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		// Exit with error
		os.Exit(1)
	}

	fmt.Println("\nAttempting to gracefully shutdown Tempest")
	fmt.Println("\nWaiting for", globals.Wg.GetCount(), "GoRoutines to finish execution. Please wait... (~15s)")

	// Wait for all goroutines to finish execution, shouldn't take longer than 15s due to 15s httpclient.Timeout
	globals.Wg.Wait()
	// Close all files/flush all writers
	handlers.Wipe()
	fmt.Println("Tempest was gracefully shut down")
}
