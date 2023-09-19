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

// Package cloudmailru contains functions that can be used to accurately extract and validate cloud.mail.ru links.
package cloudmailru

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/ax-i-om/tempest/internal/handlers"
	"github.com/ax-i-om/tempest/pkg/models"
)

// Compile RegEx expressions for extraction of links/metadata
var rLink *regexp.Regexp = regexp.MustCompile(`(https|http)://cloud.mail.ru/public/[a-zA-Z0-9]{4}/[a-zA-Z0-9]{9}`)
var rInfo *regexp.Regexp = regexp.MustCompile(`"serverSideFolders":{(.*?)"DISPATCHERS":`)

// Extract returns a slice of all cloudmailru links contained within a string, if any.
func Extract(res string) ([]string, error) {
	return rLink.FindAllString(res, -1), nil // Return all cloudmailru links found within an http response
}

// ExtractInfo takes the contents of the body response from a valid cloud.mail.ru link and extracts its metadata.
// It unmarshalls the information and returns as type *models.CMRInfo alongside an error
func ExtractInfo(cmailrContents string) (*models.CMRInfo, error) {
	roughInfo := strings.ReplaceAll(strings.ReplaceAll(rInfo.FindString(cmailrContents), `,"DISPATCHERS":`, ``), `"serverSideFolders":`, ``)
	info := new(models.CMRInfo)
	err := json.Unmarshal([]byte(roughInfo), &info)
	if err != nil {
		return nil, err
	}
	return info, nil
}

// Validate takes a Gofile link/URL and checks certain metadata patterns to identify whether or not the link is valid/online.
func Validate(x string) (bool, string, error) {
	// Perform a GET request using the Gofile URL
	res, err := handlers.GetRes(x)
	if err != nil {
		return false, "", err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return false, "", err
	}

	err = res.Body.Close()

	// Read the response, if the title contains the below specified string, then the Gofile link is not online.
	if res.StatusCode == 200 {
		return true, string(body), err
	} else {
		return false, string(body), err
	}
}

// Delegate takes a string as an argument and returns a slice of valid cloudmailru links found within the response (if any) or nil, and an error
func Delegate(res, source string) ([]models.Entry, error) {
	// Use Extract() to extract any existing cloudmailru links from the response
	x, err := Extract(res)
	if err != nil {
		handlers.LogErr(err, "error occurred on cloudmailru delegate attempt to call extract")
		return nil, err
	}
	// Check if the return slice of cloudmailru links is empty
	if len(x) > 0 {
		// Create a new, empty slice where we will append any valid cloudmailru links
		var results []models.Entry = nil
		// Loop through each cloudmailru link within the slice
		for _, v := range x {
			// Call the Validate function in order to check whether or not the link is valid
			x, contents, err := Validate(v)
			if err != nil {
				// If any error occurs during the validation process, stop the current iteration and immediately begin with the next link within the slice
				handlers.LogErr(err, "error occurred on cloudmailru delegate attempt to call validate")
				continue
			}
			// If x, the bool return by Validate(), is true: output the result to the terminal and append the link to the specified results slice.
			if x {
				cmr, err := ExtractInfo(contents)
				if err != nil {
					handlers.LogErr(err, "error occurred on cloudmailru delegate attempt to extract file metadata")
					continue
				}

				// Create type Entry and specify the respective values
				ent := models.Entry{Source: source, Link: v, Service: "CloudMailRu", Title: cmr.Name, Size: fmt.Sprint(cmr.Size), Type: cmr.Type, Mtime: fmt.Sprint(cmr.Mtime), Hash: cmr.Hash, Malware: cmr.Malware.Status}
				// Append the entry to the results slice to be returned to the main runner
				results = append(results, ent)
			}
		}
		// When the loop is finished, return the results slice
		return results, nil
	}
	// Return nothing, if nothing happens (bruh)
	return nil, nil
}
