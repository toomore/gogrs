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
	defer wg.Done()
	if prepareData(b...)[0] != true {
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
	defer wg.Done()
	return utils.ThanSumPastUint64((*b[0]).GetVolumeList(), 3, true) && (*b[0]).IsRed()
}

type check03 struct{}

func (check03) String() string {
	return "量或價走平 45 天"
}

func (check03) Mindata() int {
	return 45
}

func (check03) CheckFunc(b ...*twse.Data) bool {
	defer wg.Done()
	if !prepareData(b...)[0] {
		return false
	}
	var price = b[0].GetPriceList()
	var volume = b[0].GetVolumeList()
	return price[len(price)-1] > 10 &&
		(utils.SD(price[len(price)-46:]) < 0.25 ||
			utils.SDUint64(volume[len(volume)-46:]) < 0.25)
}

func prepareData(b ...*twse.Data) []bool {
	var result []bool
	var mindata int
	for i := range ckList {
		if ckList[i].Mindata() > mindata {
			mindata = ckList[i].Mindata()
		}
	}

	for i, _ := range b {
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
					result[i] = false
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
}
