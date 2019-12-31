// Command pkg-diff-example implements a subset of the diff command using
// github.com/pkg/diff.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/pkg/diff"
	"github.com/pkg/diff/write"
)

var color = flag.Bool("color", false, "colorize the output")

// check logs a fatal error and exits if err is not nil.
func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "pkg-diff-example [flags] file1 file2\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	log.SetPrefix("pkg-diff-example: ")
	log.SetFlags(0)

	flag.Usage = usage
	flag.Parse()

	var aName, bName string
	var a, b interface{}
	switch flag.NArg() {
	case 2:
		// A human has invoked us.
		aName, bName = flag.Arg(0), flag.Arg(1)
	case 7, 9:
		// We are a git external diff tool.
		// We have been passed the following arguments:
		//   path old-file old-hex old-mode new-file new-hex new-mode [new-path similarity-metrics]
		aName = "a/" + flag.Arg(0)
		if flag.NArg() == 7 {
			bName = "b/" + flag.Arg(0)
		} else {
			bName = "b/" + flag.Arg(7)
		}
		var err error
		a, err = ioutil.ReadFile(flag.Arg(1))
		check(err)
		b, err = ioutil.ReadFile(flag.Arg(4))
		check(err)
	default:
		flag.Usage()
	}

	var opts []write.Option
	if *color {
		opts = append(opts, write.TerminalColor())
	}

	err := diff.Text(aName, bName, a, b, os.Stdout, opts...)
	check(err)
}
