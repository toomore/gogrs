package main

import (
	"github.com/toomore/gogrs/twse"
	"github.com/toomore/gogrs/utils"
)

type checkGroup interface {
	String() string
	CheckFunc(...*twse.Data) bool
}

type check01 struct{}

func (check01) String() string {
	return "MA 3 > 6 > 18"
}

func (check01) CheckFunc(b ...*twse.Data) bool {
	defer wg.Done()
	var d = b[0]
	var start = d.Len()
	if start == 0 {
		d.Get()
	}
	for {
		if d.Len() >= 18 {
			break
		}
		d.PlusData()
		if (d.Len() - start) == 0 {
			break
		}
		start = d.Len()
	}
	if d.Len() < 18 {
		return false
	}
	var ma3 = d.MA(3)
	if days, ok := utils.CountCountineFloat64(utils.DeltaFloat64(ma3)); !ok || days == 0 {
		return false
	}
	var ma6 = d.MA(6)
	if days, ok := utils.CountCountineFloat64(utils.DeltaFloat64(ma6)); !ok || days == 0 {
		return false
	}
	var ma18 = d.MA(18)
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
func (check02) CheckFunc(b ...*twse.Data) bool {
	defer wg.Done()
	return utils.ThanSumPastUint64((*b[0]).GetVolumeList(), 3, true) && (*b[0]).IsRed()
}

type check03 struct{}

func (check03) String() string {
	return "量或價走平 45 天"
}

func (check03) CheckFunc(b ...*twse.Data) bool {
	defer wg.Done()
	if b[0].Len() < 45 {
		start := b[0].Len()
		for {
			b[0].PlusData()
			if b[0].Len() > 45 {
				break
			}
			if b[0].Len() == start {
				break
			}
			start = b[0].Len()
		}
		if b[0].Len() < 45 {
			return false
		}
	}
	var price = b[0].GetPriceList()
	var volume = b[0].GetVolumeList()
	return price[len(price)-1] > 10 &&
		(utils.SD(price[len(price)-46:]) < 0.25 ||
			utils.SDUint64(volume[len(volume)-46:]) < 0.25)
}

func init() {
	ckList.Add(checkGroup(check01{}))
	ckList.Add(checkGroup(check02{}))
	ckList.Add(checkGroup(check03{}))
}
