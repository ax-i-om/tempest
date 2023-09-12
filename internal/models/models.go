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

// Package models contains type declarations used in Tempest
package models

// Entry
type Entry struct {
	Link           string `json:"link"`
	LastValidation string `json:"lastvalidation"`

	Title       string `json:"title"`
	Description string `json:"description"`
	Service     string `json:"service"`
	Uploaded    string `json:"uploaded"`

	Type      string `json:"type"`
	Size      string `json:"size"`
	Length    string `json:"length"`
	FileCount string `json:"filecount"`

	Thumbnail string `json:"thumbnail"`
	Downloads string `json:"downloads"`
	Views     string `json:"views"`
}
