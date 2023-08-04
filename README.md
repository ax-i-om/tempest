<p align="center">
  <a><img src="./images/omegac.png" width=280 height="245"></a>
  <p align="center">
    <a href="https://goreportcard.com/report/github.com/ax-i-om/omega"><img src="https://goreportcard.com/badge/github.com/ax-i-om/omega" alt="Go Report Card"></a>
    <a><img src="https://img.shields.io/badge/mk-1-green.svg" alt="Mark 1"></a><br>
    Mass discovery of valid Mega links<br>
</a>
  </p><br>
</p>

## Table of Contents

- [Information](#information)
  - [About](#about)
  - [Disclaimer](#disclaimer)
  - [Todo](#todo)
  - [Installation and Usage](#installation-and-usage)

## Information

### About

Omega is a simple, lightweight, and cross-platform solution designed to enable individuals to efficiently discover and extract active Mega links from paste platforms such as [Rentry.co](https://rentry.co). It was created to address the notable uptick paste sites being used to distribute content that violates copyright and piracy statutes.

### Disclaimer

It is the end user's responsibility to obey all applicable local, state, and federal laws. Developers assume no liability and are not responsible for any misuse or damage caused by this program. By using Omega, you agree to the previous statements.

### Todo

- [ ] Custom modules / custom module submission (submitting new paste sites to scrape from)
- [ ] Custom expressions / custom expression submission (submitting new expressions for other potential mediums of distribution such as Google Drive or Yandex Disk)
- [ ] Better output functionality
- [ ] Proxylist support to bypass rate limiting
- [ ] Concurrent searching
- [ ] General Optimization

### Installation and Usage

1. Fetch the repository via ***git clone***: `git clone https://github.com/ax-i-om/bitcrook.git`
2. Navigate to the root directory of of the cloned repository via ***cd***: `cd omega`
3. In your preferred terminal, enter and run: `go run main.go`

 ***OR***

1. Install the repository via ***go install***: `go install github.com/ax-i-om/omega@latest`
2. In your preferred terminal, enter and run: `omega`