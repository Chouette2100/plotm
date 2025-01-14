// Copyright ©2013 The Gonum Authors. All rights reserved.
// Copyright (c) 2024 chouette2100@gmail.com
// Refer to https://github.com/gonum/plot/blob/v0.15.0/LICENSE
package main

import (
	"log"
	"strconv"
	"time"

	"gonum.org/v1/plot"
)

// UTCUnixTime converts a float64 value into a time.Time.
func JSTUnixTime(t float64) time.Time {
	return time.Unix(int64(t+60*60*9), 0).UTC()
}

// TimeTicks is suitable for axes representing time values.
type MyTimeTicks struct {
	// Ticker is used to generate a set of ticks.
	// If nil, DefaultTicks will be used.
	Ticker plot.Ticker

	// Format is the textual representation of the time value.
	// If empty, time.RFC3339 will be used
	Format string

	// Time takes a float64 value and converts it into a time.Time.
	// If nil, UTCUnixTime is used.
	Time func(t float64) time.Time
}

// var _ Ticker = MyTimeTicks{}

// Ticks implements plot.Ticker.
func (t MyTimeTicks) Ticks(min, max float64) []plot.Tick {
	//	if t.Ticker == nil {
	//		t.Ticker = plot.DefaultTicks{}
	//	}
	if t.Format == "" {
		t.Format = time.RFC3339
	}
	if t.Time == nil {
		t.Time = JSTUnixTime
	}

	// ticks := t.Ticker.Ticks(min, max)
	ticks := NewTicks(min, max)
	if ticks == nil {
		return nil
	}
	for i := range ticks {
		tick := &ticks[i]
		if tick.Label == "" {
			tick.Label = " "
			continue
		}
		tick.Label = t.Time(tick.Value).Format(t.Format)
	}
	return ticks
}

// Ticks returns Ticks in the specified range.
func NewTicks(min, max float64) []plot.Tick {
	if max <= min || min < 1727740800.0 { // 2024-10-01 00:00 UTC
		log.Printf("NewTicks(): illegal range: [%g, %g]\n", min, max)
		return nil
	}
	var ticks []plot.Tick

	// id := 60.0 * 60.0 * 3.0
	//id := (max - min ) / 10.0
	fint, idiv, _ := MakeScale(max, min)
	id := fint

	// ii := int(min/id+0.01)
	ii := 0
	for i := min; i <= max; i += id {
		// if ii%4 == 1 {
		if ii%idiv == 0 {
			ticks = append(ticks, plot.Tick{Value: i, Label: strconv.FormatFloat(i, 'f', 0, 64)})
		} else {
			ticks = append(ticks, plot.Tick{Value: i, Label: ""})
		}
		ii++
	}
	return ticks
}

func MakeScale(
	max float64,
	min float64,
) (
	fint float64,
	idiv int,
	err error,
) {

	const T08days = 60.0*60.0*24.0*8.0 + 0.001
	const T04days = 60.0*60.0*24.0*4.0 + 0.001
	const T02days = 60.0*60.0*24.0*2.0 + 0.001
	const T24hours = 60.0*60.0*24.0 + 0.001
	const T12hours = 60.0*60.0*12.0 + 0.001
	const T06hours = 60.0*60.0*6.0 + 0.001
	const T03hours = 60.0*60.0*3.0 + 0.001
	const T02hours = 60.0*60.0*2.0 + 0.001
	const T60minutes = 60.0*60.0 + 0.001
	const T30minutes = 60.0*30.0 + 0.001
	const T20minutes = 60.0*20.0 + 0.001
	const T10minutes = 60.0*10.0 + 0.001
	const T04minutes = 60.0*4.0 + 0.001

	prd := max - min

	switch {
	case prd < T04minutes:
		fint = 60.0
		idiv = 2
	case prd < T10minutes:
		fint = 60.0
		idiv = 5
	case prd < T20minutes:
		fint = 60.0 * 2.0
		idiv = 5
	case prd < T30minutes:
		fint = 60.0 * 5.0
		idiv = 3
	case prd < T60minutes:
		fint = 60.0 * 10.0
		idiv = 3
	case prd < T02hours:
		fint = 60.0 * 15.0
		idiv = 2
	case prd < T03hours:
		fint = 60.0 * 15.0
		idiv = 4
	case prd < T06hours:
		fint = 60.0 * 30.0
		idiv = 4
	case prd < T12hours:
		fint = 60.0 * 60.0 * 1.0
		idiv = 3
	case prd < T24hours:
		fint = 60.0 * 60.0 * 3.0
		idiv = 2
	case prd < T02days:
		fint = 60.0 * 60.0 * 6.0
		idiv = 2
	case prd < T04days:
		fint = 60.0 * 60.0 * 12.0
		idiv = 2
	case prd < T08days:
		fint = 60.0 * 60.0 * 24.0 * 1.0
		idiv = 2
	default:
	}
	return
}
