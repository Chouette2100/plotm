// Copyright © 2024 chouette2100@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php
package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"crypto/sha256"

	"html/template"
	"net/http"

	"gopkg.in/yaml.v3"
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
	// Nterm    int64
	Period   string
	Fperiod  float64
	Fint     float64
	Div      int
	Method   string
	Rhmax    float64
	Rhmin    float64
	Vhmax    float64
	Vhmin    float64
	Ymlfiles []string
	Nextyml  string
	Err      error
}

func HandlerMeasurements(w http.ResponseWriter, r *http.Request) {

	_, _, isallow := GetUserInf(r)
	if !isallow {
		w.Write([]byte("Access Denied\n"))
		return
	}

	graphd := Graph{}
	var err error

	sess := globalSessions.SessionStart(w, r)
	cfgfn := sess.Get("cfgfn") // ②
	if cfgfn == nil {
		cfgfn = <-Chcfgfn
		sess.Set("cfgfn", cfgfn)
	} else {
		cfgfn = cfgfn.(int)
	}
	nextyml := "User_" + fmt.Sprintf("%05d.yml", cfgfn)

	sha256.Sum256([]byte(sess.SessionID()))
	tmpyml := fmt.Sprintf("%0x.yml", sha256.Sum256([]byte(sess.SessionID())))


	fnc := r.FormValue("fnc")
	yml := r.FormValue("yml")
	seutime := r.FormValue("uetime")
	if seutime != "" {
		ieutime, _ := strconv.Atoi(seutime)
		graphd.Uetime = int64(ieutime)
	}

	if fnc == "" || fnc == "RS" || graphd.Uetime == 0 {
		if yml == "" {
			yml = "Default_000.yml"
		}
		err = ReadYAML("YmlFiles/"+yml, &graphd)
		if err != nil {
			fmt.Printf("ReadYAML() error=%s\n", err.Error())
			err = fmt.Errorf("ReadYAML(): %w", err)
			w.Write([]byte(err.Error()))
			return
		}
		if graphd.Uetime == 0 {
			graphd.Etime = time.Now()
			graphd.Uetime = graphd.Etime.Unix()
		}
		for _, it := range graphd.Item {
			if it.Name == "Humidity" {
				if it.Unit == "%RH" {
					graphd.Rhmax = it.Vmax
					graphd.Rhmin = it.Vmin
				} else {
					graphd.Vhmax = it.Vmax
					graphd.Vhmin = it.Vmin
				}
			}
		}
	} else {
		err = ReadYAML("tmp/"+tmpyml, &graphd)
		if err != nil {
			err = ReadYAML("YmlFiles/Default_000.yml", &graphd)
		}
		if err != nil {
			fmt.Printf("ReadYAML() error=%s\n", err.Error())
			err = fmt.Errorf("ReadYAML(): %w", err)
			w.Write([]byte(err.Error()))
			return
		}
		speriod := r.FormValue("period")
		graphd.Period = speriod
		err = SetPeriod(&graphd)
		if err != nil {
			fmt.Printf("SetPeriod() error=%s\n", err.Error())
			err = fmt.Errorf("SetPeriod(): %w", err)
			w.Write([]byte(err.Error()))
			return
		}

		seutime := r.FormValue("uetime")
		if seutime == "" {
			graphd.Etime = time.Now()
			graphd.Uetime = graphd.Etime.Unix()
		} else {
			iuetime, _ := strconv.Atoi(seutime)
			graphd.Uetime = int64(iuetime)
			graphd.Etime = time.Unix(graphd.Uetime, 0)
		}

		switch fnc {
		case "P":
			// graphd.Nterm += 1
			graphd.Etime = graphd.Etime.Add(-time.Duration(graphd.Fperiod/2.0) * time.Second)
		case "N":
			// graphd.Nterm -= 1
			// if graphd.Nterm < 0 {
			// 	graphd.Nterm = 0
			// }
			graphd.Etime = graphd.Etime.Add(time.Duration(graphd.Fperiod/2.0) * time.Second)
		case "L", "":
			// graphd.Nterm = 0
			// unt := int64(graphd.Fint) * int64(graphd.Div)
			// graphd.Etime = time.Unix(time.Now().Unix()/unt*unt+unt-unt*graphd.Nterm, 0)
			graphd.Etime = time.Now()
		case "R":
		default:
		}
		lastmethod := graphd.Method
		graphd.Method = r.FormValue("method")
		if graphd.Method == "" {
			graphd.Method = "R"
		}
		if graphd.Method != "R" {
			graphd.Method = "V"
		}

		if fnc != "" {
			for j := 0; j < len(graphd.Item); j++ {
				rngmin := "rng_" + strconv.Itoa(j) + "_min"
				if r.FormValue(rngmin) != "" {
					min, _ := strconv.ParseFloat(r.FormValue(rngmin), 64)
					graphd.Item[j].Vmin = min
				}
				rngmax := "rng_" + strconv.Itoa(j) + "_max"
				if r.FormValue(rngmax) != "" {
					max, _ := strconv.ParseFloat(r.FormValue(rngmax), 64)
					graphd.Item[j].Vmax = max
				}
				if graphd.Item[j].Name == "Humidity" && graphd.Method == lastmethod {
					if graphd.Method == "R" {
						graphd.Rhmax = graphd.Item[j].Vmax
						graphd.Rhmin = graphd.Item[j].Vmin
					} else {
						graphd.Vhmax = graphd.Item[j].Vmax
						graphd.Vhmin = graphd.Item[j].Vmin
					}
				}
				for i := range graphd.Item[j].Udev {
					devname := "dev" + strconv.Itoa(j) + "_" + strconv.Itoa(i)
					checked := r.FormValue(devname)
					if checked == "checked" {
						graphd.Item[j].Udev[i].Status = "checked"
					} else {
						graphd.Item[j].Udev[i].Status = ""
					}
				}
			}
		}

	}
	graphd.Nextyml = nextyml

	for nd, dv := range graphd.Device {
		switch dv.Tablename {
		case "aht10":
			graphd.Device[nd].Tabletype = Aht10{}
		case "scd41":
			graphd.Device[nd].Tabletype = Scd41{}
		}
	}

	/*
		graphd.Device = []Device{
			{0, "AHT10 X038", Aht10{}, "aht10"},
			{1, "AHT10 X039", Aht10{}, "aht10"},
			{98, "SCD41 single shot", Scd41{}, "scd41"},
			{4096, "SCD41 periodic", Scd41{}, "scd41"},
		}
	*/

	// snterm := r.FormValue("nterm")
	// if snterm == "" {
	// 	graphd.Nterm = 0
	// } else {
	// 	intterm, _ := strconv.Atoi(snterm)
	// 	graphd.Nterm = int64(intterm)
	// }

	unt := int64(graphd.Fint)
	graphd.Etime = time.Unix((graphd.Etime.Unix()-1)/unt*unt+unt, 0)
	graphd.Btime = graphd.Etime.Add(time.Duration(-graphd.Fperiod) * time.Second)
	graphd.Uetime = graphd.Etime.Unix()

	/*
		graphd.Item = make([]Item, 3)
		graphd.Item[0] = Item{"Temperature", "°C", 15.0, 30.0, 5,
			[]Udevice{
				{0, "", "checked"},
				{1, "", ""},
				{3, "", ""},
				{2, "", ""},
			}}
		graphd.Item[1] = Item{"Humidity", "%RH", 20.0, 60.0, 10,
			[]Udevice{
				{0, "", "checked"},
				{1, "", ""},
				{3, "", ""},
				{2, "", ""},
			}}
	*/

	for j := 0; j < len(graphd.Item); j++ {
		if graphd.Item[j].Name == "Humidity" {
			if graphd.Method == "V" {
				graphd.Item[j].Name = "Humidity"
				graphd.Item[j].Unit = "g/m3"
				graphd.Item[j].Vmin = graphd.Vhmin
				graphd.Item[j].Vmax = graphd.Vhmax
				graphd.Item[j].Vint = 1
			} else {
				graphd.Item[j].Name = "Humidity"
				graphd.Item[j].Unit = "%RH"
				graphd.Item[j].Vmin = graphd.Rhmin
				graphd.Item[j].Vmax = graphd.Rhmax
				graphd.Item[j].Vint = 1
			}
		}
	}
	/*
		graphd.Item[2] = Item{"CO2", "ppm", 0, 4000, 1000,
			[]Udevice{
				{3, "", "checked"},
				{2, "", ""},
			}}

		for j := 0; j < len(graphd.Item); j++ {
			for k := 0; k < len(graphd.Item[j].Udev); k++ {
				graphd.Item[j].Udev[k].Name = graphd.Device[graphd.Item[j].Udev[k].Devno].Name
			}

		}
	*/

	graphd.Filename, err = DrawGraph(graphd)
	if err != nil {
		fmt.Printf("DrawGraph() error=%s\n", err.Error())
		err = fmt.Errorf("DrawGraph(): %w", err)
		w.Write([]byte(err.Error()))
		return
	}

	// テンプレートをパースする
	tpl := template.Must(template.ParseFiles("Measurements.gtpl"))

	if fnc == "SV" {
		current := r.FormValue("current")
		if current == "yes" {
			graphd.Uetime = 0
		}
		err = WriteYAML("YmlFiles/"+graphd.Nextyml, graphd)
		if err != nil {
			fmt.Printf("WriteYAML() error=%s\n", err.Error())
			err = fmt.Errorf("WriteYAML(): %w", err)
			w.Write([]byte(err.Error()))
			return
		}
		graphd.Nextyml = "User_" + fmt.Sprintf("%05d.yml", <-Chcfgfn)
	}

	files, _ := os.ReadDir("YmlFiles")
	graphd.Ymlfiles = make([]string, len(files))
	user := "User_00000.yml"
	i := 0
	for _, f := range files {
		fname := f.Name()
		graphd.Ymlfiles[i] = fname
		if fname[0:5] == "User_" && fname > user {
			user = fname
		}
		i++
	}
	graphd.Ymlfiles = graphd.Ymlfiles[:i]

	if err := tpl.ExecuteTemplate(w, "Measurements.gtpl", graphd); err != nil {
		log.Printf("tpl.ExcecuteTemplate(): %s", err.Error())
	}

	err = WriteYAML("tmp/"+tmpyml, graphd)
	if err != nil {
		fmt.Printf("WriteYAML() error=%s\n", err.Error())
		err = fmt.Errorf("WriteYAML(): %w", err)
		w.Write([]byte(err.Error()))
		return
	}

}
func WriteYAML(fn string, yf interface{}) error {
	file, err := os.Create(fn)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := yaml.NewEncoder(file)
	defer encoder.Close()
	err = encoder.Encode(yf)
	if err != nil {
		return err
	}
	return nil
}
func ReadYAML(fn string, yf interface{}) error {
	file, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer file.Close()
	byteData, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(byteData, yf) //メモリ上のYAMLを構造体へ
	if err != nil {
		return err
	}
	return nil
}

func SetPeriod(graphd *Graph) (err error) {
	switch graphd.Period {
	case "8 days":
		graphd.Fperiod = 60.0 * 60.0 * 24.0 * 8.0
		graphd.Fint = 60.0 * 60.0 * 24.0 * 1.0
		// graphd.Div = 2
	case "4 days":
		graphd.Fperiod = 60.0 * 60.0 * 24.0 * 4.0
		graphd.Fint = 60.0 * 60.0 * 12.0
		// graphd.Div = 2
	case "2 days":
		graphd.Fperiod = 60.0 * 60.0 * 24.0 * 2.0
		graphd.Fint = 60.0 * 60.0 * 6.0
		// graphd.Div = 2
	case "1 day":
		graphd.Fperiod = 60.0 * 60.0 * 24.0
		graphd.Fint = 60.0 * 60.0 * 3.0
		// graphd.Div = 2
	case "12 hours":
		graphd.Fperiod = 60.0 * 60.0 * 12.0
		graphd.Fint = 60.0 * 60.0 * 1.0
		// graphd.Div = 3
	case "6 hours":
		graphd.Fperiod = 60.0 * 60.0 * 6.0
		graphd.Fint = 60.0 * 30.0
		// graphd.Div = 4
	case "3 hours":
		graphd.Fperiod = 60.0 * 60.0 * 3.0
		graphd.Fint = 60.0 * 15.0
		graphd.Div = 4
	case "2 hours":
		graphd.Fperiod = 60.0 * 60.0 * 2.0
		graphd.Fint = 60.0 * 15.0
		// graphd.Div = 2
	case "1 hour":
		graphd.Fperiod = 60.0 * 60.0
		graphd.Fint = 60.0 * 10.0
		// graphd.Div = 3
	case "30 minutes":
		graphd.Fperiod = 60.0 * 20.0
		graphd.Fint = 60.0 * 5.0
		// graphd.Div = 3
	case "20 minutes":
		graphd.Fperiod = 60.0 * 15.0
		graphd.Fint = 60.0 * 2.0
		// graphd.Div = 5
	case "10 minutes":
		graphd.Fperiod = 60.0 * 10.0
		graphd.Fint = 60.0
		graphd.Div = 5
	case "4 minutes":
		graphd.Fperiod = 60.0 * 4.0
		graphd.Fint = 60.0
		graphd.Div = 2
	default:
		graphd.Period = "2 days"
		graphd.Fperiod = 60.0 * 60.0 * 24.0 * 2.0
		graphd.Fint = 60.0 * 60.0 * 6.0
		graphd.Div = 2
	}
	return

}
