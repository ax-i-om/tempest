package main

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
	"unsafe"
)

const (
	leIndexBits = 6
	leIndexMask = 1<<leIndexBits - 1
	leIndexMax  = 63 / leIndexBits
)

var src = rand.NewSource(time.Now().UnixNano())

var filename = trueRand(16, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

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

	return *(*string)(unsafe.Pointer(&b))
}

func retrieveRaw(url string) {
	method := "GET"
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer res.Body.Close()

	if res.StatusCode == 200 {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		re := regexp.MustCompile("https://mega.nz/(folder|file)/([a-zA-Z0-9]{0,8})#([a-zA-Z0-9_-]{43}|[a-zA-Z0-9_-]{22})")
		x := re.FindAllString(string(body), -1)
		if x == nil {
			// fmt.Println("z")
		} else {
			for _, v := range x {
				x, err := check(v)
				if err != nil {
					fmt.Println(err)
					break
				}
				if x {
					fmt.Println("VALID: ", v)
					write(v)
				}
			}

		}

	}

}

func check(x string) (bool, error) {
	re := regexp.MustCompile("([a-zA-Z0-9]{8}#)")
	pre := re.FindString(x)
	post := strings.Replace(pre, "#", "", -1)

	url := "https://g.api.mega.co.nz/cs?id=5644474&n=" + post
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return false, err
	}
	req.Header.Add("Cookie", "geoip=US")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	if strings.Contains(string(body), "-2") {
		return true, nil
	} else {
		return false, nil
	}

}

func main() {
	fmt.Println("Omega Copyright (C) 2023 Axiom\nThis program comes with ABSOLUTELY NO WARRANTY.\nThis is free software, and you are welcome to redistribute it\nunder certain conditions")

	time.Sleep(5 * time.Second)

	for {
		time.Sleep(250)
		go retrieveRaw("https://rentry.co/" + trueRand(5, "abcdefghijklmnopqrstuvwxyz0123456789") + "/raw")
	}
}

func write(url string) {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	writer := bufio.NewWriter(f)

	defer writer.Flush()

	body, err := io.ReadAll(f)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if !strings.Contains(string(body), url) {
		writer.WriteString(url)
		writer.WriteString("\n")
	}
}
