package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/dwbuiten/dumpindex/v7/ffmsindex"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s file.ffindex\n", os.Args[0])
		return
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalf("Can't open %s.", os.Args[1])
	}
	defer f.Close()

	idx, err := ffmsindex.ReadIndex(f)
	if err != nil {
		log.Fatalf(err.Error())
	}

	out, err := json.Marshal(idx)
	if err != nil {
		log.Fatalf(err.Error())
	}

	var b bytes.Buffer
	err = json.Indent(&b, out, "", "    ")
	if err != nil {
		log.Fatalf(err.Error())
	}

	b.WriteTo(os.Stdout)
}
