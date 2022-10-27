package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	in := os.Stdin
	out := os.Stdout
	strictMode := flag.Bool("strict", false, "use strict mode (like error on undefined columns)")
	flag.Parse()

	if len(flag.Args()) == 0 {
		fmt.Println("Pass template in first argument, for quick start you can use following templates for this file:")
		PrintExample(in)
		os.Exit(1)
	}

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
