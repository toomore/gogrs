package main

import (
	"flag"
	"log"

	"github.com/fatih/color"
	"github.com/toomore/gogrs/utils"
)

var hCache = utils.NewHTTPCache(utils.GetOSRamdiskPath(), "utf8")
var flushall = flag.Bool("flushall", false, "Clear cache.")

func outputNote(note ...interface{}) {
	color.Set(color.FgBlue, color.Bold)
	log.Println(note...)
	color.Unset()
}

func outputDone(note ...interface{}) {
	color.Set(color.FgGreen, color.Bold)
	log.Println(note...)
	color.Unset()
}

func main() {
	flag.Parse()
	if flag.NFlag() == 0 {
		flag.PrintDefaults()
	}
	if *flushall {
		outputNote("Clear Cache:", hCache.Dir, utils.TempFolderName)
		hCache.FlushAll()
		outputDone("Done!")
	}
}
