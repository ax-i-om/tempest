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

// Package bunkr contains functions that can be used to accurately extract and validate Bunkr links.
package bunkr

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
var bLink *regexp.Regexp = regexp.MustCompile(`(https|http)://bunkrr.su/a/([a-zA-Z0-9]{8})`)                                                   // Bunkr links
var roughTitle *regexp.Regexp = regexp.MustCompile(`<title>(.*?)</title>`)                                                                     // Rough title
var size *regexp.Regexp = regexp.MustCompile(`(\d+(?:\.\d+)?)\s*([KMGTP]?B)`)                                                                  // Size
var thumb *regexp.Regexp = regexp.MustCompile(`(https|http)://((i(-pizza|-burger|\d+))|(big-taco-1img)).bunkr.ru/thumbs/(.*?).(png|jpg|jpeg)`) // Thumbnail
var roughInfo *regexp.Regexp = regexp.MustCompile(`">\d+files \((\d+(?:\.\d+)?)\s*([KMGTP]?B)\)<\/span`)                                       // Rough info
var iRight *regexp.Regexp = regexp.MustCompile(`files \((\d+(?:\.\d+)?)\s*([KMGTP]?B)\)<\/span`)                                               // Right side info
var digits *regexp.Regexp = regexp.MustCompile(`\d+`)                                                                                          // Digits

// Extract returns a slice of all Bunkr links contained within a string, if any.
func Extract(res string) ([]string, error) {
	return bLink.FindAllString(res, -1), nil // Return all Bunkr links found within an http response
}

// Convert takes a bunkr link in string format as an argument and converts the domain to one
// that is currently active (bunkrr.su right now). Result returned in string format
func Convert(res string) string {
	post := strings.ReplaceAll(res, "bunkr.is", "bunkrr.su")
	post = strings.ReplaceAll(post, "bunkr.ru", "bunkrr.su")
	post = strings.ReplaceAll(post, "bunkr.su", "bunkrr.su")
	post = strings.ReplaceAll(post, "bunkr.la", "bunkrr.su")
	return post
}

// ExtractTitle takes the body response/contents of a Bunkr page (raw source/html (formatted as string)) as
// an argument and returns the album's title as a string.
func ExtractTitle(bunkrContents string) string {
	eTitle := roughTitle.FindString(bunkrContents)      // Extract rough title via RegEx
	eTitle = strings.ReplaceAll(eTitle, `<title>`, ``)  // Strip opening tag
	eTitle = strings.ReplaceAll(eTitle, `</title>`, ``) // Strip closing tag
	return strings.ReplaceAll(eTitle, ` | Bunkr`, ``)   // Strip unneccesary text
}

// ExtractSize takes the body response/contents of a Bunkr page (raw source/html (formatted as string)) as
// an argument and returns the album's total/cumulative file size as a string.
func ExtractSize(bunkrContents string) string {
	return size.FindString(roughInfo.FindString(bunkrContents)) // Find size information within info via RegEx
}

// ExtractThumbnail takes the body response/contents of a Bunkr page (raw source/html (formatted as string)) as
// an argument. Bunkr albums do not have a dedicated thumbnail, so ExtractThumbnail() instead extracts the first
// image URL it finds, as this can grant more insight if other metadata is misleading/inconclusive. The extracted URL
// is unescaped to ensure validity and returned in string format.
func ExtractThumbnail(bunkrContents string) string {
	return html.UnescapeString(thumb.FindString(bunkrContents)) // Find picture via RegEx, and return the unescaped URL
}

// ExtractFileCount takes the body response/contents of a Bunkr page (raw source/html (formatted as string)) as
// an argument and returns the album's file count as an integer. It will return -1 in the case of a syntax error.
func ExtractFileCount(bunkrContents string) int {
	eRough := roughInfo.FindString(bunkrContents)   // Extract rough info
	eStrip := iRight.FindString(eRough)             // Find right side of string via RegEx
	eRough = strings.Replace(eRough, eStrip, ``, 1) // Strip right side of string
	eRough = strings.ReplaceAll(eRough, `">`, ``)   // Strip remnants of html
	c, e := strconv.Atoi(digits.FindString(eRough)) // Extract digits via RegEx and convert to int
	if e != nil {                                   // Error (likely syntax) returns a -1 instead
		c = -1
	}
	return c
}

// ExtractViewCount uses the Bunkr API (https://slut.bunkr.ru/slutsCount?pageUrl=${pageUrl}) to get view count. It
// takes a bunkr album link as an argument (in string format) and returns a string containing the view count,
// alongside an error. The error will be nil if everything is successful. If a failure occurs, -1 will be returned.
func ExtractViewCount(albumUrl string) (int, error) {
	// Get body contents of the bunkr link
	res, err := handlers.GetRes(`https://slut.bunkr.ru/slutsCount?pageUrl=` + albumUrl)
	if err != nil {
		return -1, err
	}

	// Read results of the *http.Response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return -1, err
	}

	views, err := strconv.Atoi(digits.FindString(string(body)))
	if err != nil {
		return -1, err
	}

	return views, nil
}

// Validate performs a GET request to the Bunkr URL and uses the response status code to identify its validity
// If the link is valid, it will return true. If not, it will return false. Validate also returns an error.
func Validate(x string) (bool, error) {
	// Perform a GET request using the Bunkr URL
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

// Delegate takes a string as an argument and returns a slice of valid Bunkr links found within the response (if any) or nil, and an error
func Delegate(res string) ([]models.Entry, error) {
	// Use Convert() to convert all bunkr link domains to bunkrr.su (the currently active one)
	c := Convert(res)
	// Use Extract() to extract any existing Bunkr links from the converted response
	x, err := Extract(c)
	if err != nil {
		return nil, err
	}
	// Check if the return slice of Bunkr links is empty
	if len(x) > 0 {
		// Create a new, empty slice where we will append any valid Bunkr links
		var results []models.Entry = nil
		// Loop through each Bunkr link within the slice
		for _, v := range x {
			// Call the Validate function in order to check whether or not the link is valid
			x, err := Validate(v)
			if err != nil {
				// If any error occurs during the validation process, stop the current iteration and immediately begin with the next link within the slice
				continue
			}
			// If x, the bool return by Validate(), is true: output the result to the terminal and append the link to the specified results slice.
			if x {
				// Get body contents of the bunkr link
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

				// Remove newline and tabs from content
				contents = strings.ReplaceAll(contents, "\n", ``)
				contents = strings.ReplaceAll(contents, "\t", ``)

				aTitle := ExtractTitle(contents)         // Extract title
				aSize := ExtractSize(contents)           // Extract size
				aThumbnail := ExtractThumbnail(contents) // Extract thumbnail
				aFileCount := ExtractFileCount(contents) // Extract file count
				aViews, _ := ExtractViewCount(v)         // Extract view count

				// Create type Entry and specify the respective values
				ent := models.Entry{Link: v, Service: "Bunkr", Title: aTitle, Size: aSize, FileCount: aFileCount, Views: aViews, Thumbnail: aThumbnail, Type: "Folder"}

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
