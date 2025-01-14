// Copyright ©2013 The Gonum Authors. All rights reserved.
// Copyright (c) 2024 chouette2100@gmail.com
// Refer to https://github.com/gonum/plot/blob/v0.15.0/LICENSE
package main

import (
	"log"
	"math"
	"strconv"

	"gonum.org/v1/plot"
)

// MyTicks
type MyTicks struct{}

// Ticks returns Ticks in the specified range.
func (MyTicks) Ticks(min, max float64) []plot.Tick {
	if max <= min {
		log.Printf("(MyTicks)Ticks(): illegal range [%f, %f]\n", min, max)
		return nil
	}
	var ticks []plot.Tick

	dt := max - min
	ndg := int(math.Log10(dt + 0.00001))
	p10 := math.Pow10(ndg)
	frstdg := (dt + 0.00001) / p10

	ie := 0
	is := 0
	st := 0.0

	switch {
	case frstdg < 1.001:
		ie = 11
		is = 5
		st = 0.1
	case frstdg < 1.501:
		ie = 7
		is = 2
		st = 0.25
	case frstdg < 2.001:
		ie = 11
		is = 5
		st = 0.2
	case frstdg < 2.501:
		ie = 11
		is = 2
		st = 0.25
	case frstdg < 3.001:
		ie = 7
		is = 2
		st = 0.5
	case frstdg < 4.001:
		ie = 9
		is = 2
		st = 0.5
	case frstdg < 5.001:
		ie = 11
		is = 2
		st = 0.5
	case frstdg < 6.001:
		ie = 13
		is = 2
		st = 0.5
	case frstdg < 7.001:
		ie = 8
		is = 1
		st = 1.0
	case frstdg < 8.001:
		ie = 17
		is = 2
		st = 0.5
	case frstdg < 10.0:
		ie = 11
		is = 5
		st = 1.0
	}

	st *= p10

	v := float64(int((min+0.001)/st)) * st
	for i := 0; i < ie; i++ {
		if i%is == 0 {
			ticks = append(ticks, plot.Tick{Value: v, Label: strconv.FormatFloat(v, 'f', 0, 64)})
		} else {
			ticks = append(ticks, plot.Tick{Value: v, Label: " "})
		}
		v += st
	}
	return ticks
}
/*
func (MyTicks) Ticks(min, max float64) []plot.Tick {
	if max <= min {
		Graphd.Err = fmt.Errorf("illegal range [%f, %f]", min, max)
		return nil
	}
	var ticks []plot.Tick

	dt := max - min
	ndg := int(math.Log10(dt + 0.00001))
	p10 := math.Pow10(ndg)
	frstdg := int((dt + 0.00001) / p10)

	type Vtck struct {
		Ib int
		Ie int
		Is int
		St float64
	}
	var vtcklist [10]Vtck = [10]Vtck{
		{0, 11, 5, 0.1},
		{1, 11, 5, 0.1},
		{2, 11, 5, 0.2},
		{3, 15, 5, 0.2},
		{4, 9, 2, 0.5},
		{5, 11, 2, 0.5},
		{6, 13, 2, 0.5},
		{8, 17, 2, 0.5},
		{8, 17, 2, 0.5},
		{10, 11, 5, 1.0},
	}

	ie := vtcklist[frstdg].Ie
	is := vtcklist[frstdg].Is
	st := vtcklist[frstdg].St * p10

	v := float64(int(min/st)) * st
	for i := 0; i < ie; i++ {
		if i%is == 0 {
			ticks = append(ticks, plot.Tick{Value: v, Label: strconv.FormatFloat(v, 'f', 0, 64)})
		} else {
			ticks = append(ticks, plot.Tick{Value: v, Label: " "})
		}
		v += st
	}
	return ticks
}
*/