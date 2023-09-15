<p align="center">
  <a><img src="./images/icon.png" width=280 height="245"></a>
    <h1 align="center">Tempest</h1>
  <p align="center">
    <a href="https://app.deepsource.com/gh/ax-i-om/tempest/" target="_blank"><img alt="DeepSource" title="DeepSource" src="https://app.deepsource.com/gh/ax-i-om/tempest.svg/?label=active+issues&show_trend=true"/></a>
    <a href="https://pkg.go.dev/github.com/ax-i-om/tempest#section-directories"><img src="https://pkg.go.dev/badge/github.com/ax-i-om/tempest.svg" alt=""></a>
    <a><img src="https://img.shields.io/badge/version-0.5.0-blue.svg" alt="v0.5.0"></a>
    <a href="https://goreportcard.com/report/github.com/ax-i-om/tempest"><img src="https://goreportcard.com/badge/github.com/ax-i-om/tempest" alt="Go Report Card"></a><br>
   Leverage paste sites as a medium for discovery of objectionable/infringing materials. <br>
</a>
  </p><br>
</p>

## Table of Contents

- [Information](#information)
  - [About](#about)
  - [Disclaimer](#disclaimer)
  - [Installation](#installation)
  - [Usage](#usage)
  - [Cloud Storage / File Sharing Platform Modules](#cloud-storage--file-sharing-platform-modules)
  - [Entry Format](#entry-format)
  - [Important Notes](#important-notes)

## Information

### About

Tempest is a simple, lightweight, and cross-platform solution designed to enable individuals to efficiently discover and extract active cloud storage/file sharing links from paste platforms such as [Rentry.co](https://rentry.co). It was created to address the notable uptick in paste sites being used to distribute content that violates copyright and piracy statutes.

### Disclaimer

It is the end user's responsibility to obey all applicable local, state, and federal laws. Developers assume no liability and are not responsible for any misuse or damage caused by this program. By using Tempest, you agree to the previous statements.

### Installation
1. Fetch the repository via ***git clone***: `git clone https://github.com/ax-i-om/tempest.git`
2. Navigate to the root directory of of the cloned repository via ***cd***: `cd tempest`
3. In your preferred terminal, enter and run: `go run main.go`

### Usage

Display Tempest usage help in the terminal via: `go run main.go help`

Tempest supports three primary methods of output, those being JSON, CSV, and plain text (output to console). 
If you want to output plain text to the console, run tempest like so: `go run main.go console`.

*Note:*
If you want to output the console results to a file, append this to the command: `2>&1 | tee results.txt` <br>
For example `go run main.go console 2>&1 | tee results.txt` (may vary depending on operating system) <br>
CAUTION: IF THE SPECIFIED OUTPUT FILE ALREADY EXISTS, THIS WILL OVERWRITE THE CONTENTS

If you want to output the results to a JSON/CSV file, the command should be formatted like so: `go run main.go <json/csv> <filename>`<br>
JSON Example: `go run main.go json results` ***VS*** CSV Example: `go run main.go csv results`<br>
*Note:* If you exclude the file extension *(.json/.csv)*, one will be automatically appended.

In order to gracefully shut down Tempest, press `Ctrl + C` in the terminal **ONCE** and wait until the remaining goroutines finish executing (typically <60s).<br>
In order to forcefully shut down Tempest press `Ctrl + C` in the terminal **TWICE**.<br>
*CAUTION:* FORCEFULLY SHUTTING DOWN TEMPEST MAY RESULT IN ISSUES INCLUDING, BUT NOT LIMITED TO, DATA LOSS AND FILE CORRUPTION.

If you decide to output the results to a JSON file specifically, it will not be valid JSON.<br>
Tempest comes bundled with a function for cleaning the resulting JSON content and can be used like so: `go run main.go clean results.json`<br>
This will be the quickest way of converting the JSON file formatting into one that is valid; however, reusing this file for results will cause further formatting issues. The clean function also removes any duplicate entries from the file.
*Note:* If you exclude the file extension *(.json/.csv)*, one will be automatically appended.

### Cloud Storage / File Sharing Platform Modules

| Module        | Status       | Information Extracted                                                            |
| :-----------: | ------------ | :------------------------------------------------------------------------------: |
| Bunkr         | Functioning  | Link, Title, Service, Type, Size, FileCount, Thumbnail, Views                    |
| Cyberdrop     | Functioning  | Link, Title, Service, Type, Size, FileCount, Thumbnail, Description, UploadDate  |
| Dood          | Functioning  | Link, Service, Type                                                              |
| Gofile        | Functioning  | Link, Title, Service, Type, FileCount, Downloads                                 |
| Google Drive  | Functioning  | Link, Title, Service, Type                                                       |
| Mega          | Functioning  | Link, Service, Type, Size, FileCount                                             |
| Sendvid       | Functioning  | Link, Title, Service, Type, Thumbnail, Views                                     |

### Entry Format

``` go
type Entry struct {
	Source string `json:"source"`
	Link   string `json:"link"`

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
```

### Important Notes

- Mega file count and size is unreliable, as the metadata specified in the Mega folder/file headers doesn't seem to accurately align with the true content's file count/size. Take with a grain of salt.
- CSV values are delimited with commas (,). Ensure that fields are not separated via other characters/delimeters such as semicolons (;) and tabs as this may cause presentation/formatting issues.