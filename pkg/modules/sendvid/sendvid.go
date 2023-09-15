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

// Package sendvid contains functions that can be used to accurately extract and validate Sendvid links.
package sendvid

import (
	"html"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/ax-i-om/tempest/internal/handlers"
	"github.com/ax-i-om/tempest/internal/models"
)

// Compile RegEx expressions for extraction of links/metadata
var sLink *regexp.Regexp = regexp.MustCompile("(https|http)://sendvid.com/([a-z0-9]{8})")                           // Extract sendvid links
var rThumb *regexp.Regexp = regexp.MustCompile(`(https|http)://thumbs(.*?).jpg`)                                    // Extract thumbnail
var roughViews *regexp.Regexp = regexp.MustCompile(`<p class="hits"><i class="icon-icn-view"></i>([0-9](.*?))</p>`) // Extract view count
var roughTitle *regexp.Regexp = regexp.MustCompile(`<title>(.*?)</title>`)                                          // Extract title

// Extract returns a slice of all Sendvid links contained within a string, if any.
func Extract(res string) ([]string, error) {
	// Return all Sendvid links found within an http response
	return sLink.FindAllString(res, -1), nil
}

// ExtractTitle takes the body response/contents of a Sendvid page (raw source/html (formatted as string)) as
// an argument and returns the title as a string.
func ExtractTitle(sendvidContents string) string {
	title := roughTitle.FindString(sendvidContents)  // Extract rough title
	title = strings.ReplaceAll(title, `<title>`, ``) // Strip opening tags
	return strings.ReplaceAll(title, `</title>`, ``) // Strip closing tags
}

// ExtractThumbnail takes the body response/contents of a Sendvid page (raw source/html (formatted as string)) as
// an argument. The extracted URL is unescaped to ensure validity and returned in string format.
func ExtractThumbnail(sendvidContents string) string {
	return html.UnescapeString(rThumb.FindString(sendvidContents)) // Extract and unescape thumbnail URL
}

// ExtractViewCount takes the body response/contents of a Sendvid page (raw source/html (formatted as string)) as
// an argument and returns a string containing the view count, alongside an error. The error will be nil if
// everything is successful. If a failure occurs, -1 will be returned.
func ExtractViewCount(sendvidContents string) int {
	eViews := roughViews.FindString(sendvidContents)                                         // Extract rough view count
	eViews = strings.ReplaceAll(eViews, `<p class="hits"><i class="icon-icn-view"></i>`, ``) // Strip unnecessary html
	viewcount, err := strconv.Atoi(strings.ReplaceAll(eViews, `</p>`, ``))                   // Strip closing tag and convert to int
	if err != nil {
		return -1 // Return -1 to signify an error occured and the filecount could not be converted to Int
	}
	return viewcount
}

// Validate performs a GET request to the Sendvid URL and uses the response status code to identify its validity
func Validate(x string) (bool, error) {
	// Perform a GET request using the Sendvid URL
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
	// Use Extract() to extract any existing Sendvid links from the response
	x, err := Extract(res)
	if err != nil {
		return nil, err
	}
	// Check if the return slice of Sendvid links is empty
	if len(x) > 0 {
		// Create a new, empty slice where we will append any valid Sendvid links
		var results []models.Entry = nil
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
				aThumbnail := ExtractThumbnail(contents)
				aViewCount := ExtractViewCount(contents)

				// Create type Entry and specify the respective values
				ent := models.Entry{Source: source, Link: v, Service: "Sendvid", Thumbnail: aThumbnail, Views: aViewCount, Title: aTitle, Type: "File"}
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
