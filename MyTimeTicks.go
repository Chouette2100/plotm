package main

import (
	"fmt"
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
		Graphd.Err = fmt.Errorf("illegal range: [%g, %g]", min, max)
		return nil
	}
	var ticks []plot.Tick

	// id := 60.0 * 60.0 * 3.0
	//id := (max - min ) / 10.0
	id := Graphd.Fint

	// ii := int(min/id+0.01)
	ii := 0
	for i := min; i <= max; i += id {
		// if ii%4 == 1 {
		if ii%Graphd.Div == 0 {
			ticks = append(ticks, plot.Tick{Value: i, Label: strconv.FormatFloat(i, 'f', 0, 64)})
		} else {
			ticks = append(ticks, plot.Tick{Value: i, Label: ""})
		}
		ii++
	}
	return ticks
}
