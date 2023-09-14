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

package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/ax-i-om/tempest/internal/models"
	"github.com/ax-i-om/tempest/internal/req"
	"github.com/ax-i-om/tempest/pkg/modules/bunkr"
	"github.com/ax-i-om/tempest/pkg/modules/cyberdrop"
	"github.com/ax-i-om/tempest/pkg/modules/dood"
	"github.com/ax-i-om/tempest/pkg/modules/gofile"
	"github.com/ax-i-om/tempest/pkg/modules/googledrive"
	"github.com/ax-i-om/tempest/pkg/modules/mega"
	"github.com/ax-i-om/tempest/pkg/modules/sendvid"
)

// Stores output mode and filename
var mode, filename string

// Checks if a CSV file with filename: filename already exists to determine whether or not to write headers
var existed bool

var src = rand.NewSource(time.Now().UnixNano())

// Custom sync.WaitGroup that implements a counter, used in run() for graceful cleanup
var wg models.WaitGroupCount = models.WaitGroupCount{}

var writeMutex sync.Mutex

// Global declaration of files/writers in order to write/flush from anywhere in main
var jsonfile *os.File = nil
var csvfile *os.File = nil
var writer *csv.Writer = nil

const (
	leIndexBits = 6
	leIndexMask = 1<<leIndexBits - 1
	leIndexMax  = 63 / leIndexBits
)

// Generate a random number for paste ID generation
func trueRand(n int, chars string) string {
	b := make([]byte, n)
	for i, cache, remain := n-1, src.Int63(), leIndexMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), leIndexMax
		}
		if idx := int(cache & leIndexMask); idx < len(chars) {
			b[i] = chars[idx]
			i--
		}
		cache >>= leIndexBits
		remain--
	}

	return string(b)
}

// write is used to write all entries from results to the specified file/output
func write(results []models.Entry) {
	// Loop through all entries in results
	for _, v := range results {
		writeMutex.Lock()
		switch mode {
		case "console": // If mode is set to console, print results to terminal
			fmt.Println(v.Service, ": ", v.Link)
		case "json": // If mode is set to json:
			// JSON encode the current iteration's accompanying entry (v)
			vByte, err := json.Marshal(v)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
				continue // immediately start next iteration
			}
			// Convert the encoded JSON (type []byte) to the previously opened JSON file.
			_, err = jsonfile.WriteString(string(vByte) + ",\n")
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
				continue // immediately start next iteration
			}
		case "csv": // If mode is set to csv:
			// Create a CSV record based on the current iteration's accompanying entry (v)
			row := []string{v.Link, v.LastValidation, v.Title, v.Description, v.Service, v.Uploaded, v.Type, v.Size, v.Length, fmt.Sprint(v.FileCount), v.Thumbnail, fmt.Sprint(v.Downloads), fmt.Sprint(v.Views)}
			// Write the record
			err := writer.Write(row)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
				continue // immediately start next iteration
			}
			// Call flush to ensure that the record is written to the CSV file
			writer.Flush()
		}
		writeMutex.Unlock()
	}
}

// wipe is used to close and flush all opened files/writers
func wipe() {
	// if jsonfile was assigned a value other than nil, close it
	if jsonfile != nil {
		jsonfile.Close()
	}
	// if writer was assigned a value other than nil, close it
	if writer != nil {
		writer.Flush()
	}
	// if csvfile was assigned a value other than nil, close it
	if csvfile != nil {
		csvfile.Close()
	}
}

// swapCheck will be used in the future to determine whether or not to swap to a new proxy
// will be migrated to an exported function, likely in internal to handle modular rate limiting
func swapCheck(err error) {
	// Convert the error to a string
	e := err.Error()
	if !strings.Contains(e, "Get") && !strings.Contains(e, "EOF") {
		// Check contents for an error that may indicate rate limiting or denied request
		if strings.Contains(e, "connection reset by peer") || strings.Contains(e, "client connection force closed via ClientConn.Close") || strings.Contains(e, "closed") {
			// SWAP HERE
			// Ignore this
		} else {
			// Handle other types of errors, these may be fatal
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}
	}
}

// fixName cleans/appends an extension (substr) to a string (str) if necessary
func fixName(str, substr string) string {
	if strings.Contains(str, substr) {
		return str
	} else {
		t := strings.ReplaceAll(str, `.`, ``) + substr
		return strings.ReplaceAll(t, ` `, ``)
	}
}

// printUsage outputs Tempest usage instructions to the terminal.
func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("\tInstructions can also be found in the README.md file")
	fmt.Println()
	fmt.Println("\tgo run main.go help\t- Prints usage information to the terminal")
	fmt.Println()
	fmt.Println("\tgo run main.go console\t- Start Tempest and output results to the terminal")
	fmt.Println()
	fmt.Println("\tgo run main.go {json/csv} <filename/filepath>\t- Start Tempest and output results to file in JSON/CSV format")
	fmt.Println("\t\tJSON Example:\tgo run main.go json results.json")
	fmt.Println("\t\tCSV Example:\tgo run main.go csv results.csv")
	fmt.Println()
	fmt.Println("\tgo run main.go clean <filename/filepath>\t- Clean/Validate JSON file created by Tempest")
	fmt.Println("\t\tExample:\tgo run main.go clean results.json")
	fmt.Println("\t\tNOTE:\tReusing a cleaned file for Tempest output will cause further formatting issues")
	fmt.Println()
	fmt.Println("\tIn order to gracefully shut down Tempest, press `Ctrl + C` in the terminal **ONCE** and wait until the remaining goroutines finish executing (typically <60s)")
	fmt.Println("\tIn order to forcefully shut down Tempest press `Ctrl + C` in the terminal **TWICE**")
	fmt.Println("\tCAUTION: FORCEFULLY SHUTTING DOWN TEMPEST MAY RESULT IN ISSUES INCLUDING, BUT NOT LIMITED TO, DATA LOSS AND FILE CORRUPTION")
	fmt.Println()
}

// worker handles the randomly generated Rentry.co URL and processes the results
func worker(renturl string) error {
	// Performs a get request on the randomly generated Rentry.co URL.
	res, err := req.GetRes(renturl)
	if err != nil {
		return err
	}
	// If a Status Code of 200 is returned, that means we randomly generated a valid Rentry.co link and can continue
	if res.StatusCode == 200 {
		// Prepare the contents of the response to be read
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		// Convert the slice of bytes to a string
		conv := string(body)

		// Create results slice
		var results []models.Entry = nil

		// Delegate the string to all specified modules

		// Mega Module
		vMega, err := mega.Delegate(conv)
		if err != nil {
			return err
		}
		results = append(results, vMega...)

		// Gofile Module
		vGofile, err := gofile.Delegate(conv)
		if err != nil {
			return err
		}
		results = append(results, vGofile...)

		// Sendvid Module
		vSendvid, err := sendvid.Delegate(conv)
		if err != nil {
			return err
		}
		results = append(results, vSendvid...)

		// Cyberdrop Module
		vCyberdrop, err := cyberdrop.Delegate(conv)
		if err != nil {
			return err
		}
		results = append(results, vCyberdrop...)

		// Bunkr Module
		vBunkr, err := bunkr.Delegate(conv)
		if err != nil {
			return err
		}
		results = append(results, vBunkr...)

		// Google Drive Module
		vGdrive, err := googledrive.Delegate(conv)
		if err != nil {
			return err
		}
		results = append(results, vGdrive...)

		// Dood Module
		vDood, err := dood.Delegate(conv)
		if err != nil {
			return err
		}
		results = append(results, vDood...)

		// Call the write function to write the results depending on specified output mode/filename
		write(results)
	}

	// Close the response body, return an error/nil
	return res.Body.Close()
}

func main() {
	// Printing license information to the terminal
	fmt.Println("Tempest Copyright (C) 2023 Axiom\nThis program comes with ABSOLUTELY NO WARRANTY.\nThis is free software, and you are welcome to redistribute it\nunder certain conditions.")
	fmt.Println()

	args := os.Args

	var err error

	// Convert all strings in args to lower case for easier handling
	for i := range args {
		args[i] = strings.ToLower(args[i])
	}

	// Check how many args were passed to Tempest
	if len(args) < 2 { // len(args) < 2 means no extra args were passed to Tempest, only `go run main.go` or `tempest`
		// Print usage information to terminal
		printUsage()
		// Exit successfully
		os.Exit(0)
	} else if len(args) == 2 { // len(args) == 2 means one extra argument was passed to Tempest, most likely "console," for example: `go run main.go console`
		if args[1] == "console" { //
			// Set mode to console
			mode = "console"
			// Set filename to "", as this variable won't be used for console output.
			filename = ""
			fmt.Println("Output Mode: Console")
			fmt.Println("")
		} else { // Argument passed was not "console"
			// Print usage information to terminal
			printUsage()
			fmt.Fprintf(os.Stderr, "%s\n", errors.New("file name/path not specified"))
			// Exit successfully
			os.Exit(0)
		}
	} else if len(args) == 3 { // len(args) == 2 means two extra arguments were passed to Tempest, for example: `go run main.go json results.json`
		switch args[1] {
		case "json": // If json was specified by args[1]
			if args[2] == "" || len(args[2]) < 1 { // If no filename was specified
				// Print usage information to the terminal
				printUsage()
				fmt.Fprintf(os.Stderr, "%s\n", errors.New("json file name/path not specified"))
				// Exit successfully
				os.Exit(0)
			}
			// Set output mode to json
			mode = "json"
			// Set filename to args[2], append .json if necessary
			filename = fixName(args[2], ".json")
			fmt.Println("Output Mode: JSON")
			fmt.Println("File Name: ", filename)
			fmt.Println()
			// Set the globally declared jsonfile variable to filename, create one if it doesn't exist
			jsonfile, err = os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
			if err != nil { // Error when attempting to open/create JSON file, meaning issues could occur when trying to call write()
				// Close all files/flush all writers
				wipe()
				fmt.Fprintf(os.Stderr, "%s\n", err)
				// Exit with error
				os.Exit(1)
			}
		case "csv": // If csv was specified by args[1]
			if args[2] == "" || len(args[2]) < 1 { // If no filename was specified
				// Print usage information to the terminal
				printUsage()
				fmt.Fprintf(os.Stderr, "%s\n", errors.New("json file name/path not specified"))
				// Exit successfully
				os.Exit(0)
			}
			// Set output mode to csv
			mode = "csv"
			// Set filename to args[2], append .csv if necessary
			filename = fixName(args[2], ".csv")
			fmt.Println("Output Mode: CSV")
			fmt.Println("File Name: ", filename)
			fmt.Println()
			// Set the globally declared jsonfile variable to filename
			csvfile, err = os.OpenFile(filename, os.O_WRONLY|os.O_APPEND, 0600)
			// Set existed to true, if it didn't exist, this will be set to false
			existed = true
			if err != nil {
				// Check if the error occurred because the file doesn't exist
				if errors.Is(err, os.ErrNotExist) {
					// If file doesn't exist, create one
					csvfile, err = os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
					if err != nil { // Error when attempting to create CSV file, meaning issues could occur when trying to call write()
						fmt.Fprintf(os.Stderr, "%s\n", err)
						// Close all files/flush all writers
						wipe()
						// Exit with error
						os.Exit(1)
					}
					// If error doesn't occur when trying to create file, this means one likely did not already exist or may have been overwritten; therefore,
					// set existed flag to false
					existed = false
				} else { // An error unrelated to a files existence/lack-thereof occured, resulting in an inability to create/open csvfile
					// Close all files/flush all writers
					wipe()
					fmt.Fprintf(os.Stderr, "%s\n", err)
					// Exit with error
					os.Exit(1)
				}
			}
			// Create a new *csv.Writer that writes to csvfile, assign to globally declared variable writer
			writer = csv.NewWriter(csvfile)
			if !existed { // Check if the specified csv file already existed by referencing the existed flag, if it did not exist:
				// Create/format headers string slice
				headers := []string{"link", "lastvalidation", "title", "description", "service", "uploaded", "type", "size", "length", "filecount", "thumbnail", "downloads", "views"}
				// Write headers
				err := writer.Write(headers)
				if err != nil { //
					// Close all files/flush all writers
					wipe()
					fmt.Fprintf(os.Stderr, "%s\n", err)
					// Exit with error
					os.Exit(1)
				}
				// Flush writer to ensure contents were written
				writer.Flush()
			}

		case "clean": // If csv was specified by args[1]
			if args[2] == "" || len(args[2]) < 1 { // If no filename was specified
				// Print usage information to the terminal
				printUsage()
				fmt.Fprintf(os.Stderr, "%s\n", errors.New("json file name/path not specified"))
				// Exit successfully
				os.Exit(0)
			}
			// Set output mode to clean
			mode = "clean"
			// Set filename to args[2], append .csv if necessary
			filename = fixName(args[2], ".json")
			fmt.Println("Output Mode: CLEAN")
			fmt.Println("File Name: ", filename)
			fmt.Println()

			// Attempt to read the specified json file
			content, err := os.ReadFile(filename)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
				// Calling to wipe() here is unnecessary, as the clean case doesn't assign any files/writers
				// Exit with error
				os.Exit(1)
			}

			// Trim newline from end of file
			middle := strings.TrimRight(string(content), "\n")
			// Trim the rightmost comma from end file
			middle = strings.TrimRight(middle, ",")
			// Append two tabs to the beginning of each entry (formatting)
			middle = strings.ReplaceAll(middle, "{\"link\":\"", "\t\t{\"link\":\"")

			// Combine the strings
			comp := "{\n\t\"content\":[\n" + middle + "\n\t]\n}"

			// Attempt to write the combined strings new a new file, with a name based on the specified filename
			err = os.WriteFile("clean-"+filename, []byte(comp), 0600)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
				// Exit with error
				os.Exit(1)
			}

			fmt.Println("Finished cleaning", filename)
			fmt.Println("Cleaned file name: clean-" + filename)
			// Exit successfully
			os.Exit(0)
		default: // args[1] was set to something other than json/csv/clean
			fmt.Fprintf(os.Stderr, "%s\n", errors.New("unrecognized output mode"))
			// Exit successfully
			os.Exit(0)
		}
	} else { // bad input
		fmt.Fprintf(os.Stderr, "%s\n", errors.New("malformed command arguments"))
		// Exit successfully
		os.Exit(0)
	}

	cntx := context.Background()
	cntx, cancel := context.WithCancel(cntx)

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

	err = run(cntx)

	func() {
		signal.Stop(sigChannel)
		cancel()
	}()

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		// Exit with error
		os.Exit(1)
	}
	// Exit successfully
	os.Exit(0)
}

func run(cntx context.Context) error {
	for {
		select {
		case <-cntx.Done():
			fmt.Println("  ->  Attempting to gracefully shutdown Tempest")
			fmt.Println("\nWaiting for", wg.GetCount(), "GoRoutines to finish execution. Please wait... (~15s)")

			// Wait for all goroutines to finish execution, shouldn't take longer than 15s due to 15s httpclient.Timeout
			wg.Wait()
			// Close all files/flush all writers
			wipe()

			fmt.Println("Tempest was gracefully shut down")

			// Iterate through rest of func main() and eventually exit (a few lines in this case)
			return nil
		default:
			time.Sleep(1 * time.Millisecond) // Sleeps for 1 millisecond (lol)

			// Waitgroup count ++
			wg.Add(1)

			go func() {
				// When func execution is complete, subtract 1 from waitgroup
				defer wg.Done()
				// Generate a random, 5 char long string [a-z0-9] and place it in the rentry.co string, pass as an argument to worker()
				err := worker("https://rentry.co/" + trueRand(5, "abcdefghijklmnopqrstuvwxyz0123456789") + "/raw")
				if err != nil {
					// Call swapcheck on any errors to check if a swap is necessary
					swapCheck(err)
				}
			}()
		}
	}
}
