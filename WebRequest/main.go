// Go contains rich function for grab web contents. _net/http_ is the major
// library.
// Ref: [golang.org](http://golang.org/pkg/net/http/#pkg-examples).
package main

import (
	"net/http"
)
import "io/ioutil"
import "fmt"
import "strings"

// keep first n lines
func keepLines(s string, n int) string {
	result := strings.Join(strings.Split(s, "\n")[:n], "\n")
	return strings.Replace(result, "\r", "", -1)
}

func main() {
	// We can use GET form to get result.
	resp, err := http.Get("http://www.somesweetsite.com")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println("get:\n", keepLines(string(body), 3))
}
