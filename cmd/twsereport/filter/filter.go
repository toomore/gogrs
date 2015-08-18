// Package filter - all check list.
package filter

import (
	"github.com/toomore/gogrs/twse"
	"github.com/toomore/gogrs/utils"
)

// CheckGroup for base filter interface.
type CheckGroup interface {
	No() uint64
	String() string
	CheckFunc(...*twse.Data) bool
	Mindata() int
}

type checkGroupList []CheckGroup

func (c *checkGroupList) Add(f CheckGroup) {
	if (*c)[0] == nil {
		(*c)[0] = f
	} else {
		*c = append(*c, f)
	}
}

// Check01 MA 3 > 6 > 18
type Check01 struct{}

// No for check no.
func (Check01) No() uint64 {
	return 1
}

// String to string.
func (Check01) String() string {
	return "MA 3 > 6 > 18"
}

// Mindata is filter required a minimum of data.
func (Check01) Mindata() int {
	return 18
}

// CheckFunc func to check.
func (Check01) CheckFunc(b ...*twse.Data) bool {
	if !prepareData(b...)[0] {
		return false
	}
	var ma3 = b[0].MA(3)
	if days, ok := utils.CountCountineFloat64(utils.DeltaFloat64(ma3)); !ok || days == 0 {
		return false
	}
	var ma6 = b[0].MA(6)
	if days, ok := utils.CountCountineFloat64(utils.DeltaFloat64(ma6)); !ok || days == 0 {
		return false
	}
	var ma18 = b[0].MA(18)
	if days, ok := utils.CountCountineFloat64(utils.DeltaFloat64(ma18)); !ok || days == 0 {
		return false
	}
	//log.Println(ma3[len(ma3)-1], ma6[len(ma6)-1], ma18[len(ma18)-1])
	if ma3[len(ma3)-1] > ma6[len(ma6)-1] && ma6[len(ma6)-1] > ma18[len(ma18)-1] {
		return true
	}
	return false
}

// Check02 量大於前三天 K 線收紅
type Check02 struct{}

// No for check no.
func (Check02) No() uint64 {
	return 2
}

// String to string.
func (Check02) String() string {
	return "量大於前三天 K 線收紅"
}

// Mindata is filter required a minimum of data.
func (Check02) Mindata() int {
	return 4
}

// CheckFunc func to check.
func (Check02) CheckFunc(b ...*twse.Data) bool {
	if !prepareData(b...)[0] {
		return false
	}
	return utils.ThanSumPastUint64((*b[0]).GetVolumeList(), 3, true) && ((*b[0]).IsRed() || (*b[0]).IsThanYesterday())
}

// Check03 量或價走平 45 天
type Check03 struct{}

// No for check no.
func (Check03) No() uint64 {
	return 3
}

// String to string.
func (Check03) String() string {
	return "量或價走平 45 天"
}

// Mindata is filter required a minimum of data.
func (Check03) Mindata() int {
	return 45
}

// CheckFunc func to check.
func (Check03) CheckFunc(b ...*twse.Data) bool {
	if !prepareData(b...)[0] {
		return false
	}

	var (
		price  = b[0].GetPriceList()
		volume = b[0].GetVolumeList()
	)

	return price[len(price)-1] > 10 &&
		(utils.SD(price[len(price)-45:]) < 0.25 ||
			utils.SDUint64(volume[len(volume)-45:]) < 0.25)
}

// Check04 (MA3 < MA6) > MA18 and MA3UP(1)
type Check04 struct{}

// No for check no.
func (Check04) No() uint64 {
	return 4
}

// String to string.
func (Check04) String() string {
	return "(MA3 < MA6) > MA18 and MA3UP(1)"
}

// Mindata is filter required a minimum of data.
func (Check04) Mindata() int {
	return 18
}

// CheckFunc func to check.
func (Check04) CheckFunc(b ...*twse.Data) bool {
	if !prepareData(b...)[0] {
		return false
	}
	var ma3 = b[0].MA(3)
	if days, up := utils.CountCountineFloat64(utils.DeltaFloat64(ma3)); up && days == 1 {
		var (
			ma6      = b[0].MA(6)
			ma18     = b[0].MA(18)
			ma3Last  = len(ma3) - 1
			ma6Last  = len(ma6) - 1
			ma18Last = len(ma18) - 1
		)
		return (ma3[ma3Last] > ma18[ma18Last] && ma6[ma6Last] > ma18[ma18Last]) && ma3[ma3Last] < ma6[ma6Last]
	}
	return false
}

// Check05 三日內最大量 K 線收紅 收在 MA18 之上
type Check05 struct{}

// No for check no.
func (Check05) No() uint64 {
	return 5
}

// String to string.
func (Check05) String() string {
	return "三日內最大量 K 線收紅 收在 MA18 之上"
}

// Mindata is filter required a minimum of data.
func (Check05) Mindata() int {
	return 18
}

// CheckFunc func to check.
func (Check05) CheckFunc(b ...*twse.Data) bool {
	if !prepareData(b...)[0] {
		return false
	}
	var (
		vols        = b[0].GetVolumeList()
		volsFloat64 = make([]float64, 3)
	)
	for i, v := range vols[len(vols)-3:] {
		volsFloat64[i] = float64(v)
	}
	if days, up := utils.CountCountineFloat64(utils.DeltaFloat64(volsFloat64)); up && days >= 1 && b[0].IsRed() {
		var (
			ma18      = b[0].MA(18)
			priceList = b[0].GetPriceList()
		)

		if priceList[len(priceList)-1] > ma18[len(ma18)-1] {
			return true
		}
	}
	return false
}

// Check06 漲幅 7% 以上
type Check06 struct{}

// No for check no.
func (Check06) No() uint64 {
	return 6
}

// String to string.
func (Check06) String() string {
	return "漲幅 7% 以上"
}

// Mindata is filter required a minimum of data.
func (Check06) Mindata() int {
	return 1
}

// CheckFunc func to check.
func (Check06) CheckFunc(b ...*twse.Data) bool {
	if !prepareData(b...)[0] {
		return false
	}

	priceList := b[0].GetPriceList()
	openList := b[0].GetOpenList()
	price := priceList[len(priceList)-1]
	open := openList[len(openList)-1]

	if price > open && (price-open)/open > 0.068 {
		return true
	}

	return false
}

// Check07 多方力道 > 0.75
type Check07 struct{}

// No for check no.
func (Check07) No() uint64 {
	return 7
}

// String to string.
func (Check07) String() string {
	return "多方力道 > 0.75"
}

// Mindata is filter required a minimum of data.
func (Check07) Mindata() int {
	return 1
}

// CheckFunc func to check.
func (Check07) CheckFunc(b ...*twse.Data) bool {
	if !prepareData(b...)[0] {
		return false
	}

	var power []float64
	power = utils.CalLHPower(b[0].GetPriceList(), b[0].GetLowList(), b[0].GetHighList())

	if power[len(power)-1] > 0.75 {
		return true
	}

	return false
}

func prepareData(b ...*twse.Data) []bool {
	var (
		result  []bool
		mindata int
	)

	for i := range AllList {
		if AllList[i].Mindata() > mindata {
			mindata = AllList[i].Mindata()
		}
	}

	for i := range b {
		result = make([]bool, len(b))
		b[i].Get()
		if b[i].Len() < mindata {
			start := b[i].Len()
			for {
				b[i].PlusData()
				if b[i].Len() > mindata {
					result[i] = true
					break
				}
				if b[i].Len() == start {
					break
				}
				start = b[i].Len()
			}
			if b[i].Len() < mindata {
				result[i] = false
			}
		} else {
			result[i] = true
		}
	}
	return result
}

// AllList is all check group.
var AllList = make(checkGroupList, 1)

func init() {
	AllList.Add(CheckGroup(Check01{}))
	AllList.Add(CheckGroup(Check02{}))
	AllList.Add(CheckGroup(Check03{}))
	AllList.Add(CheckGroup(Check04{}))
	AllList.Add(CheckGroup(Check05{}))
	AllList.Add(CheckGroup(Check06{}))
	AllList.Add(CheckGroup(Check07{}))
}
