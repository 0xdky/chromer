package main

// Heavily adopted and modified from: https://gist.github.com/nathankerr/38d8b0d45590741b57f5f79be336f07c/revisions
// Get Chrome profile names from: chrome://version/

/*
#cgo CFLAGS: -I/usr/local/include/macosx -x objective-c
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

	"github.com/fsnotify/fsnotify"

	"github.com/andlabs/ui"
)

var labelText chan string
var updateConfig chan bool

type configBlock struct {
	profile string
	regex   *regexp.Regexp
}

func main() {
	cfg := os.Getenv("HOME") + "/.chromer"
	logfh, err := os.Create(os.Getenv("HOME") + "/.chromer.log")
	if err != nil {
		log.Fatal(err)
	}
	defer logfh.Close()
	logger := log.New(logfh, "chromer: ", log.LstdFlags)

	// Load the mandatory config data
	configs, err := loadConfig(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Prepare to receive the clicked URL
	labelText = make(chan string, 1)
	updateConfig = make(chan bool, 1)

	monitorConfig(cfg, updateConfig, logger)
	C.StartURLHandler()

	wg := sync.WaitGroup{}
	wg.Add(1)
	if err := ui.Main(func() {
		go func() {
			defer wg.Done()

			for {
				select {
				case url := <-labelText:
					ui.QueueMain(func() {
						launchURL(configs, url)
					})
				case <-updateConfig:
					configs, _ = loadConfig(cfg)
				}
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

func monitorConfig(cfg string, ch chan bool, logger *log.Logger) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					logger.Fatal("bailing out due to fsnotify event listening error")
				}

				logger.Println("event:", event)
				if event.Op&(fsnotify.Write|fsnotify.Rename) != 0 {
					logger.Println("modified file:", event.Name)
					ch <- true
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					logger.Fatal("bailing out due to fsnotify error listening error")
				}
				logger.Println("error:", err)
			}
		}
	}()

	return watcher.Add(cfg)
}

func loadConfig(cfg string) ([]configBlock, error) {
	var err error
	var fh *os.File

	if fh, err = os.Open(cfg); err != nil {
		return nil, err
	}

	var profile string
	var patterns []string
	var config []configBlock

	scanner := bufio.NewScanner(fh)
	for scanner.Scan() {
		// Sanitized line from config
		line := strings.TrimSpace(scanner.Text())

		// Extract the profile name block
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			if len(profile) > 0 {
				if len(patterns) > 0 {
					config = append(config, configBlock{profile,
						regexp.MustCompile(fmt.Sprintf("(?i)\\b(%s)\\b", strings.Join(patterns, "|")))})
					patterns = nil
				} else {
					config = append(config, configBlock{profile, nil})
				}
			}

			profile = strings.Trim(line, "[]")
		} else if len(line) > 0 && !strings.HasPrefix(line, "#") {
			patterns = append(patterns, line)
		}
	}
	fh.Close()

	// Catch the last config block
	if len(profile) > 0 && len(patterns) > 0 {
		config = append(config, configBlock{profile,
			regexp.MustCompile(fmt.Sprintf("(?i)\\b(%s)\\b", strings.Join(patterns, "|")))})
		patterns = nil
	}

	return config, nil
}

func launchURL(configs []configBlock, url string) error {
	args := []string{"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
		fmt.Sprintf("--profile-directory=%s", getProfile(configs, url)),
		"-t", url}

	if _, err := syscall.ForkExec(args[0], args, nil); err != nil {
		return err
	}

	return nil
}

func getProfile(configs []configBlock, url string) string {
	urlBytes := []byte(url)
	var profile string
	for _, cfg := range configs {
		if len(profile) == 0 && cfg.regex == nil {
			profile = cfg.profile
		} else if cfg.regex != nil && cfg.regex.Match(urlBytes) {
			profile = cfg.profile
			break
		}
	}

	return profile
}
