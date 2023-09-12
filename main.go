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
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
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

func fixName(str, substr string) string {
	if strings.Contains(str, substr) {
		return str
	} else {
		return strings.ReplaceAll(str, `.`, ``) + substr
	}
}

func runner(renturl string) error {
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
			if mode == "console" {
				fmt.Println(v.Service, ": ", v.Link)
			} else if mode == "json" {
				vByte, err := json.Marshal(v)
				if err != nil {
					fmt.Println(err)
					continue
				}
				_, err = jsonfile.WriteString(string(vByte[:]) + ",\n")
				if err != nil {
					fmt.Println(err)
					continue
				}
			} else if mode == "csv" {
				row := []string{v.Link, v.LastValidation, v.Title, v.Description, v.Service, v.Uploaded, v.Type, v.Size, v.Length, v.FileCount, v.Thumbnail, v.Downloads, v.Views}
				err := writer.Write(row)
				if err != nil {
					fmt.Println(err)
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
		fmt.Println("Usage")
		return
	} else if len(args) == 2 {
		if args[1] == "console" {
			mode = "console"
			filename = ""
			fmt.Println("Output Mode: Console")
			fmt.Println("")
		} else {
			fmt.Println("Please specify a file name.")
			return
		}
	} else if len(args) == 3 {
		if args[1] == "json" {
			if args[2] == "" || len(args[2]) < 1 {
				fmt.Println("Please specify a file name.")
				return
			} else {
				mode = "json"
				filename = fixName(args[2], ".json")
				fmt.Println("Output Mode: JSON")
				fmt.Println("File Name: ", filename)
				fmt.Println()
				jsonfile, err = os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
				if err != nil {
					log.Fatal(err)
				}

			}
		} else if args[1] == "csv" {
			if args[2] == "" || len(args[2]) < 1 {
				fmt.Println("Please specify a file name.")
				return
			} else {
				mode = "csv"
				filename = fixName(args[2], ".csv")
				fmt.Println("Output Mode: CSV")
				fmt.Println("File Name: ", filename)
				fmt.Println()
				csvfile, err = os.OpenFile(filename, os.O_WRONLY|os.O_APPEND, 0644)
				existed = true
				if err != nil {
					if errors.Is(err, os.ErrNotExist) {
						csvfile, err = os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
						if err != nil {
							log.Fatal(err)
						}
						existed = false
					} else {
						log.Fatal(err)
					}
				}
				writer = csv.NewWriter(csvfile)
				if !existed {
					headers := []string{"link", "lastvalidation", "title", "description", "service", "uploaded", "type", "size", "length", "filecount", "thumbnail", "downloads", "views"}
					err := writer.Write(headers)
					if err != nil {
						fmt.Println(err)
					}
					writer.Flush()
				}
			}
		} else if args[1] == "clean" {
			if args[2] == "" || len(args[2]) < 1 {
				fmt.Println("Please specify a file name.")
				return
			} else {
				mode = "clean"
				filename = fixName(args[2], ".json")
				fmt.Println("Output Mode: CLEAN")
				fmt.Println("File Name: ", filename)
				fmt.Println()
				// START DOING CLEANING STUFF HERE
				fmt.Println("Done!")
				return
			}
		} else {
			fmt.Println("Please properly specify the mode of output")
			return
		}
	} else {
		fmt.Println("Something went wrong when trying to start Tempest.\nPlease check your input and internet connection\nand try again.")
		return
	}

	if jsonfile != nil {
		defer jsonfile.Close()
	}

	if csvfile != nil {
		defer csvfile.Close()
	}

	if writer != nil {
		defer writer.Flush()
	}

	for {
		time.Sleep(1 * time.Millisecond) // Sleeps for 1 millisecond (lol)
		go func() {

			err := runner("https://rentry.co/" + trueRand(5, "abcdefghijklmnopqrstuvwxyz0123456789") + "/raw")
			if err != nil {
				e := err.Error()
				if !strings.Contains(e, "Get") && !strings.Contains(e, "EOF") {
					if strings.Contains(e, "connection reset by peer") || strings.Contains(e, "client connection force closed via ClientConn.Close") || strings.Contains(e, "closed") {
						// SWAP HERE
						// Ignore this
					} else {
						fmt.Println(e)
						os.Exit(1)
					}
				}
			}

		}()

	}
}
