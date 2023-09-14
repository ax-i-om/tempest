<p align="center">
  <a><img src="./images/icon.png" width=280 height="245"></a>
    <h1 align="center">Tempest</h1>
  <p align="center">
    <a href="https://goreportcard.com/report/github.com/ax-i-om/tempest"><img src="https://goreportcard.com/badge/github.com/ax-i-om/tempest" alt="Go Report Card"></a>
    <a><img src="https://img.shields.io/badge/version-0.3.2-blue.svg" alt="v0.3.2"></a><br>
    Leveraging paste sites as a medium for discovery<br>
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
  - [Notes](#notes)

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
This will be the quickest way of converting the JSON file formatting into one that is valid; however, reusing this file for results will cause further formatting issues.
*Note:* If you exclude the file extension *(.json/.csv)*, one will be automatically appended.

### Cloud Storage / File Sharing Platform Modules

| Module        | Status       |
| :-----------: | -------------|
| Bunkr         | Functioning  |
| Cyberdrop     | Functioning  |
| Dood          | Functioning  |
| Gofile        | Functioning  |
| Google Drive  | Functioning  |
| Mega          | Functioning  |
| Sendvid       | Functioning  |

### Entry Format

``` go
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
```

### Notes

- Mega file count and size is unreliable, as the metadata specified in the Mega folder/file headers doesn't seem to accurately align with the true content's file count/size. Take with a grain of salt.
