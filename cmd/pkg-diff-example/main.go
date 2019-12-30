// Command pkg-diff-example implements a subset of the diff command using
// github.com/pkg/diff.
package main

import (
	"flag"
	"fmt"
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
	if len(flag.Args()) != 2 {
		flag.Usage()
	}

	aName := flag.Arg(0)
	bName := flag.Arg(1)

	var opts []write.Option
	if *color {
		opts = append(opts, write.TerminalColor())
	}

	err := diff.Text(aName, bName, nil, nil, os.Stdout, opts...)
	check(err)
}
