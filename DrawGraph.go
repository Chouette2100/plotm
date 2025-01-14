// Copyright © 2024 chouette2100@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php
package main

import (
	"fmt"
	"log"
	"math"
	"os"
	// "strings"
	// "time"

	// "github.com/go-gorp/gorp"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"image/color"
	// "gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

type Device struct {
	Device    int
	Name      string
	Tabletype interface{}
	Tablename string
}

func DrawGraph(graphd Graph) (
	filename string,
	err error,
) {
	rs := make([][]interface{}, len(graphd.Device))

	// for k := 0; k < len(graphd.Device); k++ {
	for k, dev := range graphd.Device {
		// if strings.Contains(graphd.Device[k].Name, "AHT10") {
		// 	rs[k], err = Dbmap.Select(Aht10{}, sqlaht10, graphd.Device[k].Device, graphd.Btime, graphd.Etime)
		// } else {
		// 	rs[k], err = Dbmap.Select(Scd41{}, sqlscd41, graphd.Device[k].Device, graphd.Btime, graphd.Etime)
		// }
		sqlst := "SELECT * FROM showroom." + dev.Tablename + " WHERE device = ? "
		sqlst += " AND ts BETWEEN ? AND ? "
		sqlst += " ORDER BY ts "
		rs[k], err = Dbmap.Select(dev.Tabletype, sqlst, graphd.Device[k].Device, graphd.Btime, graphd.Etime)
		if err != nil {
			log.Printf("Database error. err=%s.\n", err.Error())
			err = fmt.Errorf("Dbmap.Select() %w", err)
			return
		}
	}
	const cols = 1
	rows := len(graphd.Item)
	plots := make([][]*plot.Plot, rows)

	// xticks := MyTimeTicks{Format: "01-02\n15:04"}
	// if graphd.Err != nil {
	// 	err = fmt.Errorf("MyTimeTicks(): %w", graphd.Err)
	// 	return
	// }

	for j, it := range graphd.Item {
		dvmax := len(graphd.Item[j].Udev)
		plots[j] = make([]*plot.Plot, cols)
		for i := 0; i < cols; i++ {

			// Create a new plot, set its title and
			// axis labels.
			p := plot.New()

			p.Y.Label.Text = it.Name + "[" + it.Unit + "]"
			if j == len(graphd.Item)-1 {
				p.X.Label.Text = "Date and Time"
			}
			// }
			p.Y.Tick.Marker = MyTicks{}
			if graphd.Err != nil {
				err = fmt.Errorf("MyTicks(): %w", graphd.Err)
				return
			}

			unt := int64(graphd.Fint)
			p.X.Max = float64(graphd.Etime.Unix() / unt * unt)
			p.X.Min = float64(graphd.Btime.Unix() / unt * unt)
			xticks := MyTimeTicks{Format: "01-02\n15:04"}
			if graphd.Err != nil {
				err = fmt.Errorf("MyTimeTicks(): %w", graphd.Err)
				return
			}
			p.X.Tick.Marker = xticks

			// Draw a grid behind the data
			p.Add(plotter.NewGrid())

			type pts []plotter.XYs
			ptslist := make([]pts, dvmax)

			// for k := 0; k < dvmax; k++ {
			for k, vd := range graphd.Item[j].Udev {
				// if graphd.Item[j].Udev[k].Status != "checked" {
				// 	continue
				// }
				if vd.Status != "checked" || len(rs[vd.Devno]) == 0 {
					continue
				}
				// pts[k] = make(plotter.XYs, len(rs[graphd.Item[j].Udev[k].Devno]))
				ptslist[k] = make(pts, 0)
			}
			vmax := it.Vmax - 0.00001
			vmin := it.Vmin
			vint := it.Vint
			for k, vd := range graphd.Item[j].Udev {
				// if graphd.Item[j].Udev[k].Status != "checked" { // 	continue
				// }
				if vd.Status != "checked" || len(rs[vd.Devno]) == 0 {
					continue
				}
				ptslist[k] = make(pts, 1)
				ptslist[k][0] = make(plotter.XYs, len(rs[vd.Devno]))
				tlast := int64(0)
				if vd.Devno < 2 {
					tlast = rs[vd.Devno][0].(*Aht10).Ts.Unix()
				} else {
					tlast = rs[vd.Devno][0].(*Scd41).Ts.Unix()
				}

				lenrs := 0
				ip := 0
				im := 0
				m := 0
				var r interface{}
				for m, r = range rs[vd.Devno] {
					var v float64
					v, err = GetMeasurementResults(graphd, it.Name, graphd.Method, vd.Devno, r)
					if err != nil {
						fmt.Printf("GetMeasurementResults() err=%s\n", err.Error())
						return
					}

					tnew := int64(0)
					if vd.Devno < 2 {
						tnew = r.(*Aht10).Ts.Unix()
					} else {
						tnew = r.(*Scd41).Ts.Unix()
					}
					// if tnew-tlast > 900 {
					if tnew-tlast > 900 || v > vmax || v < vmin {
						// No line is drawn where it is missing.
						// ptslist[k][ip] = ptslist[k][ip][:m-lenrs]
						if im != 0 {
							ptslist[k][ip] = ptslist[k][ip][:im-1]
							lenrs = m
							ptslist[k] = append(ptslist[k], make(plotter.XYs, len(rs[vd.Devno])-lenrs))
							ip++
							im = 0
						}
					}
					tlast = tnew
					// log.Printf("k=%d, m=%d, ip=%d, lenrs=%d\n", k, m, ip, lenrs)
					if v <= vmax && v >= vmin {
						// ptslist[k][ip][m-lenrs].X = float64(tnew)
						// ptslist[k][ip][m-lenrs].Y = v
						if tnew < 1 {
							err = fmt.Errorf(" j=%d, i=%d, k=%d, m=%d, ip=%d, im=%d, tnew=%v", j, i, k, m, ip, im, tnew)
							return
						}
						if v < 0 {
							err = fmt.Errorf(" j=%d, i=%d, k=%d, m=%d, ip=%d, im=%d, v=%f", j, i, k, m, ip, im, v)
							return
						}
						ptslist[k][ip][im].X = float64(tnew)
						ptslist[k][ip][im].Y = v
						im++
					}
					// if v > vmax {
					// 	vmax = v
					// } else if v < vmin {
					// 	vmin = v
					// }
				}
				// ptslist[k][ip] = ptslist[k][ip][:m-lenrs]
				ptslist[k][ip] = ptslist[k][ip][:im]
			}
			// make sure the horizontal scales match
			p.Y.Max = float64(int(vmax)/vint*vint + vint)
			p.Y.Min = float64(int(vmin) / vint * vint)
			// p.Y.Max = it.Vmax
			// p.Y.Min = it.Vmin
			// p.Y.AutoRescale = false

			// Make a line plotter and set its style.
			llist := make([][]*plotter.Line, dvmax)
			// l := make([]*plotter.Line, dvmax)

			// for k := 0; k < dvmax; k++ {
			for k, vd := range graphd.Item[j].Udev {
				// if graphd.Item[j].Udev[k].Status != "checked" {
				// 	continue
				// }
				if vd.Status != "checked" || len(rs[vd.Devno]) == 0 {
					continue
				}

				llist[k] = make([]*plotter.Line, len(ptslist[k]))
				ip := 0
				for ; ip < len(ptslist[k]); ip++ {
					// if graphd.Item[j].Udev[k].Status != "checked" {
					// 	continue
					// }
					if vd.Status != "checked" || len(rs[vd.Devno]) == 0 {
						continue
					}
					llist[k][ip], err = plotter.NewLine(ptslist[k][ip])
					if err != nil {
						err = fmt.Errorf("plotter.NewLine() %w", err)
						return
					}
					llist[k][ip].LineStyle.Width = vg.Points(1)

					llist[k][ip], err = plotter.NewLine(ptslist[k][ip])
					if err != nil {
						err = fmt.Errorf("plotter.NewLine() %w", err)
						return
					}
					llist[k][ip].LineStyle.Width = vg.Points(1)

					switch graphd.Item[j].Udev[k].Devno {
					case 0:
						llist[k][ip].LineStyle.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255}
					case 1:
						llist[k][ip].LineStyle.Color = color.RGBA{R: 0, G: 255, B: 0, A: 255}
					case 2:
						llist[k][ip].LineStyle.Color = color.RGBA{R: 0, G: 0, B: 255, A: 255}
					case 3:
						// llist[k][ip].LineStyle.Color = color.RGBA{R: 0, G: 95, B: 95, A: 191}
						llist[k][ip].LineStyle.Color = color.RGBA{R: 0, G: 128, B: 0, A: 255}
					}

					// Add the plotters to the plot, with a legend
					// entry for each

					p.Add(llist[k][ip]) // TOFIX:
				}
				if ip != 0 {
					p.Legend.Add(graphd.Device[graphd.Item[j].Udev[k].Devno].Name, llist[k][ip-1])
				}
			}

			plots[j][i] = p
		}
	}

	img := vgimg.New(vg.Points(800), vg.Points(500))
	dc := draw.New(img)

	t := draw.Tiles{
		Rows: rows,
		Cols: cols,
	}

	canvases := plot.Align(plots, t, dc)
	for j := 0; j < rows; j++ {
		for i := 0; i < cols; i++ {
			if plots[j][i] != nil {
				plots[j][i].Draw(canvases[j][i])
			}
		}
	}
	var w *os.File
	filename = fmt.Sprintf("TAH%04d.png", <-Chimgfn)
	w, err = os.Create("public/" + filename)
	if err != nil {
		err = fmt.Errorf("cannot create file %s", filename)
		return
	}

	png := vgimg.PngCanvas{Canvas: img}
	if _, err = png.WriteTo(w); err != nil {
		err = fmt.Errorf("cannot convert to image file")
		return
	}
	return
}

// Saturated vapor pressure of water at t°C
func CalcE(t float64) (e float64) {
	e = 6.1078 * math.Pow(10.0, 7.5*t/(t+237.3)) / 10.0
	return
}

// Saturated water vapor content at t°C
func CalcSVC(t float64) (vh float64) {
	vh = CalcE(t) * 10.0 / (t + 273.5) * 216.7
	return
}

func GetMeasurementResults(graphd Graph, item string, method string, devno int, r interface{}) (vm float64, err error) {
	switch item {
	case "Temperature":
		switch graphd.Device[devno].Tabletype.(type) {
		case Aht10:
			vm = r.(*Aht10).Temperature
		case Scd41:
			vm = r.(*Scd41).Temperature
		default:
			err = fmt.Errorf("unknown device type. check tabletype of device. devno=%d", devno)
			return
		}
	case "Humidity":
		if method == "R" {
			// Relative humidity
			switch graphd.Device[devno].Tabletype.(type) {
			case Aht10:
				vm = r.(*Aht10).Humidity
			case Scd41:
				vm = r.(*Scd41).Humidity
			default:
				err = fmt.Errorf("unknown device type. check tabletype of device. devno=%d", devno)
				return
			}
		} else {

			// absolute humidity by volume
			switch graphd.Device[devno].Tabletype.(type) {
			case Aht10:
				t := r.(*Aht10).Temperature
				rh := r.(*Aht10).Humidity
				vm = CalcSVC(t) * rh / 100
			case Scd41:
				t := r.(*Scd41).Temperature
				rh := r.(*Scd41).Humidity
				vm = CalcSVC(t) * rh / 100
			default:
				err = fmt.Errorf("unknown device type. check tabletype of device. devno=%d", devno)
				return
			}
		}
	case "CO2":
		switch graphd.Device[devno].Tabletype.(type) {
		case Scd41:
			vm = float64(r.(*Scd41).Co2)
		default:
			err = fmt.Errorf("unknown device type. check tabletype of device. devno=%d", devno)
			return
		}
	}
	return
}
