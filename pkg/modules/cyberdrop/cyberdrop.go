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
	"io"
	"regexp"
	"strings"
	"time"

	"github.com/ax-i-om/tempest/internal/hdl"
	"github.com/ax-i-om/tempest/internal/models"
	"github.com/ax-i-om/tempest/internal/req"
)

// Compile the RegEx expression to be used in the identification and extraction of the Cyberdrop links
var cLink *regexp.Regexp = regexp.MustCompile("(https|http)://cyberdrop.me/a/([a-zA-Z0-9]{8})")

// Compile the RegEx expression for extracting the area that contains the title
var roughTitle *regexp.Regexp = regexp.MustCompile(`<title>(.*?)</title>`)

// Compile the RegEx expression for extracting the area that lies to the right of the actual title
var rightTitle *regexp.Regexp = regexp.MustCompile(` \[\d+ files(.*?)\| CyberDrop`)

// Compile the RegEx expression for extracting the area that contains the upload date
var roughUploaded *regexp.Regexp = regexp.MustCompile(`<p class="heading">Uploaded</p>  <p class="title">(.*?)</p>`)

// Compile the RegEx expression for extracting the thumbnail URL
var rThumb *regexp.Regexp = regexp.MustCompile(`(https|http)://i0.wp.com(.*?).png`)

// Compile the RegEx expression for extracting the area that contains the description (specified by author)
var roughDesc *regexp.Regexp = regexp.MustCompile(`\[Reg: CLOSED] - (.*?)" />`)

// Compile the RegEx expression for dirty extraction of file count from title
var rFiles *regexp.Regexp = regexp.MustCompile(`\[(.*?) files ::`)

// Compile the RegEx expression for extraction of digits
var digits *regexp.Regexp = regexp.MustCompile(`\d+`)

// Size RegEx expression
var size *regexp.Regexp = regexp.MustCompile(`(\d+(?:\.\d+)?)\s*([KMGTP]?B)`)

// Extract returns a slice of all Cyberdrop links contained within a string, if any.
func Extract(res string) ([]string, error) {
	// Return all Cyberdrop links found within an http response
	return cLink.FindAllString(res, -1), nil
}

// Validate performs a GET request to the Cyberdrop URL and uses the response status code to identify its validity
func Validate(x string) (bool, error) {
	// Perform a GET request using the Cyberdrop URL
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
func Delegate(res string) ([]models.Entry, error) {
	// Use Extract() to extract any existing Cyberdrop links from the response
	x, err := Extract(res)
	if err != nil {
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
				rt := roughTitle.FindString(contents)
				r1 := strings.ReplaceAll(rt, `<title>`, ``)
				r2 := strings.ReplaceAll(r1, `</title>`, ``)
				r3 := strings.ReplaceAll(r2, `Album: `, ``)
				tStrip := rightTitle.FindString(r3)
				title := strings.ReplaceAll(r3, tStrip, ``)

				// Extract size
				fsize := size.FindString(r3)

				// Extract file count
				rf := rFiles.FindString(r3)
				filecount := digits.FindString(rf)

				// Extract description
				rd := roughDesc.FindString(contents)
				d1 := strings.ReplaceAll(rd, `[Reg: CLOSED] - `, ``)
				desc := strings.ReplaceAll(d1, `" />`, ``)
				if strings.Contains(desc, "A privacy-focused censorship-resistant file sharing platform free for everyone. Upload files up to 200MB. Keep your uploads safe and secure with us") {
					desc = ""
				}

				uploaded := ""

				// Extract upload date
				ru := roughUploaded.FindString(contents)
				u1 := strings.ReplaceAll(ru, `<p class="heading">Uploaded</p>  <p class="title">`, ``)
				u2 := strings.ReplaceAll(u1, `</p>`, ``)
				u3 := strings.ReplaceAll(u2, `.`, `-`)
				u4, err := time.Parse("02-01-2006", u3)
				if err != nil {
					uploaded = u3
				} else {
					uploaded = u4.Format("Jan 02, 2006")
				}

				// Extract thumbnail URL
				thumb := rThumb.FindString(contents)

				// Create type Entry and specify the respective values
				ent := models.Entry{Link: v, Service: "Cyberdrop", LastValidation: hdl.Time(), Thumbnail: thumb, Description: desc, Title: title, FileCount: filecount, Size: fsize, Type: "Folder", Uploaded: uploaded}
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
