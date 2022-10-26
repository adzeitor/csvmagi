package main

import (
	"flag"
	"log"
	"os"
)

func main() {
	in := os.Stdin
	out := os.Stdout
	strictMode := flag.Bool("strict", false, "use strict mode (like error on undefined columns)")
	flag.Parse()

	magi, err := New(flag.Arg(0), Config{
		StrictMode: *strictMode,
	})
	if err != nil {
		log.Fatalln(err)
	}

	err = magi.ReadAndExecute(in, out)
	if err != nil {
		log.Fatalln(err)
	}
}
