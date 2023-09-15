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

// Package googledrive contains functions that can be used to accurately extract and validate Google Drive links.
package googledrive

import (
	"io"
	"regexp"
	"strings"

	"github.com/ax-i-om/tempest/internal/handlers"
	"github.com/ax-i-om/tempest/internal/models"
)

// Compile RegEx expressions for extraction of links/metadata

// Extract Google Drive links
var gLink *regexp.Regexp = regexp.MustCompile("(https|http)://drive.google.com/(folder|file|drive)/(d|folders)/(1[a-zA-Z0-9_-]{32}|0[a-zA-Z0-9_-]{27})")
var roughTitle *regexp.Regexp = regexp.MustCompile(`<title>(.*?)</title>`) // Extract Title

// Extract returns a slice of all Google Drive links contained within a string, if any.
func Extract(res string) ([]string, error) {
	// Return all Google Drive links found within an http response
	return gLink.FindAllString(res, -1), nil
}

// ExtractTitle takes the body response/contents of a Google Drive page (raw source/html (formatted as string)) as
// an argument and returns the title as a string.
func ExtractTitle(googledriveContents string) string {
	eTitle := roughTitle.FindString(googledriveContents)     // Extract rough title
	eTitle = strings.ReplaceAll(eTitle, `<title>`, ``)       // Strip opening tags
	eTitle = strings.ReplaceAll(eTitle, `</title>`, ``)      // Strip closing tags
	return strings.ReplaceAll(eTitle, ` - Google Drive`, ``) // Strip extra text
}

// Validate performs a GET request to the Google Drive URL and uses the response status code to identify its validity
func Validate(x string) (bool, error) {
	// Perform a GET request using the Google Drive URL
	res, err := handlers.GetRes(x)
	if err != nil {
		return false, err
	}

	if res.StatusCode == 200 {
		return true, res.Body.Close()
	} else {
		return false, res.Body.Close()
	}
}

// Delegate takes a string as an argument and returns a slice of valid Senvid links found within the response (if any) and an error
func Delegate(res, source string) ([]models.Entry, error) {
	// Use Extract() to extract any existing Google Drive links from the response
	x, err := Extract(res)
	if err != nil {
		return nil, err
	}
	// Check if the return slice of Google Drive links is empty
	if len(x) > 0 {
		// Create a new, empty slice where we will append any valid Google Drive links
		var results []models.Entry = nil
		// Loop through each Google Drive link within the slice
		for _, v := range x {
			// Call the Validate function in order to check whether or not the link is valid
			x, err := Validate(v)
			if err != nil {
				// If any error occurs during the validation process, stop the current iteration and immediately begin with the next link within the slice
				continue
			}
			// If x, the bool return by Validate(), is true: output the result to the terminal and append the link to the specified results slice.
			if x {
				// Get body contents of the sendvid link
				res, err := handlers.GetRes(v)
				if err != nil {
					continue
				}

				// Read results of the *http.Response body
				body, err := io.ReadAll(res.Body)
				if err != nil {
					continue
				}

				// Convert read results to a string
				contents := string(body)

				aTitle := ExtractTitle(contents)

				// Create type Entry and specify the respective values
				ent := models.Entry{Source: source, Link: v, Service: "Google Drive", Title: aTitle}

				if strings.Contains(v, `/file/`) {
					ent.Type = "File"
				} else {
					ent.Type = "Folder"
				}

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
