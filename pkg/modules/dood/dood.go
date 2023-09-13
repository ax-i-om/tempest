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

// Package dood contains functions that can be used to accurately extract and validate Dood links.
package dood

import (
	"io"
	"regexp"
	"strings"

	"github.com/ax-i-om/tempest/internal/hdl"
	"github.com/ax-i-om/tempest/internal/models"
	"github.com/ax-i-om/tempest/internal/req"
)

// Compile the RegEx expression to be used in the identification and extraction of the Bunkr links
var dLink *regexp.Regexp = regexp.MustCompile("(https|http)://doods.pro/((f/[a-z0-9]{10})|((d/[a-z0-9]{32}|(d/[a-z0-9]{31})|(d/[a-z0-9]{12})))|e/[a-z0-9]{12})")

/*

Cannot read valid dood link body as a 403 is returned upon request. Dood requires the client to have JS enabled. Maybe achievable through headless-browser?

// Compile the RegEx expression for extracting the thumbnail URL
var rThumb *regexp.Regexp = regexp.MustCompile(`(https|http)://img(.*?).jpg`)

// Compile the RegEx expression for extracting the area that contains the title
var roughTitle *regexp.Regexp = regexp.MustCompile(`<title>(.*?)</title>`)

// Compile the RegEx expression for extracting the area that contains the upload date
var roughDate *regexp.Regexp = regexp.MustCompile(`<div class="uploadate"> <i class="fad fa-calendar-alt mr-1"></i> (.*?) </div>`)

// Size RegEx expression
var size *regexp.Regexp = regexp.MustCompile(`(\d+(?:\.\d+)?)\s*([KMGTP]?B)`)

// Compile the RegEx expression for extracting the length
var length *regexp.Regexp = regexp.MustCompile(`(\d+):(\d+)`)

*/

// Extract returns a slice of all Dood links contained within a string, if any.
func Extract(res string) ([]string, error) {
	// Return all Dood links found within an http response
	return dLink.FindAllString(res, -1), nil
}

// Convert takes a slice of Dood links that use varying domains and converts them to an active domain that can be used.
func Convert(res string) string {
	post := strings.ReplaceAll(res, "dood.la", "doods.pro")
	post = strings.ReplaceAll(post, "dood.re", "doods.pro")
	post = strings.ReplaceAll(post, "dood.wf", "doods.pro")
	post = strings.ReplaceAll(post, "dood.yt", "doods.pro")
	post = strings.ReplaceAll(post, "dood.so", "doods.pro")
	return post
}

// Validate performs a GET request to the Dood URL and uses the response status code to identify its validity
func Validate(x string) (bool, error) {
	// Perform a GET request using the Dood URL
	res, err := req.GetRes(x)
	if err != nil {
		return false, err
	}

	// Prepare the contents of the response to be read
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return false, err
	}

	// Read the response, if the title contains the below specified string, then the Dood link is not online.
	if strings.Contains(string(body), "<title>Video not found | DoodStream</title>") {
		return false, res.Body.Close()
	} else {
		return true, res.Body.Close()
	}
}

// Delegate takes a string as an argument and returns a slice of valid Senvid links found within the response (if any) and an error
func Delegate(res string) ([]models.Entry, error) {
	// Use Convert() to convert all Dood link domains to doods.pro (the currently active one)
	c := Convert(res)
	// Use Extract() to extract any existing Dood links from the converted response
	x, err := Extract(c)
	if err != nil {
		return nil, err
	}
	// Check if the return slice of Dood links is empty
	if len(x) > 0 {
		// Create a new, empty slice where we will append any valid Dood links
		var results []models.Entry = nil
		// Loop through each Dood link within the slice
		for _, v := range x {
			// Call the Validate function in order to check whether or not the link is valid
			x, err := Validate(v)
			if err != nil {
				// If any error occurs during the validation process, stop the current iteration and immediately begin with the next link within the slice
				continue
			}
			// If x, the bool return by Validate(), is true: output the result to the terminal and append the link to the specified results slice.
			if x {
				/*

					// Extract title
					rt := roughTitle.FindString(contents)
					r1 := strings.ReplaceAll(rt, `<title>`, ``)
					title := strings.ReplaceAll(r1, `</title>`, ``)

					// Extract length
					flength := length.FindString(contents)

					// Extract thumbnail URL
					thumb := rThumb.FindString(contents)

					// Extract size
					fsize := size.FindString(contents)

					// Extract date
					rdate := roughDate.FindString(contents)
					d1 := strings.ReplaceAll(rdate, `<div class="uploadate"> <i class="fad fa-calendar-alt mr-1"></i> `, ``)
					uploaded := strings.ReplaceAll(d1, ` </div>`, ``)

				*/

				// Create type Entry and specify the respective values
				ent := models.Entry{Link: v, Service: "Dood", LastValidation: hdl.Time(), Type: "File"}
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
