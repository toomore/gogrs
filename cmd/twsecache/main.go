package main

import (
	"flag"
	"log"

	"github.com/toomore/gogrs/utils"
)

var hCache = utils.NewHTTPCache(utils.GetOSRamdiskPath(), "utf8")
var flushall = flag.Bool("flushall", false, "Clear cache.")

func main() {
	flag.Parse()
	if flag.NFlag() == 0 {
		flag.PrintDefaults()
	}
	if *flushall {
		log.Println("Clear Cache:", hCache.Dir, utils.TempFolderName)
		hCache.FlushAll()
	}
}
