// Copyright © 2017 Toomore Chiang
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"log"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/toomore/gogrs/utils"
)

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

// cacheCmd represents the cache command
var cacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "cache system",
	Long:  `快取機制`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var cacheFlushAllCmd = &cobra.Command{
	Use:   "flushall",
	Short: "flushall cache",
	Long:  `清除所有快取`,
	Run: func(cmd *cobra.Command, args []string) {
		var hCache = utils.NewHTTPCache(utils.GetOSRamdiskPath(""), "utf8")
		outputNote("Clear Cache:", hCache.Dir, utils.TempFolderName)
		hCache.FlushAll()
		outputDone("Done!")
	},
}

func init() {
	RootCmd.AddCommand(cacheCmd)
	cacheCmd.AddCommand(cacheFlushAllCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cacheCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cacheCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
