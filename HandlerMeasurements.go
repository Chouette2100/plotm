// Copyright © 2024 chouette2100@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php
package main

import (
	"fmt"
	"log"
	// "os"
	"strconv"
	"time"

	"html/template"
	"net/http"
)

type Udevice struct {
	Devno  int
	Name   string
	Status string
}

type Item struct {
	Name string
	Unit string
	Vmin float64
	Vmax float64
	Vint int
	Udev []Udevice
}

type Graph struct {
	Btime    time.Time
	Etime    time.Time
	Uetime   int64
	Item     []Item
	Device   []Device
	Filename string
	Nterm    int64
	Period   string
	Fperiod  float64
	Fint     float64
	Div      int
	Method   string
	Err      error
}

var Graphd Graph

func HandlerMeasurements(w http.ResponseWriter, r *http.Request) {


	_, _, isallow := GetUserInf(r)
	if !isallow {
		w.Write([]byte("Access Denied\n"))
		return
	}


	Graphd.Device = []Device{
		{0, "AHT10 X038", Aht10{}, "aht10"},
		{1, "AHT10 X039", Aht10{}, "aht10"},
		{98, "SCD41 single shot", Scd41{}, "scd41"},
		{4096, "SCD41 periodic", Scd41{}, "scd41"},
	}

	var err error

	snterm := r.FormValue("nterm")
	if snterm == "" {
		Graphd.Nterm = 0
	} else {
		intterm, _ := strconv.Atoi(snterm)
		Graphd.Nterm = int64(intterm)
	}
	Graphd.Method = r.FormValue("method")
	if Graphd.Method == "" {
		Graphd.Method = "R"
	}

	speriod := r.FormValue("period")
	Graphd.Period = speriod
	switch speriod {
	case "8 days":
		Graphd.Fperiod = 60.0 * 60.0 * 24.0 * 8.0
		Graphd.Fint = 60.0 * 60.0 * 24.0 * 1.0
		Graphd.Div = 2
	case "4 days":
		Graphd.Fperiod = 60.0 * 60.0 * 24.0 * 4.0
		Graphd.Fint = 60.0 * 60.0 * 12.0
		Graphd.Div = 2
	case "2 days":
		Graphd.Fperiod = 60.0 * 60.0 * 24.0 * 2.0
		Graphd.Fint = 60.0 * 60.0 * 6.0
		Graphd.Div = 2
	case "1 day":
		Graphd.Fperiod = 60.0 * 60.0 * 24.0
		Graphd.Fint = 60.0 * 60.0 * 3.0
		Graphd.Div = 2
	case "12 hours":
		Graphd.Fperiod = 60.0 * 60.0 * 12.0
		Graphd.Fint = 60.0 * 60.0 * 1.0
		Graphd.Div = 3
	case "6 hours":
		Graphd.Fperiod = 60.0 * 60.0 * 6.0
		Graphd.Fint = 60.0 * 60.0 * 1.0
		Graphd.Div = 2
	case "3 hours":
		Graphd.Fperiod = 60.0 * 60.0 * 3.0
		Graphd.Fint = 60.0 * 20.0
		Graphd.Div = 3
	case "2 hours":
		Graphd.Fperiod = 60.0 * 60.0 * 2.0
		Graphd.Fint = 60.0 * 15.0
		Graphd.Div = 2
	case "1 hour":
		Graphd.Fperiod = 60.0 * 60.0
		Graphd.Fint = 60.0 * 10.0
		Graphd.Div = 2
	case "30 minutes":
		Graphd.Fperiod = 60.0 * 20.0
		Graphd.Fint = 60.0 * 2.0
		Graphd.Div = 2
	case "15 minutes":
		Graphd.Fperiod = 60.0 * 15.0
		Graphd.Fint = 60.0
		Graphd.Div = 5
	case "10 minutes":
		Graphd.Fperiod = 60.0 * 10.0
		Graphd.Fint = 60.0
		Graphd.Div = 5
	case "5 minutes":
		Graphd.Fperiod = 60.0 * 10.0
		Graphd.Fint = 60.0
		Graphd.Div = 5
	default:
		Graphd.Period = "2 days"
		Graphd.Fperiod = 60.0 * 60.0 * 24.0 * 2.0
		Graphd.Fint = 60.0 * 60.0 * 6.0
		Graphd.Div = 2
	}

	seutime := r.FormValue("uetime")
	if seutime == "" {
		Graphd.Etime = time.Now()
		Graphd.Uetime = Graphd.Etime.Unix()
	} else {
		iuetime, _ := strconv.Atoi(seutime)
		Graphd.Uetime = int64(iuetime)
		Graphd.Etime = time.Unix(Graphd.Uetime, 0)
	}

	fnc := r.FormValue("fnc")

	switch fnc {
	case "P":
		// Graphd.Nterm += 1
		Graphd.Etime = Graphd.Etime.Add(-time.Duration(Graphd.Fperiod/2.0) * time.Second)
	case "N":
		// Graphd.Nterm -= 1
		// if Graphd.Nterm < 0 {
		// 	Graphd.Nterm = 0
		// }
		Graphd.Etime = Graphd.Etime.Add(time.Duration(Graphd.Fperiod/2.0) * time.Second)
	case "L", "":
		// Graphd.Nterm = 0
		// unt := int64(Graphd.Fint) * int64(Graphd.Div)
		// Graphd.Etime = time.Unix(time.Now().Unix()/unt*unt+unt-unt*Graphd.Nterm, 0)
		Graphd.Etime = time.Now()
	case "R":
	default:
	}
	unt := int64(Graphd.Fint)
	Graphd.Etime = time.Unix((Graphd.Etime.Unix()-1)/unt*unt+unt, 0)
	Graphd.Btime = Graphd.Etime.Add(time.Duration(-Graphd.Fperiod) * time.Second)
	Graphd.Uetime = Graphd.Etime.Unix()

	Graphd.Item = make([]Item, 3)
	Graphd.Item[0] = Item{"Temperature", "°C", 15.0, 25.0, 5,
		[]Udevice{
			{0, "", "checked"},
			{1, "", ""},
			{3, "", ""},
			{2, "", ""},
		}}
	Graphd.Item[1] = Item{"Humidity", "%RH", 30.0, 60.0, 10,
		[]Udevice{
			{0, "", "checked"},
			{1, "", ""},
			{3, "", ""},
			{2, "", ""},
		}}

	if Graphd.Method == "V" {
		Graphd.Item[1].Name = "VH"
		Graphd.Item[1].Unit = "g/m3"
		Graphd.Item[1].Vmin = 0
		Graphd.Item[1].Vmax = 10
		Graphd.Item[1].Vint = 10
	}
	Graphd.Item[2] = Item{"CO2", "ppm", 0, 4000, 1000,
		[]Udevice{
			{3, "", "checked"},
			{2, "", ""},
		}}

	for j := 0; j < len(Graphd.Item); j++ {
		for k := 0; k < len(Graphd.Item[j].Udev); k++ {
			Graphd.Item[j].Udev[k].Name = Graphd.Device[Graphd.Item[j].Udev[k].Devno].Name
		}

	}

	// tnow := time.Now().Truncate(time.Second)
	// // etime := tnow.Add(9 * time.Hour).Truncate(24 * time.Hour).Add(15 * time.Hour)
	// // ih := int64(60 * 60)
	// // etime := time.Unix((tnow.Unix()+(ih*9))/(ih*12)*(ih*12)+(ih*3)-ih*24*graph.Term, 0)
	// // btime := etime.Add(-48 * time.Hour)
	// unt := int64(Graphd.Fint) * int64(Graphd.Div)
	// Graphd.Etime = time.Unix(tnow.Unix()/unt*unt+unt-unt*Graphd.Nterm, 0)
	// Graphd.Btime = Graphd.Etime.Add(time.Duration(-Graphd.Fperiod) * time.Second)
	// Graphd.Uetime = Graphd.Etime.Unix()

	if fnc != "" {
		for j := 0; j < len(Graphd.Item); j++ {
			for i := range Graphd.Item[j].Udev {
				devname := "dev" + strconv.Itoa(j) + "_" + strconv.Itoa(i)
				checked := r.FormValue(devname)
				if checked == "checked" {
					Graphd.Item[j].Udev[i].Status = "checked"
				} else {
					Graphd.Item[j].Udev[i].Status = ""
				}
			}
		}
	}
	Graphd.Filename, err = DrawGraph()
	if err != nil {
		fmt.Printf("DrawGraph() error=%s\n", err.Error())
		err = fmt.Errorf("DrawGraph(): %w", err)
		w.Write([]byte(err.Error()))
		return
	}

	// Graphd.Filename = "/" + Graphd.Filename

	// テンプレートをパースする
	tpl := template.Must(template.ParseFiles("Measurements.gtpl"))

	// テンプレートに出力する値をマップにセット
	/*
	   values := map[string]string{
	           "filename": req.FormValue("FileName"),
	   }
	*/

	// マップを展開してテンプレートを出力する
	if err := tpl.ExecuteTemplate(w, "Measurements.gtpl", Graphd); err != nil {
		log.Printf("tpl.ExcecuteTemplate(): %s", err.Error())
	}

}
