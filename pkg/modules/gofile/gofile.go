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
	"strings"

	"github.com/ax-i-om/tempest/internal/hdl"
	"github.com/ax-i-om/tempest/internal/models"
	"github.com/ax-i-om/tempest/internal/req"
)

// Compile the RegEx expression to be used in the identification and extraction of the Gofile links
var gLink *regexp.Regexp = regexp.MustCompile("(https|http)://gofile.io/d/([a-zA-Z0-9]{6})")

// Compile the RegEx expression for extracting the area that contains the title
var roughTitle *regexp.Regexp = regexp.MustCompile(`<title>(.*?)</title>`)

// Compile the RegEx expression for extracting the area that contains the description (file count if folder, download count if file)
var roughDesc *regexp.Regexp = regexp.MustCompile(`<meta name='description' content='(.*?)' />`)

// Extract returns a slice of all Gofile links contained within a string, if any.
func Extract(res string) ([]string, error) {
	// Return all Gofile links found within an http response
	return gLink.FindAllString(res, -1), nil
}

// Validate takes a Gofile link/URL and checks certain metadata patterns to identify whether or not the link is valid/online.
func Validate(x string) (bool, string, error) {
	// Perform a GET request using the Gofile URL
	res, err := req.GetRes(x)
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
func Delegate(res string) ([]models.Entry, error) {
	// Use Extract() to extract any existing Gofile links from the response
	x, err := Extract(res)
	if err != nil {
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
				continue
			}
			// If x, the bool return by Validate(), is true: output the result to the terminal and append the link to the specified results slice.
			if x {
				// Extract title
				rt := roughTitle.FindString(contents)
				r1 := strings.ReplaceAll(rt, `<title>`, ``)
				title := strings.ReplaceAll(r1, `</title>`, ``)

				filecount := ""
				downloads := ""
				fTyp := ""

				rd := roughDesc.FindString(contents)
				d1 := strings.ReplaceAll(rd, `<meta name='description' content='`, ``)
				d2 := strings.ReplaceAll(d1, `' />`, ``)

				if strings.Contains(title, "Folder") {
					title = strings.ReplaceAll(title, `Folder `, ``)
					filecount = strings.ReplaceAll(d2, ` files`, ``)
					fTyp = "Folder"
				} else {
					downloads = strings.ReplaceAll(d2, ` downloads`, ``)
					fTyp = "File"
				}

				// Create type Entry and specify the respective values
				ent := models.Entry{Link: v, Service: "GoFile", LastValidation: hdl.Time(), Title: title, FileCount: filecount, Downloads: downloads, Type: fTyp}
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
