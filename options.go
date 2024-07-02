package main

import (
	"os"
	"errors"
	"github.com/projectdiscovery/goflags"
	"sync"
)

var onceOptions sync.Once
var options = &Options{}

type Options struct {
	URL string
	// Threads int
	// Verbose   bool
	SkipSave   bool
	Calendar     string

	User     string
	Password string

	Proxy string
}

func ParseOptions() (*Options, error) {
	var err error

	onceOptions.Do(func() {

		flagSet := goflags.NewFlagSet()
		flagSet.SetDescription("Example - description TODO")

		// flagSet.CreateGroup("input", "Input",

		// flagSet.IntVarP(&options.Threads, "t", "threads", 10, "threads to run"), //TODO add estimate counter to packets/s
		// flagSet.StringVarP(&options.URL, "u", "url", "", "verbose"),
		// )

		flagSet.CreateGroup("debug", "WebDAV Debug",
			flagSet.StringVarP(&options.URL, "u", "url", "", "target's url"),
			flagSet.StringVarP(&options.User, "l", "login", "", "WebDAV login username"),
			flagSet.StringVarP(&options.Password, "p", "password", "", "WebDAV password (forbid filesystem access in WebDAV and don't forget to clean shell history!)"),
			flagSet.StringVarP(&options.Calendar, "c", "calendar", "", "CalDAV calendar (to-do list) name to use (works only with -u,-l,-p flags)"),
			flagSet.StringVarP(&options.Proxy, "P", "proxy", "", "HTTP proxy to debug errors"),
			// flagSet.BoolVarP(&options.SkipSave, "s", "no-save", false, "skip save to keyring"), //TODO
		)

		// flagSet.CreateGroup("debug", "Debug",
		// 	flagSet.BoolVarP(&options.Verbose, "v", "verbose", false, "verbose output with debugging information"),
		// )

		_ = flagSet.Parse()

		err = options.SanityCheck()

	})

	return options, err
}

func (options *Options) SanityCheck() error {

	if (options.URL != "") || (options.User != "") || (options.Password != "") {
		if (options.URL != "") && (options.User != "") && (options.Password != "") {
		} else {
			return errors.New("-u,-l,-p flags must present")
		}
	}

	if options.Proxy != "" {os.Setenv("HTTP_PROXY", options.Proxy)}
	

	return nil
}
