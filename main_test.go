package main

import "testing"

func TestEnvrive(t *testing.T) {
	in := `<img src="https://enrive.org/media/sys/download.png" title="Download File" alt="Download File"</a> <a " class="hyperlinkMediaFileName" href="https://enrive.org/media/images/masterchanisretarded_530.jpg" target="_blank"><span class="mediaFileName">1.jpg</span></a> 24.54 KB</span><span class="username">Anonymous</span><span class="dateAndTime">`

	l := filter_enrive(in)
	found := false
	for _, u := range l {
		if u == "https://enrive.org/media/images/masterchanisretarded_530.jpg" {
			found = true
		}
	}

	if !found {
		t.Fatalf("no match found")
	}
}
