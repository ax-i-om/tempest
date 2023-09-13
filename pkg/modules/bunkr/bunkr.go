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
	"strings"

	"github.com/ax-i-om/tempest/internal/hdl"
	"github.com/ax-i-om/tempest/internal/models"
	"github.com/ax-i-om/tempest/internal/req"
)

// Compile the RegEx expression to be used in the identification and extraction of the Bunkr links
var bLink *regexp.Regexp = regexp.MustCompile(`(https|http)://bunkrr.su/a/([a-zA-Z0-9]{8})`)

// Compile the RegEx expression for extracting the area that contains the title
var roughTitle *regexp.Regexp = regexp.MustCompile(`<title>(.*?)</title>`)

// Size RegEx expression
var size *regexp.Regexp = regexp.MustCompile(`(\d+(?:\.\d+)?)\s*([KMGTP]?B)`)

// Compile the RegEx expression for extracting a thumbnail
var thumb *regexp.Regexp = regexp.MustCompile(`(https|http)://((i(-pizza|-burger|\d+))|(big-taco-1img)).bunkr.ru/thumbs/(.*?).(png|jpg|jpeg)`)

// Compile the RegEx expression for dirty extraction of size and file count
var roughInfo *regexp.Regexp = regexp.MustCompile(`">\d+files \((\d+(?:\.\d+)?)\s*([KMGTP]?B)\)<\/span`)

// Compile the RegEx expression for stripping right side of info
var iRight *regexp.Regexp = regexp.MustCompile(`files \((\d+(?:\.\d+)?)\s*([KMGTP]?B)\)<\/span`)

// Compile the RegEx expression for extraction of digits
var digits *regexp.Regexp = regexp.MustCompile(`\d+`)

// Extract returns a slice of all Bunkr links contained within a string, if any.
func Extract(res string) ([]string, error) {
	// Return all Bunkr links found within an http response
	return bLink.FindAllString(res, -1), nil
}

// Convert takes a slice of bunkr links that use varying domains and converts them to an active domain that can be used.
func Convert(res string) string {
	post := strings.ReplaceAll(res, "bunkr.is", "bunkrr.su")
	post = strings.ReplaceAll(post, "bunkr.ru", "bunkrr.su")
	post = strings.ReplaceAll(post, "bunkr.su", "bunkrr.su")
	post = strings.ReplaceAll(post, "bunkr.la", "bunkrr.su")
	return post
}

// Validate performs a GET request to the Bunkr URL and uses the response status code to identify its validity
func Validate(x string) (bool, error) {
	// Perform a GET request using the Bunkr URL
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

// FetchViews uses the Bunkr API (https://slut.bunkr.ru/slutsCount?pageUrl=${pageUrl}) to get view count
func FetchViews(albumUrl string) (string, error) {
	// Get body contents of the sendvid link
	res, err := req.GetRes(`https://slut.bunkr.ru/slutsCount?pageUrl=` + albumUrl)
	if err != nil {
		return "", err
	}

	// Read results of the *http.Response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return digits.FindString(string(body)), nil
}

// Delegate takes a string as an argument and returns a slice of valid Senvid links found within the response (if any) and an error
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
				// Get body contents of the sendvid link
				res, err := req.GetRes(v)
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

				// Extract title
				title := roughTitle.FindString(contents)
				title = strings.ReplaceAll(title, `<title>`, ``)
				title = strings.ReplaceAll(title, `</title>`, ``)
				title = strings.ReplaceAll(title, ` | Bunkr`, ``)

				// Extract info
				info := roughInfo.FindString(contents)

				// Extract size from info
				fsize := size.FindString(info)

				// Select albumn thumbnail by extracting thumbnail of first image in album
				rt := thumb.FindString(contents)
				thumbnail := html.UnescapeString(rt)

				// Extract file count from info
				rStrip := iRight.FindString(info)
				cRough := strings.Replace(info, rStrip, ``, 1)
				fcount := strings.ReplaceAll(cRough, `">`, ``)

				// Extract views using API
				views, _ := FetchViews(v)

				// Create type Entry and specify the respective values
				ent := models.Entry{Link: v, Service: "Bunkr", LastValidation: hdl.Time(), Title: title, Size: fsize, FileCount: fcount, Views: views, Thumbnail: thumbnail, Type: "Folder"}

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
