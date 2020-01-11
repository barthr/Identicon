package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/barthr/identicon"
)

var welcomeSignature = `
Usage of Identicon made By Bart
_______________________________
	< Identicon >
-------------------------------

-name string:
	Set the name where you want to generate a identicon for

`

func main() {
	var (
		name = flag.String("name", "", "Set the name where you want to generate a identicon for")
	)
	flag.Parse()

	if *name == "" {
		flag.Usage = func() {
			fmt.Println(welcomeSignature)
		}
		flag.Usage()
		os.Exit(0)
	}

	generatedIdenticon := identicon.Generate([]byte(*name))

	f, err := os.Create(generatedIdenticon.Name + ".png")
	if err != nil {
		fmt.Printf("error: failed creating file for output png")
		return
	}
	defer f.Close()

	if err := generatedIdenticon.WriteImage(f); err != nil {
		fmt.Printf("error failed writing image to file")
		return
	}
}
