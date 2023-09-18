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
// Package handlers contains functions used throughout tempest
package handlers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ax-i-om/tempest/internal/globals"
	"github.com/ax-i-om/tempest/pkg/models"
)

const (
	leIndexBits = 6
	leIndexMask = 1<<leIndexBits - 1
	leIndexMax  = 63 / leIndexBits
)

// TrueRand generates a random number for paste ID generation
func TrueRand(n int, chars string) string {

	b := make([]byte, n)
	for i, cache, remain := n-1, globals.Src.Int63(), leIndexMax; i >= 0; {
		if remain == 0 {
			cache, remain = globals.Src.Int63(), leIndexMax
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

// SwapCheck will be used in the future to determine whether or not to swap to a new proxy
// will be migrated to an exported function, likely in internal to handle modular rate limiting
func SwapCheck(err error) {
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

// FixName cleans/appends an extension (substr) to a string (str) if necessary
func FixName(str, substr string) string {
	if strings.Contains(str, substr) {
		return str
	} else {
		t := strings.ReplaceAll(str, `.`, ``) + substr
		return strings.ReplaceAll(t, ` `, ``)
	}
}

// Write is used to write all entries from results to the specified file/output
func Write(results []models.Entry) {
	// Loop through all entries in results
	for _, v := range results {
		globals.WriteMutex.Lock()
		switch globals.Mode {
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
			_, err = globals.Jsonfile.WriteString(string(vByte) + ",\n")
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
				continue // immediately start next iteration
			}
		case "csv": // If mode is set to csv:
			// Create a CSV record based on the current iteration's accompanying entry (v)
			row := []string{v.Source, v.Link, v.Title, v.Description, v.Service, v.Uploaded, v.Type, v.Size, fmt.Sprint(v.FileCount), v.Thumbnail, fmt.Sprint(v.Downloads), fmt.Sprint(v.Views)}
			// Write the record
			err := globals.Writer.Write(row)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
				continue // immediately start next iteration
			}
			// Call flush to ensure that the record is written to the CSV file
			globals.Writer.Flush()
		}
		globals.WriteMutex.Unlock()
	}
}

// Wipe is used to close and flush all opened files/writers
func Wipe() {
	// if jsonfile was assigned a value other than nil, close it
	if globals.Jsonfile != nil {
		globals.Jsonfile.Close()
	}
	// if writer was assigned a value other than nil, close it
	if globals.Writer != nil {
		globals.Writer.Flush()
	}
	// if csvfile was assigned a value other than nil, close it
	if globals.Csvfile != nil {
		globals.Csvfile.Close()
	}
}

// Deduplicate removes any duplicate lines from a string and writes the results to a cleaned file
func Deduplicate(filename string) error {
	// Open file
	opened, err := os.Open(filename)
	if err != nil {
		return err
	}
	// New scanner, set split function to scan lines
	scanner := bufio.NewScanner(opened)
	scanner.Split(bufio.ScanLines)

	var entries []string

	// Append scanned lines to entries
	for scanner.Scan() {
		entries = append(entries, scanner.Text())
	}

	// Close opened file, as we have already read all the contents
	opened.Close()

	emk := make(map[string]bool)
	// Results slice
	var deduped []string
	for _, item := range entries {
		if _, value := emk[item]; !value {
			emk[item] = true
			deduped = append(deduped, item)
		}
	}

	// Emptry string (will be written to file)
	dedupedToString := ""

	// Append each entry to the string, with a newline
	for _, entry := range deduped {
		dedupedToString += entry + "\n"
	}

	// Write string to clean file, create if it doesn't exist
	err = os.WriteFile("clean-"+filename, []byte(dedupedToString), 0600)
	if err != nil {
		return err
	}
	return nil
}

// GetRes sends a request to a specified URL and returns and *http.Response and error that will be nil if successful.
func GetRes(link string) (*http.Response, error) {
	method := "GET"
	client := &http.Client{
		Timeout: time.Second * 15, // connection timeout after 15 seconds
	}

	req, err := http.NewRequest(method, link, nil)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
