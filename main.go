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

package main

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/ax-i-om/vigor/internal/req"
	"github.com/ax-i-om/vigor/pkg/modules/gofile"
	"github.com/ax-i-om/vigor/pkg/modules/mega"
)

var src = rand.NewSource(time.Now().UnixNano())

const (
	leIndexBits = 6
	leIndexMask = 1<<leIndexBits - 1
	leIndexMax  = 63 / leIndexBits
)

func trueRand(n int, chars string) string {
	b := make([]byte, n)
	for i, cache, remain := n-1, src.Int63(), leIndexMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), leIndexMax
		}
		if idx := int(cache & leIndexMask); idx < len(chars) {
			b[i] = chars[idx]
			i--
		}
		cache >>= leIndexBits
		remain--
	}

	return string(b)
}

func runner(renturl string) error {
	// Performs a get request on the randomly generated Rentry.co URL.
	res, err := req.GetRes(renturl)
	if err != nil {
		return err
	}
	// If a Status Code of 200 is returned, that means we randomly generated a valid Rentry.co link and can continue
	if res.StatusCode == 200 {
		// Prepare the contents of the response to be read
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		// Convert the slice of bytes to a string
		conv := string(body)

		// Delegate the string to all specified modules
		// Mega Module
		_, err = mega.Delegate(conv)
		if err != nil {
			return err
		}
		// Gofile Module
		_, err = gofile.Delegate(conv)
		if err != nil {
			return err
		}
	}
	return res.Body.Close()
}

func main() {
	// Printing license information to the terminal
	fmt.Println("Vigor Copyright (C) 2023 Axiom\nThis program comes with ABSOLUTELY NO WARRANTY.\nThis is free software, and you are welcome to redistribute it\nunder certain conditions")

	for {
		time.Sleep(1 * time.Millisecond) // Sleeps for 1 millisecond (lol)
		go func() {
			err := runner("https://rentry.co/" + trueRand(5, "abcdefghijklmnopqrstuvwxyz0123456789") + "/raw")
			if err != nil {
				e := err.Error()
				if !strings.Contains(e, "Get") && !strings.Contains(e, "EOF") {
					if strings.Contains(e, "connection reset by peer") || strings.Contains(e, "client connection force closed via ClientConn.Close") || strings.Contains(e, "closed") {
						// SWAP HERE
						// Ignore this
					} else {
						fmt.Println(e)
						os.Exit(1)
					}
				}
			}
		}()
	}
}
