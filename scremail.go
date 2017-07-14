package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	reg = regexp.MustCompile(`\w+@[\w.-]+|\{(?:\w+, *)+\w+\}@[\w.-]+`)

	urlArg       = kingpin.Arg("url", "URL of website.").Required().String()
	usernameFlag = kingpin.Flag("username", "Output just usernames instead of full emails.").Short('u').Bool()
	domainFlag   = kingpin.Flag("domain", "Output just domains instead of full emails.").Short('d').Bool()
	threadsflag  = kingpin.Flag("threads", "Number of threads to use, default: 1.").Short('t').Default("1").Int()
)

func main() {
	kingpin.CommandLine.HelpFlag.Short('h')
	kingpin.Parse()

	body := reqURL(*urlArg)

	outputResults(reg.FindAllString(string(body), -1))
}

func outputResults(results []string) {
	if *domainFlag {
		domains := make([]string, len(results))
		for _, email := range results {
			domains = append(domains, strings.Split(email, "@")[1])
		}
		for _, dom := range uniq(domains) {
			if dom != "" {
				fmt.Printf("%s\n", dom)
			}
		}
		return
	}

	if *usernameFlag {
		for _, email := range results {
			fmt.Printf("%s\n", strings.Split(email, "@")[0])
		}
		return
	}

	for _, email := range results {
		fmt.Printf("%s\n", email)
	}
}

func reqURL(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting to URL\n")
		panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading body of html request\n")
		panic(err)
	}
	return body
}

func uniq(elements []string) (result []string) {
	encountered := map[string]bool{}
	// Create a map of all unique elements.
	for v := range elements {
		encountered[elements[v]] = true
	}
	// Place all keys from the map into a slice.
	result = []string{}
	for key, _ := range encountered {
		result = append(result, key)
	}
	return
}
