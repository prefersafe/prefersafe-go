package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

// emptyGIF is the tiniest possible 1-pixel GIF image (handtinyblack.gif) from
// http://probablyprogramming.com/2009/03/15/the-tiniest-gif-ever
var emptyGIF = []byte("GIF89a\x01\x00\x01\x00\x00\xff\x00,\x00\x00\x00\x00\x01\x00\x01\x00\x00\x02\x00;")

// headerContains looks for the value in all headers with the given name name,
// and returns true if the value was found.
//
// Examples when it returns true for "Prefer", "safe" lookup:
//
// - a single header value:
//
//  Prefer: safe
//
// - a single header containing space-separated values:
//
//  Prefer: safe crazy
//
// - multiple headers:
//
//  Prefer: safe
//  Prefer: crazy
//
func headerContains(h http.Header, name, value string) bool {
	for _, s := range h[name] {
		for _, it := range strings.Fields(s) {
			if it == value {
				return true
			}
		}
	}
	return false
}

// hasPreferSafe returns true if the given header contains Prefer:safe.
func hasPreferSafe(h http.Header) bool {
	return headerContains(h, "Prefer", "safe")
}

// imageHandler serves a valid GIF image if Prefer:safe header is set,
// otherwise returns a zero-length response, causing image decoding
// error in browsers.
func imageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/gif")
	if hasPreferSafe(r.Header) {
		w.Write(emptyGIF)
	}
}

// jsonpHandler serves a script that calls PreferSafe() global function with
// boolean value indicating whether the Prefer:safe header was set.
func jsonpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript")
	fmt.Fprintf(w, "PreferSafe(%v);\n", hasPreferSafe(r.Header))
}

// serveMux returns a new http.Handler which serves content.
func serveMux() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/image.gif", imageHandler)
	mux.HandleFunc("/jsonp.js", jsonpHandler)
	return mux
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.Handle("/", serveMux())
	log.Printf("Serving from port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
