// Copyright ©2013 The Gonum Authors. All rights reserved.
// Copyright (c) 2024 chouette2100@gmail.com
// Refer to https://github.com/gonum/plot/blob/v0.15.0/LICENSE
package main

import (
	"fmt"
	"math"
	"strconv"

	"gonum.org/v1/plot"
)

// MyTicks
type MyTicks struct{}

// Ticks returns Ticks in the specified range.
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
		{3, 13, 4, 0.25},
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
