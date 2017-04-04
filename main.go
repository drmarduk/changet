package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

var errDupe error = errors.New("dupe")

func exitf(msg string, args ...interface{}) {
	if !strings.HasSuffix(msg, "\n") {
		msg += "\n"
	}
	fmt.Fprintf(os.Stderr, msg, args...)
	os.Exit(1)
}

func main() {
	flagURL := flag.String("thread", "", "the chan thread to download")

	flag.Parse()

	if *flagURL == "" {
		exitf("No chan link found")
	}

	var total int64
	var start time.Time = time.Now()
	// Housten, we have a thread
	src, err := downloadToString(*flagURL)
	if err != nil {
		exitf("Error while downloading page", err)
	}

	links := filter(src)
	fmt.Printf("Found %d images.\n", len(links))

	x := len(links)
	for i, s := range links {
		var t float64
		var size int64
		t, size, err = downloadToDisk(s) // 4chan only has //i.4cdn... as uri
		if err == errDupe {
			//fmt.Printf("[%d/%d] %s\n", i+1, x, s+" skipped")
			continue
		}
		if err != nil {
			fmt.Printf("Error while downloading %s: %s", s, err.Error())
		}
		fmt.Printf("[%d/%d] %s\t%.2f KB (%.2f KBytes/s)\n", i+1, x, s, (float64(size) / 1000), t/1000)
		total += size
	}

	end := time.Now()

	fmt.Printf("\nDownloaded %d files. %.2f MB in %s Seconds\n", len(links), float64(total/1000/1000), end.Sub(start).String())
}

func fileexists(file string) bool {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	}
	return true
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

// download a webpage to memory
func downloadToString(url string) (string, error) {
	buf := bytes.NewBuffer(nil)
	_, _, err := download(buf, url)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func download(dst io.Writer, url string) (float64, int64, error) {
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

func filter(src string) []string {
	var result []string
	// todo, support more chans
	if strings.Contains(src, "//i.4cdn.org/") {
		result = filter_4chan_s(src)
	}

	if strings.Contains(src, "krautchan") {
		result = filter_krautchan_s(src)
	}

	if strings.Contains(src, "taychan") {
		result = filter_taychan_b(src)
	}

	if strings.Contains(src, "7chan.org") {
		result = filter_7chan_s(src)
	}

	return removeDuplicates(result)
}

func filter_4chan_s(src string) []string {
	var result []string
	r := regexp.MustCompile("//i.4cdn.org/[a-zA-Z]{1,4}/[0-9]{1,15}.(jpg|jpeg|png|gif|webm)")
	for _, x := range r.FindAllString(src, -1) {
		result = append(result, "http:"+x)
	}
	return result
}

func filter_krautchan_s(src string) []string {
	var result []string
	r := regexp.MustCompile("/files/[0-9]{7,15}.(jpg|jpeg|png|gif|webm|gifv)")
	for _, x := range r.FindAllString(src, -1) {
		result = append(result, "http://krautchan.net"+x)
	}
	return result
}

func filter_taychan_b(src string) []string {
	var result []string
	// r := regexp.MustCompile("<p class="fileinfo">File: <a href="/b/src/1470058508938-0.webm">1470058508938-0.webm</a>")
	r := regexp.MustCompile("/b/src/[0-9]{7,15}(|-[0-9]{1,3}).(jpg|jpeg|png|gif|webm|gifv)")
	for _, x := range r.FindAllString(src, -1) {
		result = append(result, "http://taychan.eu"+x)
	}
	return result
}

func filter_7chan_s(src string) []string {
	var result []string
	r := regexp.MustCompile("")
	for _, x := range r.FindAllString(src, -1) {
		result = append(result, "http://7chan.org/"+x)
	}
	return result
}

func removeDuplicates(slice []string) []string {
	var result []string
	tmp := map[string]bool{}

	for v := range slice {
		tmp[slice[v]] = true
	}

	for k, _ := range tmp {
		result = append(result, k)
	}
	return result
}
