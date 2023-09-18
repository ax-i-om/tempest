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

// Package cyberdrop contains functions that can be used to accurately extract and validate Cyberdrop links.
package cyberdrop

import (
	"html"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ax-i-om/tempest/internal/handlers"
	"github.com/ax-i-om/tempest/pkg/models"
)

// Compile RegEx expressions for extraction of links/metadata
var cLink *regexp.Regexp = regexp.MustCompile("(https|http)://cyberdrop.me/a/([a-zA-Z0-9]{8})")                      // Extract Cyberdrop link
var roughTitle *regexp.Regexp = regexp.MustCompile(`<title>(.*?)</title>`)                                           // Extract title area/node
var rightTitle *regexp.Regexp = regexp.MustCompile(` \[\d+ files(.*?)\| CyberDrop`)                                  // Extract rightmost area of title
var roughUploaded *regexp.Regexp = regexp.MustCompile(`<p class="heading">Uploaded</p>  <p class="title">(.*?)</p>`) // Extract date area
var thumb *regexp.Regexp = regexp.MustCompile(`(https|http)://i0.wp.com(.*?).png`)                                   // Thumbnail extraction
var roughDesc *regexp.Regexp = regexp.MustCompile(`\[Reg: CLOSED] - (.*?)" />`)                                      // Area that contains description
var rFiles *regexp.Regexp = regexp.MustCompile(`\[(.*?) files ::`)                                                   // Rough file count extraction
var digits *regexp.Regexp = regexp.MustCompile(`\d+`)                                                                // Digit extraction
var size *regexp.Regexp = regexp.MustCompile(`(\d+(?:\.\d+)?)\s*([KMGTP]?B)`)                                        // Size extraction

// Extract returns a slice of all Cyberdrop links contained within a string, if any.
func Extract(res string) ([]string, error) {
	return cLink.FindAllString(res, -1), nil // Return all Cyberdrop links found within an http response
}

// ExtractTitle takes the body response/contents of a Cyberdrop page (raw source/html (formatted as string)) as
// an argument and returns the album's title as a string.
func ExtractTitle(doodContents string) string {
	eTitle := roughTitle.FindString(doodContents)                          // Extract rough title
	eTitle = strings.ReplaceAll(eTitle, `<title>`, ``)                     // Strip opening title tag
	eTitle = strings.ReplaceAll(eTitle, `</title>`, ``)                    // Strip closing title tag
	eTitle = strings.ReplaceAll(eTitle, `Album: `, ``)                     // Strip unnecessary text
	stripExtra := rightTitle.FindString(eTitle)                            // Identify extra information appended to title
	return html.UnescapeString(strings.ReplaceAll(eTitle, stripExtra, ``)) // Strip extra information identified via RegEx and unescape
}

// ExtractSize takes the body response/contents of a Cyberdrop page (raw source/html (formatted as string)) as
// an argument and returns the album's total/cumulative file size as a string.
func ExtractSize(doodContents string) string {
	eTitle := roughTitle.FindString(doodContents) // Extract rough title
	return size.FindString(eTitle)                // Find size information identified via RegEx and return results
}

// ExtractFileCount takes the body response/contents of a Cyberdrop page (raw source/html (formatted as string)) as
// an argument and returns the album's file count as an integer. It will return -1 in the case of a syntax error.
func ExtractFileCount(doodContents string) int {
	eTitle := roughTitle.FindString(doodContents)   // Extract rough title
	eCount := rFiles.FindString(eTitle)             // Extract rough file count
	c, e := strconv.Atoi(digits.FindString(eCount)) // Find file count information identified via RegEx and convert to int
	if e != nil {                                   // Error (likely syntax) returns a -1 instead
		c = -1
	}
	return c
}

// ExtractDescription takes the body response/contents of a Cyberdrop page (raw source/html (formatted as string)) and
// the excludeDefault flag (bool). The exclude default flag will return an empty string in the case that the album
// creator did not specify a description, rather than Cyberdrop's ~150 character default description. This should be
// set to true in most cases. If the album creator did specify a description, it will be returned as a string.
func ExtractDescription(doodContents string, excludeDefault bool) string {
	eDesc := roughDesc.FindString(doodContents)               // Extract rough description
	eDesc = strings.ReplaceAll(eDesc, `[Reg: CLOSED] - `, ``) // Strip unnecessary text
	eDesc = strings.ReplaceAll(eDesc, `" />`, ``)             // Strip unnecessary text
	// If default description is returned (in the case the album creator did not specify a description)
	// and the excludeDefault flag is set to true, return an empty string
	if excludeDefault && strings.Contains(eDesc, "A privacy-focused censorship-resistant file sharing platform free for everyone. Upload files up to 200MB. Keep your uploads safe and secure with us") {
		return ""
	}
	return html.UnescapeString(eDesc) // Otherwise, return the unescaped description
}

// ExtractThumbnail takes the body response/contents of a Cyberdrop page (raw source/html (formatted as string)) as
// an argument. Cyberdrop albums do not have a dedicated thumbnail, so ExtractThumbnail() instead extracts the first
// image URL it finds, as this can grant more insight if other metadata is misleading/inconclusive. The extracted URL
// is unescaped to ensure validity and returned in string format.
func ExtractThumbnail(doodContents string) string {
	return html.UnescapeString(thumb.FindString(doodContents)) // Find picture via RegEx, and return the unescaped URL
}

// ExtractUploadDate takes the body response/contents of a Cyberdrop page (raw source/html (formatted as string)) as
// and a dateFormat string as it's arguments. The dateFormat string will be used to format a time.Time into a more
// favorable format. The default dateFormat used by tempest is "Jan 02, 2006". It will return a date in string format,
// which has been formatted in accordance with the specified dateFormat, and an error which will be nil if successful.
func ExtractUploadDate(doodContents string, dateFormat string) (string, error) {
	eUdate := roughUploaded.FindString(doodContents)                                              // Extract rough upload date
	eUdate = strings.ReplaceAll(eUdate, `<p class="heading">Uploaded</p>  <p class="title">`, ``) // Strip tags and attributes
	eUdate = strings.ReplaceAll(eUdate, `</p>`, ``)                                               // Strip closing paragraph tag
	eUdate = strings.ReplaceAll(eUdate, `.`, `-`)                                                 // Convert any periods (.) to dashes/hyphens (-)
	parsed, err := time.Parse("02-01-2006", eUdate)                                               // Convert string to *time.Time
	if err != nil {
		return "", err
	}
	return parsed.Format(dateFormat), nil // Format the date based on the format specified in dateFormat and return as string
}

// Validate performs a GET request to the Cyberdrop URL and uses the response status code to identify its validity
func Validate(x string) (bool, error) {
	// Perform a GET request using the Cyberdrop URL
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

// Delegate takes a string as an argument and returns a slice of valid Cyberdrop links found within the response (if any) or nil, and an error
func Delegate(res, source string) ([]models.Entry, error) {
	// Use Extract() to extract any existing Cyberdrop links from the response
	x, err := Extract(res)
	if err != nil {
		handlers.LogErr(err, "error occurred on cyberdrop delegate attempt to call extract")
		return nil, err
	}
	// Check if the return slice of Cyberdrop links is empty
	if len(x) > 0 {
		// Create a new, empty slice where we will append any valid Cyberdrop links
		var results []models.Entry = nil
		// Loop through each Cyberdrop link within the slice
		for _, v := range x {
			// Call the Validate function in order to check whether or not the link is valid
			x, err := Validate(v)
			if err != nil {
				// If any error occurs during the validation process, stop the current iteration and immediately begin with the next link within the slice
				handlers.LogErr(err, "error occurred on cyberdrop delegate attempt to call validate")
				continue
			}
			// If x, the bool return by Validate(), is true: output the result to the terminal and append the link to the specified results slice.
			if x {
				// Get body contents of the cyberdrop link
				res, err := handlers.GetRes(v)
				if err != nil {
					handlers.LogErr(err, "error occurred on cyberdrop delegate attempt to call getres")
					continue
				}

				// Read results of the *http.Response body
				body, err := io.ReadAll(res.Body)
				if err != nil {
					handlers.LogErr(err, "error occurred on cyberdrop delegate attempt to call readall")
					continue
				}

				// Convert read results to a string
				contents := string(body)

				// Remove newline and tabs from content
				contents = strings.ReplaceAll(contents, "\n", ``)
				contents = strings.ReplaceAll(contents, "\t", ``)

				aTitle := ExtractTitle(contents)                              // Extract title
				aSize := ExtractSize(contents)                                // Extract size
				aCount := ExtractFileCount(contents)                          // Extract file count
				aDescription := ExtractDescription(contents, true)            // Extract description
				aThumbnail := ExtractThumbnail(contents)                      // Extract thumbnail
				aUploadDate, _ := ExtractUploadDate(contents, "Jan 02, 2006") // Extract upload date

				// Create type Entry and specify the respective values
				ent := models.Entry{Source: source, Link: v, Service: "Cyberdrop", Thumbnail: aThumbnail, Description: aDescription, Title: aTitle, FileCount: aCount, Size: aSize, Type: "Folder", Uploaded: aUploadDate}
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
