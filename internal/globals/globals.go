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
// package globals contains global variables used throughout Tempest
package globals

import (
	"encoding/csv"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/ax-i-om/tempest/pkg/models"
)

// Stores output mode and filename
var Mode, Filename string

// Checks if a CSV file with filename: filename already exists to determine whether or not to write headers
var Existed bool

// Custom sync.WaitGroup that implements a counter, used in run() for graceful cleanup
var Wg models.WaitGroupCount = models.WaitGroupCount{}

// Mutex lock for write()
var WriteMutex sync.Mutex

// Global declaration of files/writers in order to write/flush from anywhere in main
var Jsonfile *os.File = nil
var Csvfile *os.File = nil
var Writer *csv.Writer = nil

var Src = rand.NewSource(time.Now().UnixNano())
