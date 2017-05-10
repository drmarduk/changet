package main

import (
	"bytes"
	"io"
	"net/http"
	derp "net/url"
	"os"
	"strings"
	"time"

	"golang.org/x/net/proxy"
)

func downloadOnion(dst io.Writer, url string) (float64, int64, error) {
	var start, end time.Time
	start = time.Now()

	p, err := derp.Parse("socks5://127.0.0.1:9050")
	if err != nil {
		return 0.0, 0, err
	}

	dialer, err := proxy.FromURL(p, proxy.Direct)
	if err != nil {
		return 0.0, 0, err
	}

	t := &http.Transport{Dial: dialer.Dial}
	client := &http.Client{Transport: t}

	resp, err := client.Get(url)
	if err != nil {
		return 0.0, 0, err
	}
	defer resp.Body.Close()

	var size int64
	size, err = io.Copy(dst, resp.Body)
	if err != nil {
		return 0.0, 0, err
	}

	end = time.Now()
	bps := float64(size) / end.Sub(start).Seconds()
	return bps, size, nil
}

// download saves an url to io.Writer and returns the bits/s, size and error
func download(dst io.Writer, url string) (float64, int64, error) {
	if strings.Contains(url, ".onion") || strings.Contains(url, "ccluster") {
		return downloadOnion(dst, url)
	}
	var start, end time.Time
	start = time.Now()
	resp, err := http.Get(url)
	if err != nil {
		return 0.0, 0, err
	}
	defer resp.Body.Close()

	var size int64
	size, err = io.Copy(dst, resp.Body)
	if err != nil {
		return 0.0, 0, err
	}

	end = time.Now()
	bps := float64(size) / end.Sub(start).Seconds()
	return bps, size, nil
}

// download a webpage to memory
func downloadToString(url string) (string, error) {
	buf := bytes.NewBuffer(nil)
	_, _, err := download(buf, url)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

// download a file to disk, does overwrite existing files
func downloadToDisk(url string) (float64, int64, error) {
	x := strings.Split(url, "/")
	fn := x[len(x)-1] // get filename

	if fileexists(fn) {
		return 0.0, 0, errDupe
	}

	out, err := os.Create(fn)
	if err != nil {
		return 0.0, 0, err
	}
	defer out.Close()
	return download(out, url)
}

func fileexists(file string) bool {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	}
	return true
}
