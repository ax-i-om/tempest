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

// Package models contains type declarations and associated functions used in Tempest
package models

import (
	"sync"
	"sync/atomic"
)

// Entry represents the extracted link and it's accompanying data.
type Entry struct {
	Link string `json:"link"`

	Title       string `json:"title"`
	Description string `json:"description"`
	Service     string `json:"service"`
	Uploaded    string `json:"uploaded"`

	Type      string `json:"type"`
	Size      string `json:"size"`
	Length    string `json:"length"`
	FileCount int    `json:"filecount"`

	Thumbnail string `json:"thumbnail"`
	Downloads int    `json:"downloads"`
	Views     int    `json:"views"`
}

// WaitGroupCount represents a countable sync.WaitGroup
type WaitGroupCount struct {
	sync.WaitGroup
	count int64
}

// Add ...
func (wg *WaitGroupCount) Add(delta int) {
	atomic.AddInt64(&wg.count, int64(delta))
	wg.WaitGroup.Add(delta)
}

// Done ...
func (wg *WaitGroupCount) Done() {
	atomic.AddInt64(&wg.count, -1)
	wg.WaitGroup.Done()
}

// GetCount ...
func (wg *WaitGroupCount) GetCount() int {
	return int(atomic.LoadInt64(&wg.count))
}
