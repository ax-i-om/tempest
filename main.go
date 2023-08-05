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

	return string(b[:])
}

func write(url string) error {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(f)

	body, err := io.ReadAll(f)
	if err != nil {
		writer.Flush()
		return err
	}

	if !strings.Contains(string(body), url) {
		writer.WriteString(url)
		writer.WriteString("\n")
	}
	writer.Flush()
	return nil
}

func getRes(url string) (*http.Response, error) {
	method := "GET"
	client := &http.Client{}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func checkRes(res *http.Response) error {
	if res.StatusCode == 200 {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		re := regexp.MustCompile("https://mega.nz/(folder|file)/([a-zA-Z0-9]{0,8})#([a-zA-Z0-9_-]{43}|[a-zA-Z0-9_-]{22})")
		x := re.FindAllString(string(body), -1)
		if len(x) > 0 {
			for _, v := range x {
				x, err := check(v)
				if err != nil {
					return err
				}
				if x {
					fmt.Println("VALID: ", v)
					err := write(v)
					if err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func check(x string) (bool, error) {
	re := regexp.MustCompile("([a-zA-Z0-9]{8}#)")
	pre := re.FindString(x)
	post := strings.Replace(pre, "#", "", -1)

	url := "https://g.api.mega.co.nz/cs?id=5644474&n=" + post

	res, err := getRes(url)
	if err != nil {
		return false, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return false, err
	}
	if strings.Contains(string(body), "-2") {
		return true, res.Body.Close()
	} else {
		return false, res.Body.Close()
	}

}

func runner(renturl string) error {
	res, err := getRes(renturl)
	if err != nil {
		return err
	}

	err = checkRes(res)
	if err != nil {
		return err
	}

	return res.Body.Close()
}

func main() {
	fmt.Println("Omega Copyright (C) 2023 Axiom\nThis program comes with ABSOLUTELY NO WARRANTY.\nThis is free software, and you are welcome to redistribute it\nunder certain conditions")

	time.Sleep(5 * time.Second)

	for {
		time.Sleep(250)
		go func() {
			err := runner("https://rentry.co/" + trueRand(5, "abcdefghijklmnopqrstuvwxyz0123456789") + "/raw")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}()
	}
}
