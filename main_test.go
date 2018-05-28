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

func TestTerchan(t *testing.T) {
	out := "https://terchan.xyz/src/1526856270818.jpg"
	in := `<a href="13602.html#13606">No.</a><a href="13602.html#q13606" onclick="javascript:quotePost('13606')">13606</a>
	</span><br><span class="filesize"><a href="../src/1526856270818.jpg" download="nfv.jpg">File</a>: (247.49KB, 807x605, nfv.jpg)</span><br><div id="thumbfile13606"><a href="../src/1526856270818.jpg" target="_blank" onclick="return expandFile(event, '13606');">
		<img src="../thumb/1526856270818s.jpg" alt="13606" class="thumb" id="thumbnail13606" width="250" height="187">
		</a></div><div id="expand13606" style="display: none;">%3Ca%20href%3D%22..%2Fsrc%2F1526856270818.jpg%22%20onclick%3D%22return%20expandFile%28event%2C%20%2713606%27%29%3B%22%3E%3Cimg%20src%3D%22..%2Fsrc%2F1526856270818.jpg%22%20width%3D%22807%22%20style%3D%22max-width%3A%20100%25%3Bheight%3A%20auto%3B%22%3E%3C%2Fa%3E</div>
		<div id="file13606" class="thumb" style="display: none;"></div><div class="message">`
	l := filter_terchan(in)
	found := false
	for _, u := range l {
		if u == out {
			found = true
		}
	}

	if !found {
		t.Fatalf("found: %vwant: %s\n", l, out)
		t.Fatalf("no match found")
	}
}
