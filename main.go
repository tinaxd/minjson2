package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/tinaxd/minjson2/min"
)

func main() {
	fs := flag.NewFlagSet("minjson2", flag.ExitOnError)
	mode := fs.String("mode", "", "minjson mode")
	infile := fs.String("in", "", "input filename")
	outfile := fs.String("out", "", "output filename")

	fs.Parse(os.Args[1:])

	switch *mode {
	case "minify":
		if *infile == "" || *outfile == "" {
			fmt.Fprintln(os.Stderr, "input or output file is not set.")
			os.Exit(1)
		}

		ifile, err := os.Open(*infile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		defer ifile.Close()

		ofile, err := os.Create(*outfile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		defer ofile.Close()

		ireader := bufio.NewReader(ifile)
		owriter := bufio.NewWriter(ofile)
		min.MinifyJSON(ireader, owriter)

	case "pretty":
		if *infile == "" || *outfile == "" {
			fmt.Fprintln(os.Stderr, "input or output file is not set.")
			os.Exit(1)
		}

		ifile, err := os.Open(*infile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		defer ifile.Close()

		ofile, err := os.Create(*outfile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		defer ofile.Close()

		setting := min.PrettySetting{IndentWidth: 2}

		ireader := bufio.NewReader(ifile)
		owriter := bufio.NewWriter(ofile)
		min.PrettyJSON(ireader, owriter, setting)

	default:
		if *mode == "" {
			fmt.Fprintln(os.Stderr, "mode is not set.")
			fmt.Fprintln(os.Stderr, "available modes: minify")
			os.Exit(1)
		} else {
			fmt.Fprintf(os.Stderr, "unknown mode: %s\n", *mode)
			fmt.Fprintln(os.Stderr, "available modes: minify")
			os.Exit(1)
		}
	}
}
