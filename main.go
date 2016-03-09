package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

var flagUrl *string = flag.String("thread", "", "the chan thread to download")

func exitf(msg string, args ...interface{}) {
	if !strings.HasSuffix(msg, "\n") {
		msg += "\n"
	}
	fmt.Fprintf(os.Stderr, msg, args...)
	os.Exit(1)
}

func main() {
	flag.Parse()

	if *flagUrl == "" {
		exitf("No chan link found")
	}

	// Housten, we have a thread
	src, err := downloadToString(*flagUrl)
	if err != nil {
		exitf("Error while downloading page", err)
	}

	links := filter(src)
	fmt.Printf("Found %d images.\n", len(links))

	x := len(links)
	for i, s := range links {
		fmt.Printf("[%d/%d] %s\n", i, x, s)
		err = downloadToDisk(s) // 4chan only has //i.4cdn... as uri
		if err != nil {
			fmt.Printf("Error while downloading %s: %s", s, err.Error())
		}
	}
}

// download a file to disk, does overwrite existing files
func downloadToDisk(url string) error {
	x := strings.Split(url, "/")
	fn := x[len(x)-1] // get filename

	out, err := os.Create(fn)
	if err != nil {
		return err
	}
	defer out.Close()
	return download(out, url)
}

// download a webpage to memory
func downloadToString(url string) (string, error) {
	buf := bytes.NewBuffer(nil)
	err := download(buf, url)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func download(dst io.Writer, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(dst, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func filter(src string) []string {
	var result []string
	// todo, support more chans
	if strings.Contains(src, "//i.4cdn.org/") {
		result = filter_4chan_s(src)
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
