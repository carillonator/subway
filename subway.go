package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var cis = ComplexInfoSet{}
var mtaKey = ""

const usage = `
  Usage: subway [ server | complexId,... ]
`

func init() {
	var err error
	cis, err = newComplexInfoSet()
	if err != nil {
		log.Fatal(err)
	}

	mtaKey = os.Getenv("MTA_API_KEY")
	if mtaKey == "" {
		log.Fatal("MTA_API_KEY not set")
	}
}

func main() {

	var arg string
	if len(os.Args) > 1 {
		arg = os.Args[1]
	} else {
		fmt.Println(usage)
		os.Exit(2)
	}

	var cids []int
	if arg == "server" {
		startServer()

	} else {
		for _, s := range strings.Split(arg, ",") {
			id, err := strconv.Atoi(s)
			if err != nil {
				fmt.Println(usage)
				os.Exit(2)
			}
			cids = append(cids, id)
		}
	}

	ss, err := NewStationSet(cids, cis)
	if err != nil {
		log.Fatal(err)
	}

	printText(ss, os.Stdout)
}
