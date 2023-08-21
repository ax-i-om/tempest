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

// Package sendvid contains functions that can be used to accurately extract and validate Sendvid links.
package sendvid

import (
	"fmt"
	"regexp"

	"github.com/ax-i-om/vigor/internal/req"
)

// Extract returns a slice of all Sendvid links contained within a string, if any.
func Extract(res string) ([]string, error) {
	// Compile the RegEx expression to be used in the identification and extraction of the Sendvid links
	re := regexp.MustCompile("^(https|http)://sendvid.com/([a-z0-9]{8})")
	// Return all Sendvid links found within an http response
	return re.FindAllString(res, -1), nil
}

// Validate performs a GET request to the Sendvid URL and uses the response status code to identify its validity
func Validate(x string) (bool, error) {
	// Perform a GET request using the Sendvid URL
	res, err := req.GetRes(x)
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
func Delegate(res string) ([]string, error) {
	// Use Extract() to extract any existing Sendvid links from the response
	x, err := Extract(res)
	if err != nil {
		return nil, err
	}
	// Check if the return slice of Sendvid links is empty
	if len(x) > 0 {
		// Create a new, empty slice where we will append any valid Sendvid links
		var results []string = nil
		// Loop through each Sendvid link within the slice
		for _, v := range x {
			// Call the Validate function in order to check whether or not the link is valid
			x, err := Validate(v)
			if err != nil {
				// If any error occurs during the validation process, stop the current iteration and immediately begin with the next link within the slice
				continue
			}
			// If x, the bool return by Validate(), is true: output the result to the terminal and append the link to the specified results slice.
			if x {
				fmt.Println("SENDVID: ", v)
				results = append(results, v)
			}
		}
		// When the loop is finished, return the results slice
		return results, nil
	}
	// Return nothing, if nothing happens (bruh)
	return nil, nil
}
