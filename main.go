package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

var thread *string = flag.String("thread", "", "the chan thread to download")

func main() {
	flag.Parse()

	if *thread == "" {
		fmt.Println("Please provide a chan link.")
		os.Exit(0)
	}
	// Housten, we have a thread
	src, err := dl(*thread)
	if err != nil {
		fmt.Printf("[!] Error while downloading page: %s\n", err.Error())
		os.Exit(1)
	}
	fmt.Printf("Got %d bytes of html\n", len(src))

	links := filter(src)
	fmt.Printf("Got %d images\n", len(links))

	for _, s := range links {
		fmt.Println("Download: " + s)

		dl2file("http:"+s, "")
	}
}

func dl2file(url, file string) error {

	x := strings.Split(url, "/")
	fn := x[len(x)-1]

	out, err := os.Create(fn)
	if err != nil {
		fmt.Println("[!] Error while creating file " + file + ": " + err.Error())
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("[!] Error while HTTP GET for %s: %s\n", url, err.Error())

		return err
	}
	defer resp.Body.Close()
	if err != nil {
		fmt.Printf("[!] Error while downloading file %s: %s\n", url, err.Error())
		return err
	}
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Println("[!] Error while transfering file to disk: " + err.Error())
		return err
	}
	return nil
}

func dl(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func filter(src string) []string {
	// todo, support more chans
	l := filter_4chan_s(src)

	// filter out dups
	tmp := map[string]bool{}
	for v := range l {
		tmp[l[v]] = true
	}

	result := []string{}
	for k, _ := range tmp {
		result = append(result, k)
	}
	return result
}

func filter_4chan_s(src string) []string {
	r := regexp.MustCompile("//i.4cdn.org/[a-zA-Z]{1,4}/[0-9]{1,15}.(jpg|jpeg|png|gif|webm)")
	return r.FindAllString(src, -1)
}
