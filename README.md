<p align="center">
  <a><img src="./images/icon.png" width=280 height="245"></a>
    <h1 align="center">Tempest</h1>
  <p align="center">
    <a href="https://goreportcard.com/report/github.com/ax-i-om/tempest"><img src="https://goreportcard.com/badge/github.com/ax-i-om/tempest" alt="Go Report Card"></a>
    <a><img src="https://img.shields.io/badge/version-0.2.0-blue.svg" alt="v0.2.0"></a><br>
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

Tempest supports three primary methods of output, those being JSON, CSV, and plain text (output to console). 
If you want to output plain text to the console, run tempest like so: `go run main.go console`.

If you want to output the results to a JSON/CSV file, the command should be formatted like so: `go run main.go <filetype> <filename>`<br>
JSON Example: `go run main.go json results` ***VS*** CSV Example: `go run main.go csv results`<br>
*Note:* If you exclude the file extension *(.json/.csv), one will be automatically appended.*

*Note:*
If you want to output the console results to a file, append this to the command: `2>&1 | tee results.txt` <br>
For example `go run main.go console 2>&1 | tee results.txt` (may vary depending on operating system) <br>
CAUTION: IF THE SPECIFIED OUTPUT FILE ALREADY EXISTS, THIS WILL OVERWRITE THE CONTENTS

If you decide to output the results to a JSON file specifically, it will not be valid JSON.<br>
Tempest comes bundled with a function for cleaning the resulting JSON content and can be used like so: `go run main.go clean results.json`<br>
This will be the quickest way of converting the JSON file formatting into one that is valid; however, reusing this file for results will cause further formatting issues.

### Cloud Storage / File Sharing Platform Modules

| Module    | Expression                                                                                    |   Validation Method   | Domain Variations? | Status |
| :-------: | --------------------------------------------------------------------------------------------- | :------: | -------- | :----: |
| Bunkr      | (https\|http)://bunkrr.su/a/([a-zA-Z0-9]{8}) |  Status Code  | Yes       | Functioning | 
| Cyberdrop      | (https\|http)://cyberdrop.me/a/([a-zA-Z0-9]{8}) |  Status Code  | No       | Functioning | 
| Dood | (https\|http)://doods.pro/((f/[a-z0-9]{10})\|((d/[a-z0-9]{32}\|(d/[a-z0-9]{31})\|(d/[a-z0-9]{12})))\|e/[a-z0-9]{12}) | Body Contents | Yes | Functioning | 
| Gofile      | (https\|http)://gofile.io/d/([a-zA-Z0-9]{6}) |  Body Contents  | No       | Functioning | 
| Google Drive | (https\|http)://drive.google.com/(folder\|file\|drive)/(d\|folders)/(1[a-zA-Z0-9_-]{32}\|0[a-zA-Z0-9_-]{27}) | Status Code | No | Functioning |
| Mega      | (https\|http)://mega.nz/(folder\|file)/([a-zA-Z0-9]{0,8})#([a-zA-Z0-9_-]{43}\|[a-zA-Z0-9_-]{22}) |  Body Contents  | No       | Functioning | 
| Sendvid | (https\|http)://sendvid.com/([a-z0-9]{8}) | Status code | No | Functioning

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
