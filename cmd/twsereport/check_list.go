package main

import (
	"github.com/toomore/gogrs/twse"
	"github.com/toomore/gogrs/utils"
)

type checkGroup interface {
	String() string
	CheckFunc(...*twse.Data) bool
	Mindata() int
}

type check01 struct{}

func (check01) String() string {
	return "MA 3 > 6 > 18"
}

func (check01) Mindata() int {
	return 18
}

func (check01) CheckFunc(b ...*twse.Data) bool {
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

type check02 struct{}

func (check02) String() string {
	return "量大於前三天 K 線收紅"
}

func (check02) Mindata() int {
	return 4
}

func (check02) CheckFunc(b ...*twse.Data) bool {
	if !prepareData(b...)[0] {
		return false
	}
	return utils.ThanSumPastUint64((*b[0]).GetVolumeList(), 3, true) && ((*b[0]).IsRed() || (*b[0]).IsThanYesterday())
}

type check03 struct{}

func (check03) String() string {
	return "量或價走平 45 天"
}

func (check03) Mindata() int {
	return 45
}

func (check03) CheckFunc(b ...*twse.Data) bool {
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

type check04 struct{}

func (check04) String() string {
	return "(MA3 < MA6) > MA18 and MA3UP(1)"
}

func (check04) Mindata() int {
	return 18
}

func (check04) CheckFunc(b ...*twse.Data) bool {
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

type check05 struct{}

func (check05) String() string {
	return "三日內最大量 K 線收紅 收在 MA18 之上"
}

func (check05) Mindata() int {
	return 18
}

func (check05) CheckFunc(b ...*twse.Data) bool {
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

func prepareData(b ...*twse.Data) []bool {
	var (
		result  []bool
		mindata int
	)

	for i := range ckList {
		if ckList[i].Mindata() > mindata {
			mindata = ckList[i].Mindata()
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

func init() {
	ckList.Add(checkGroup(check01{}))
	ckList.Add(checkGroup(check02{}))
	ckList.Add(checkGroup(check03{}))
	ckList.Add(checkGroup(check04{}))
	ckList.Add(checkGroup(check05{}))
}
