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

// Package gofile contains functions that can be used to accurately extract and validate Gofile links.
package gofile

import (
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/ax-i-om/tempest/internal/handlers"
	"github.com/ax-i-om/tempest/pkg/models"
)

// Compile RegEx expressions for extraction of links/metadata
var gLink *regexp.Regexp = regexp.MustCompile("(https|http)://gofile.io/d/([a-zA-Z0-9]{6})")     // Extract Gofile link
var roughTitle *regexp.Regexp = regexp.MustCompile(`<title>(.*?)</title>`)                       // Extract title
var roughDesc *regexp.Regexp = regexp.MustCompile(`<meta name='description' content='(.*?)' />`) // Extract file count if folder, download count if file

// Extract returns a slice of all Gofile links contained within a string, if any.
func Extract(res string) ([]string, error) {
	// Return all Gofile links found within an http response
	return gLink.FindAllString(res, -1), nil
}

// ExtractTitle takes the body response/contents of a Gofile page (raw source/html (formatted as string)) as
// an argument and returns the title as a string.
func ExtractTitle(gofileContents string) string {
	rt := roughTitle.FindString(gofileContents)   // Extract rough title
	r1 := strings.ReplaceAll(rt, `<title>`, ``)   // Strip opening tag
	return strings.ReplaceAll(r1, `</title>`, ``) // Strip closing tag and return
}

// ExtractFileCount takes the body response/contents of a Gofile page (raw source/html (formatted as string)) as
// an argument and returns the album's file count as an integer. It will return -1 in the case of a syntax error.
func ExtractFileCount(gofileContents string) int {
	eDesc := roughDesc.FindString(gofileContents)                               // Extract rough description
	eDesc = strings.ReplaceAll(eDesc, `<meta name='description' content='`, ``) // Strip unnecessary html
	eDesc = strings.ReplaceAll(eDesc, `' />`, ``)                               // Strip unnecessary html
	eCount := strings.ReplaceAll(eDesc, ` files`, ``)                           // Strip unnecessary text
	eCount = strings.ReplaceAll(eCount, ",", "")                                // Strip commas (,)
	eCount = strings.ReplaceAll(eCount, ".", "")                                // Strip periods (.)
	fileCount, err := strconv.Atoi(eCount)
	if err != nil {
		return -1 // Return -1 to signify an error occurred and the filecount could not be converted to Int
	}
	return fileCount
}

// ExtractDownloadCount takes the body response/contents of a Gofile page (raw source/html (formatted as string)) as
// an argument and returns the album's file count as an integer. It will return -1 in the case of a syntax error.
func ExtractDownloadCount(gofileContents string) int {
	eDesc := roughDesc.FindString(gofileContents)                                   // Extract rough description
	eDesc = strings.ReplaceAll(eDesc, `<meta name='description' content='`, ``)     // Strip unnecessary html
	eDesc = strings.ReplaceAll(eDesc, `' />`, ``)                                   // Strip unnecessary html
	downloadCount, err := strconv.Atoi(strings.ReplaceAll(eDesc, ` downloads`, ``)) // Strip text, convert to string
	if err != nil {
		return -1 // Return -1 to signify an error occurred and the downloadcount could not be converted to Int
	}
	return downloadCount
}

// Validate takes a Gofile link/URL and checks certain metadata patterns to identify whether or not the link is valid/online.
func Validate(x string) (bool, string, error) {
	// Perform a GET request using the Gofile URL
	res, err := handlers.GetRes(x)
	if err != nil {
		return false, "", err
	}

	// Prepare the contents of the response to be read
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return false, "", err
	}

	// Read the response, if the title contains the below specified string, then the Gofile link is not online.
	if strings.Contains(string(body), "<title>Gofile - Free Unlimited File Sharing and Storage</title>") {
		return false, string(body), res.Body.Close()
	} else {
		return true, string(body), res.Body.Close()
	}
}

// Delegate takes a string as an argument and returns a slice of valid Gofile links found within the response (if any) and an error
func Delegate(res, source string) ([]models.Entry, error) {
	// Use Extract() to extract any existing Gofile links from the response
	x, err := Extract(res)
	if err != nil {
		handlers.LogErr(err, "error occurred on gofile delegate attempt to call extract")
		return nil, err
	}
	// Check if the return slice of Gofile links is empty
	if len(x) > 0 {
		// Create a new, empty slice where we will append any valid Gofile links
		var results []models.Entry = nil
		// Loop through each Gofile link within the slice
		for _, v := range x {
			// Call the Validate function in order to check whether or not the link is valid
			x, contents, err := Validate(v)
			if err != nil {
				// If any error occurs during the validation process, stop the current iteration and immediately begin with the next link within the slice
				handlers.LogErr(err, "error occurred on gofile delegate attempt to call validate")
				continue
			}
			// If x, the bool return by Validate(), is true: output the result to the terminal and append the link to the specified results slice.
			if x {
				title := ExtractTitle(contents) // Extract title

				// Create type Entry and specify the respective values
				ent := models.Entry{Source: source, Link: v, Service: "GoFile"}

				if strings.Contains(title, "Folder") {
					ent.Title = strings.ReplaceAll(title, `Folder `, ``)
					ent.FileCount = ExtractFileCount(contents)
					ent.Type = "Folder"
				} else {
					ent.Title = title
					ent.Downloads = ExtractDownloadCount(contents)
					ent.Type = "File"
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
