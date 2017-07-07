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
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/toomore/gogrs/tradingdays"
	"github.com/toomore/gogrs/twse"
	"github.com/toomore/gogrs/utils"
)

// exampleCmd represents the example command
var exampleCmd = &cobra.Command{
	Use:   "example",
	Short: "Show example",
	Long:  `Show gogrs example.`,
	Run: func(cmd *cobra.Command, args []string) {
		var stock = twse.NewTWSE("2618", time.Date(2015, 3, 20, 0, 0, 0, 0, time.Local))
		stock.Get()
		showAll(stock)
		fmt.Println("-----------------------------")
		stock.PlusData()
		showAll(stock)
		fmt.Println("-----------------------------")
		fmt.Println(tradingdays.IsOpen(2015, 5, 1))
	},
}

func showAll(stock *twse.Data) {
	fmt.Println(stock.RawData)
	fmt.Println(stock.MA(6))
	fmt.Println(stock.MAV(6))
	fmt.Println(stock.GetPriceList())
	fmt.Println(utils.ThanPastFloat64(stock.GetPriceList(), 3, true))
	fmt.Println(utils.ThanPastFloat64(stock.GetPriceList(), 3, false))
	fmt.Println(stock.GetVolumeList())
	fmt.Println(utils.ThanPastUint64(stock.GetVolumeList(), 3, true))
	fmt.Println(utils.ThanPastUint64(stock.GetVolumeList(), 3, false))
	fmt.Println(stock.GetRangeList())
	fmt.Println(stock.IsRed())
}

func init() {
	RootCmd.AddCommand(exampleCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// exampleCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// exampleCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
