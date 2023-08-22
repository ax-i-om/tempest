<p align="center">
  <a><img src="./images/icon.png" width=280 height="245"></a>
    <h1 align="center">Vigor</h1>
  <p align="center">
    <a href="https://goreportcard.com/report/github.com/ax-i-om/vigor"><img src="https://goreportcard.com/badge/github.com/ax-i-om/vigor" alt="Go Report Card"></a>
    <a><img src="https://img.shields.io/badge/version-0.1.3-blue.svg" alt="v0.1.3"></a><br>
    Leveraging paste sites as a medium for discovery<br>
</a>
  </p><br>
</p>

## Table of Contents

- [Information](#information)
  - [About](#about)
  - [Disclaimer](#disclaimer)
  - [Installation and Usage](#installation-and-usage)
  - [Cloud Storage / File Sharing Platform Modules](#cloud-storage--file-sharing-platform-modules)
## Information

### About

Vigor is a simple, lightweight, and cross-platform solution designed to enable individuals to efficiently discover and extract active cloud storage/file sharing links from paste platforms such as [Rentry.co](https://rentry.co). It was created to address the notable uptick in paste sites being used to distribute content that violates copyright and piracy statutes.

### Disclaimer

It is the end user's responsibility to obey all applicable local, state, and federal laws. Developers assume no liability and are not responsible for any misuse or damage caused by this program. By using Vigor, you agree to the previous statements.

### Installation and Usage

1. Fetch the repository via ***git clone***: `git clone https://github.com/ax-i-om/vigor.git`
2. Navigate to the root directory of of the cloned repository via ***cd***: `cd vigor`
3. In your preferred terminal, enter and run: `go run main.go`

 ***OR***

1. Install the repository via ***go install***: `go install github.com/ax-i-om/vigor@latest`
2. In your preferred terminal, enter and run: `vigor`

<br>

If you want to output the results to a file, append this to the command: `2>&1 | tee results.txt` <br>
For example `vigor 2>&1 | tee results.txt`   ***or***   `go run main.go 2>&1 | tee results.txt` (may vary depending on operating system) <br>
CAUTION: IF AN ALREADY EXISTING FILE IS SPECIFIED, THIS WILL OVERWRITE THE CONTENTS

### Cloud Storage / File Sharing Platform Modules

| Module    | Expression                                                                                    |   Validation Method   | Domain Variations? | Status |
| :-------: | --------------------------------------------------------------------------------------------- | :------: | -------- | :----: |
| Bunkr      | ^(https\|http)://bunkrr.su/a/([a-zA-Z0-9]{8}) |  Status Code  | Yes       | Functioning | 
| Cyberdrop      | ^(https\|http)://cyberdrop.me/a/([a-zA-Z0-9]{8}) |  Status Code  | No       | Functioning | 
| Dood | ^(https\|http)://doods.pro/((f/[a-z0-9]{10})\|((d/[a-z0-9]{32}\|(d/[a-z0-9]{31})\|(d/[a-z0-9]{12})))\|e/[a-z0-9]{12}) | Body Contents | Yes | Functioning | 
| Gofile      | ^(https\|http)://gofile.io/d/([a-zA-Z0-9]{6}) |  Body Contents  | No       | Functioning | 
| Google Drive | ^(https\|http)://drive.google.com/(folder\|file\|drive)/(d\|folders)/(1[a-zA-Z0-9_-]{32}\|0[a-zA-Z0-9_-]{27}) | Status Code | No | Functioning |
| Mega      | ^(https\|http)://mega.nz/(folder\|file)/([a-zA-Z0-9]{0,8})#([a-zA-Z0-9_-]{43}\|[a-zA-Z0-9_-]{22}) |  Body Contents  | No       | Functioning | 
| Sendvid | ^(https\|http)://sendvid.com/([a-z0-9]{8}) | Status code | No | Functioning