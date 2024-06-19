package main

import (
	"errors"
	"github.com/projectdiscovery/goflags"
	"sync"
)

var onceOptions sync.Once
var options = &Options{}

type Options struct {
	URL     string
	Threads int
	// Verbose   bool

	User     string
	Password string
}

func ParseOptions() (*Options, error) {
	var err error

	onceOptions.Do(func() {

		flagSet := goflags.NewFlagSet()
		flagSet.SetDescription("Example - description TODO")

		flagSet.CreateGroup("input", "Input",
			flagSet.StringVarP(&options.URL, "u", "url", "", "target's url"),
			flagSet.IntVarP(&options.Threads, "t", "threads", 10, "threads to run"), //TODO add estimate counter to packets/s
			// flagSet.StringVarP(&options.URL, "u", "url", "", "verbose"),
		)

		flagSet.CreateGroup("debug", "WebDAV Debug",
			flagSet.StringVarP(&options.User, "l", "login", "", "WebDAV login"),
			flagSet.StringVarP(&options.Password, "p", "password", "", "WebDAV password (forbid filesystem access!!!)"),
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

	if options.URL == "" {
		return errors.New("-u flag must present")
	}

	return nil
}
