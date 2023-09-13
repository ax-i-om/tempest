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

var mode, filename string

var existed bool

var src = rand.NewSource(time.Now().UnixNano())

var wg models.WaitGroupCount = models.WaitGroupCount{}

const (
	leIndexBits = 6
	leIndexMask = 1<<leIndexBits - 1
	leIndexMax  = 63 / leIndexBits
)

var jsonfile *os.File = nil

var csvfile *os.File = nil
var writer *csv.Writer = nil

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

func wipe() {
	if jsonfile != nil {
		jsonfile.Close()
	}
	if writer != nil {
		writer.Flush()

	}
	if csvfile != nil {
		csvfile.Close()
	}
}

func swapCheck(err error) {
	e := err.Error()
	if !strings.Contains(e, "Get") && !strings.Contains(e, "EOF") {
		if strings.Contains(e, "connection reset by peer") || strings.Contains(e, "client connection force closed via ClientConn.Close") || strings.Contains(e, "closed") {
			// SWAP HERE
			// Ignore this
		} else {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			wipe()
			os.Exit(1)
		}
	}
}

func fixName(str, substr string) string {
	if strings.Contains(str, substr) {
		return str
	} else {
		t := strings.ReplaceAll(str, `.`, ``) + substr
		return strings.ReplaceAll(t, ` `, ``)
	}
}

func printUsage() {
	fmt.Println("Usage:")
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

		for _, v := range results {
			switch mode {
			case "console":
				fmt.Println(v.Service, ": ", v.Link)
			case "json":
				vByte, err := json.Marshal(v)
				if err != nil {
					fmt.Fprintf(os.Stderr, "%s\n", err)
					continue
				}
				_, err = jsonfile.WriteString(string(vByte) + ",\n")
				if err != nil {
					fmt.Fprintf(os.Stderr, "%s\n", err)
					continue
				}
			case "csv":
				row := []string{v.Link, v.LastValidation, v.Title, v.Description, v.Service, v.Uploaded, v.Type, v.Size, v.Length, v.FileCount, v.Thumbnail, v.Downloads, v.Views}
				err := writer.Write(row)
				if err != nil {
					fmt.Fprintf(os.Stderr, "%s\n", err)
					continue
				}
				writer.Flush()
			}
		}
	}

	return res.Body.Close()
}

func main() {
	// Printing license information to the terminal
	fmt.Println("Tempest Copyright (C) 2023 Axiom\nThis program comes with ABSOLUTELY NO WARRANTY.\nThis is free software, and you are welcome to redistribute it\nunder certain conditions.")
	fmt.Println()

	args := os.Args

	var err error

	for i := range args {
		args[i] = strings.ToLower(args[i])
	}

	if len(args) < 2 {
		printUsage()
		os.Exit(0)
	} else if len(args) == 2 {
		if args[1] == "console" {
			mode = "console"
			filename = ""
			fmt.Println("Output Mode: Console")
			fmt.Println("")
		} else {
			printUsage()
			fmt.Fprintf(os.Stderr, "%s\n", errors.New("file name/path not specified"))
			os.Exit(0)
		}
	} else if len(args) == 3 {
		switch args[1] {
		case "json":
			if args[2] == "" || len(args[2]) < 1 {
				printUsage()
				fmt.Fprintf(os.Stderr, "%s\n", errors.New("json file name/path not specified"))
				os.Exit(0)
			} else {
				mode = "json"
				filename = fixName(args[2], ".json")
				fmt.Println("Output Mode: JSON")
				fmt.Println("File Name: ", filename)
				fmt.Println()
				jsonfile, err = os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
				if err != nil {
					wipe()
					fmt.Fprintf(os.Stderr, "%s\n", err)
					os.Exit(1)
				}

			}
		case "csv":
			if args[2] == "" || len(args[2]) < 1 {
				printUsage()
				fmt.Fprintf(os.Stderr, "%s\n", errors.New("json file name/path not specified"))
				os.Exit(0)
			} else {
				mode = "csv"
				filename = fixName(args[2], ".csv")
				fmt.Println("Output Mode: CSV")
				fmt.Println("File Name: ", filename)
				fmt.Println()
				csvfile, err = os.OpenFile(filename, os.O_WRONLY|os.O_APPEND, 0600)
				existed = true
				if err != nil {
					if errors.Is(err, os.ErrNotExist) {
						csvfile, err = os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
						if err != nil {
							fmt.Fprintf(os.Stderr, "%s\n", err)
							wipe()
							os.Exit(1)
						}
						existed = false
					} else {
						wipe()
						fmt.Fprintf(os.Stderr, "%s\n", err)
						os.Exit(1)
					}
				}
				writer = csv.NewWriter(csvfile)
				if !existed {
					headers := []string{"link", "lastvalidation", "title", "description", "service", "uploaded", "type", "size", "length", "filecount", "thumbnail", "downloads", "views"}
					err := writer.Write(headers)
					if err != nil {
						wipe()
						fmt.Fprintf(os.Stderr, "%s\n", err)
						os.Exit(1)
					}
					writer.Flush()
				}
			}
		case "clean":
			if args[2] == "" || len(args[2]) < 1 {
				printUsage()
				fmt.Fprintf(os.Stderr, "%s\n", errors.New("json file name/path not specified"))
				os.Exit(0)
			} else {
				mode = "clean"
				filename = fixName(args[2], ".json")
				fmt.Println("Output Mode: CLEAN")
				fmt.Println("File Name: ", filename)
				fmt.Println()

				content, err := os.ReadFile(filename)
				if err != nil {
					fmt.Fprintf(os.Stderr, "%s\n", err)
					os.Exit(1)
				}

				middle := strings.TrimRight(string(content), "\n")
				middle = strings.TrimRight(middle, ",")
				middle = strings.ReplaceAll(middle, "{\"link\":\"", "\t\t{\"link\":\"")

				comp := "{\n\t\"content\":[\n" + middle + "\n\t]\n}"

				err = os.WriteFile("clean-"+filename, []byte(comp), 0600)
				if err != nil {
					fmt.Fprintf(os.Stderr, "%s\n", err)
					os.Exit(1)
				}

				fmt.Println("Finished cleaning", filename)
				fmt.Println("Cleaned file name: clean-" + filename)
				os.Exit(0)
			}
		default:
			fmt.Fprintf(os.Stderr, "%s\n", errors.New("unrecognized output mode"))
			os.Exit(0)
		}
	} else {
		fmt.Fprintf(os.Stderr, "%s\n", errors.New("malformed command arguments"))
		os.Exit(0)
	}

	cntx := context.Background()
	cntx, cancel := context.WithCancel(cntx)

	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, os.Interrupt)

	defer func() {
		signal.Stop(sigChannel)
		cancel()
	}()

	go func() {
		select {
		case <-sigChannel: // graceful
			cancel()
		case <-cntx.Done():
		}
		<-sigChannel // forceful
		os.Exit(2)
	}()

	err = run(cntx, os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func run(cntx context.Context, args []string) error {
	for {
		select {
		case <-cntx.Done():
			fmt.Println("  ->  Attempting to gracefully shutdown Tempest")
			fmt.Println("\nWaiting for", wg.GetCount(), "GoRoutines to finish execution. Please wait... (~15s)")

			wg.Wait()
			wipe()

			fmt.Println("Tempest was gracefully shut down")

			return nil
		default:
			time.Sleep(1 * time.Millisecond) // Sleeps for 1 millisecond (lol)

			wg.Add(1)

			go func() {
				defer wg.Done()
				err := worker("https://rentry.co/" + trueRand(5, "abcdefghijklmnopqrstuvwxyz0123456789") + "/raw")
				if err != nil {
					swapCheck(err)
				}
			}()
		}
	}
}
