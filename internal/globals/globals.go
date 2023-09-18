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
	"github.com/rs/zerolog"
)

// DebugFlag stores information about whether or not to print debug information
var DebugFlag bool

// Mode store the specified output mode
var Mode string

// Filename stores the output file name
var Filename string

// Wg is a sync.WaitGroup that implements a counter, used in run() for graceful cleanup
var Wg models.WaitGroupCount = models.WaitGroupCount{}

// WriteMutex provides the ability to lock write()
var WriteMutex sync.Mutex

// Jsonfile is the opened JSON file where the results are written
var Jsonfile *os.File = nil

// Csvfile is the opened CSV file where results are written
var Csvfile *os.File = nil

// Writer is used in writing to Csvfile
var Writer *csv.Writer = nil

// Src is used by the TrueRand function
var Src = rand.NewSource(time.Now().UnixNano())

// Logger is a global zerolog logger used for printing & debugging
var Logger zerolog.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
	Level(zerolog.TraceLevel).
	With().
	Timestamp().
	Caller().
	Logger()
