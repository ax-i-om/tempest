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

	"github.com/ax-i-om/tempest/internal/handlers"
	"github.com/ax-i-om/tempest/pkg/models"
)

// Compile the RegEx expression to be used in the identification and extraction of the Bunkr links
var dLink *regexp.Regexp = regexp.MustCompile("(https|http)://(doods|dood).(la|re|wf|so|yt|pm|sh|to|ws|one|watch|pro|stream)/((f/[a-z0-9]{10})|((d/[a-z0-9]{32}|(d/[a-z0-9]{31})|(d/[a-z0-9]{12})))|e/[a-z0-9]{12})")

// Extract returns a slice of all Dood links contained within a string, if any.
func Extract(res string) ([]string, error) {
	// Return all Dood links found within an http response
	return dLink.FindAllString(res, -1), nil
}

// Validate performs a GET request to the Dood URL and uses the response status code to identify its validity
func Validate(x string) (bool, error) {
	// Perform a GET request using the Dood URL
	res, err := handlers.GetRes(x)
	if err != nil {
		return false, err
	}

	// Prepare the contents of the response to be read
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return false, err
	}

	// Read the response, if the title contains the below specified string, then the Dood link is not online.
	if !strings.Contains(string(body), "<title>Video not found | DoodStream</title>") || res.StatusCode == 403 {
		return true, res.Body.Close()
	} else {
		return false, res.Body.Close()
	}
}

// Convert takes a slice of Dood links that use varying domains and converts them to an active domain that can be used.
func Convert(res string) string {
	post := strings.ReplaceAll(res, "dood.la", "doods.pro")
	post = strings.ReplaceAll(post, "dood.re", "doods.pro")
	post = strings.ReplaceAll(post, "dood.wf", "doods.pro")
	post = strings.ReplaceAll(post, "dood.yt", "doods.pro")
	post = strings.ReplaceAll(post, "dood.so", "doods.pro")
	post = strings.ReplaceAll(post, "dood.pm", "doods.pro")
	post = strings.ReplaceAll(post, "dood.sh", "doods.pro")
	post = strings.ReplaceAll(post, "dood.to", "doods.pro")
	post = strings.ReplaceAll(post, "dood.ws", "doods.pro")
	post = strings.ReplaceAll(post, "dood.one", "doods.pro")
	post = strings.ReplaceAll(post, "dood.watch", "doods.pro")
	return post
}

// Delegate takes a string as an argument and returns a slice of valid Senvid links found within the response (if any) or nil, and an error
func Delegate(res, source string) ([]models.Entry, error) {
	// Use Convert() to convert all Dood link domains to doods.pro (the currently active one)
	c := Convert(res)
	// Use Extract() to extract any existing Dood links from the converted response
	x, err := Extract(c)
	if err != nil {
		handlers.LogErr(err, "error occurred on dood delegate attempt to call extract")
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
				handlers.LogErr(err, "error occurred on dood delegate attempt to call validate")
				continue
			}
			// If x, the bool return by Validate(), is true: output the result to the terminal and append the link to the specified results slice.
			if x {
				// Create type Entry and specify the respective values
				ent := models.Entry{Source: source, Link: v, Service: "Dood", Type: "File"}
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
