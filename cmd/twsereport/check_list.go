package main

import "github.com/toomore/gogrs/utils"

type base interface {
	MA(days int) []float64
	Len() int
	Get() ([][]string, error)
	PlusData()
}

type check01 struct{}

func (check01) String() string {
	return "MA 3 > 6 > 18"
}

func (check01) CheckFunc(b ...base) bool {
	defer wg.Done()
	var d = b[0]
	var start = d.Len()
	d.Get()
	for {
		if d.Len() >= 18 {
			break
		}
		d.PlusData()
		if (d.Len() - start) == 0 {
			break
		}
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
	return "check02"
}
func (check02) CheckFunc(b ...base) bool {
	defer wg.Done()
	if b[0].Len() < 3 {
		return false
	}
	days, up := utils.CountCountineFloat64(utils.DeltaFloat64(b[0].MA(3)))
	if up && days > 1 {
		return true
	}
	return false
}

func init() {
	ckList.Add(checkGroup(check01{}))
	ckList.Add(checkGroup(check02{}))
}
