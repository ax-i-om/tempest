/*
Vigor - Leveraging paste sites as a medium for discovery
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

// Package mega contains functions that can be used to accurately extract and validate Mega links.
package mega

import (
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/ax-i-om/vigor/internal/req"
)

// Extract returns a slice of all Mega links contained within a string, if any.
func Extract(res string) ([]string, error) {
	// Compile the RegEx expression to be used in the identification and extraction of the Mega links
	re := regexp.MustCompile("(https|http)://mega.nz/(folder|file)/([a-zA-Z0-9]{0,8})#([a-zA-Z0-9_-]{43}|[a-zA-Z0-9_-]{22})")
	// Return all Mega links found within an http response
	return re.FindAllString(res, -1), nil
}

// Validate takes a Mega link/URL and passes it to the Mega API to check whether or not it is online.
func Validate(x string) (bool, error) {
	// Compile the RegEx expression to be used in the extraction of the ID
	re := regexp.MustCompile("([a-zA-Z0-9]{8}#)")
	// Extracts the eight-character identifier from the URL using the pound (#) symbol as context.
	pre := re.FindString(x)
	// Remove the pound (#) symbol from the extracted ID
	post := strings.ReplaceAll(pre, "#", "")

	// Append the ID (post) to the URL
	url := "https://g.api.mega.co.nz/cs?id=5644474&n=" + post

	// Perform a GET request using the pre-formatted URL
	res, err := req.GetRes(url)
	if err != nil {
		return false, err
	}

	// Prepare the contents of the response to be read
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return false, err
	}

	// Read the response, if the body response contains a -2, then the Mega link is valid. If it contains a -9 or anything else, it is invalid.
	if strings.Contains(string(body), "-2") {
		return true, res.Body.Close()
	} else {
		return false, res.Body.Close()
	}
}

// Takes a string as an argument and returns a slice of valid Mega links found within the response (if any) and an error
func Delegate(res string) ([]string, error) {
	// Use Extract() to extract any existing Mega links from the *http.Response
	x, err := Extract(res)
	if err != nil {
		return nil, err
	}
	// Check if the return slice of Mega links is empty
	if len(x) > 0 {
		// Create a new, empty slice where we will append any valid Mega links
		var results []string = nil
		// Loop through each Mega link within the slice
		for _, v := range x {
			// Call the Validate function in order to check whether or not the link is valid
			x, err := Validate(v)
			if err != nil {
				// If any error occurs during the validation process, stop the current iteration and immediately begin with the next link within the slice
				continue
			}
			// If x, the bool return by Validate(), is true: output the result to the terminal and append the link to the specified results slice.
			if x {
				fmt.Println("VALID MEGA LINK: ", v)
				results = append(results, v)
			}
		}
		// When the loop is finished, return the results slice
		return results, nil
	}
	// Return nothing, if nothing happens (bruh)
	return nil, nil
}
