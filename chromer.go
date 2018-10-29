package main

// Heavily adopted and modified from: https://gist.github.com/nathankerr/38d8b0d45590741b57f5f79be336f07c/revisions

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation
#include "handler.h"
*/
import "C"

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
	"syscall"

	"github.com/andlabs/ui"
)

var labelText chan string
var profiles = []string{"Profile 3", "Profile 1"}

func main() {
	// Defaults to personal profile
	patterns := []string{"atlassian", "atl-paas", "bitbucket", "datadog", "splunk"}

	if fh, err := os.Open(os.Getenv("HOME") + "/" + ".chromer"); err == nil {
		var custom []string
		scanner := bufio.NewScanner(fh)
		for scanner.Scan() {
			custom = append(custom, scanner.Text())
		}
		fh.Close()

		if len(custom) > 0 {
			patterns = custom
		}
	}

	pattern := fmt.Sprintf("(?i)\\b(%s)\\b", strings.Join(patterns, "|"))
	regx, _ := regexp.Compile(pattern)

	labelText = make(chan string, 1)
	C.StartURLHandler()

	wg := sync.WaitGroup{}
	wg.Add(1)
	if err := ui.Main(func() {
		go func() {
			defer wg.Done()
			for url := range labelText {
				ui.QueueMain(func() {
					launchURL(regx, url)
				})
			}
		}()
	}); err != nil {
		log.Fatal(err)
	}

	wg.Wait()
}

//export HandleURL
func HandleURL(u *C.char) {
	labelText <- C.GoString(u)
}

func launchURL(regx *regexp.Regexp, url string) error {
	args := []string{"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
		fmt.Sprintf("--profile-directory=%s", getProfile(regx, url)),
		"-t", url}

	if _, err := syscall.ForkExec(args[0], args, nil); err != nil {
		return err
	}

	return nil
}

func getProfile(reg *regexp.Regexp, url string) string {
	if reg != nil && reg.Match([]byte(url)) {
		return profiles[1]
	}

	return profiles[0]
}
