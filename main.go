// Copyright © 2024 chouette2100@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php
package main

//
//

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-gorp/gorp"

	"github.com/astaxie/session"
	_ "github.com/astaxie/session/providers/memory"
)

/*
v1.0.1	Make the link description in English.
v1.0.2	Enable SSL in ServerConfig.yml.
v1.0.3	Create and update license.
v1.1.0	Improve the drawing method of Y-axis.
		Allows saving and loading of configuration data.
*/

const version="v010100"
// AHT10 measurement results
type Aht10 struct {
	Device      int
	Ts          time.Time
	Temperature float64
	Humidity    float64
	Status      int
}

// SCD41 measurement results
type Scd41 struct {
	Device      int
	Ts          time.Time
	Co2         int
	Temperature float64
	Humidity    float64
	Status      int
}

/*
v1.0.1	Make the link description in English.
v1.0.2	Enable SSL in ServerConfig.yml.
v1.0.3	Create and update license.
v1.1.0	Improve the drawing method of Y-axis.
		Allows saving and loading of configuration data.
v1.1.1	Fix a bug that humidity upper and lower limits in the configuration file (YmlFiles/*.yml)
		were not reflected in the graph.
*/

var Chimgfn chan int
var Chcfgfn chan int

func main() {

	logfilename := "plotm_" + version + "_" + time.Now().Format("20060102") + ".txt"
	logfile, err := os.OpenFile(logfilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic("cannnot open logfile: " + logfilename + err.Error())
	}
	defer logfile.Close()
	log.SetOutput(logfile)
	// log.SetOutput(io.MultiWriter(logfile, os.Stdout))
	// log.SetOutput(os.Stdout)

	LoadDenyIp("DenyIp.txt")

	//  =============================================

	// Open a connection to the database
	// Refer to the configuration file(DBConfig.yml for the database connection details)
	dbconfig, err := OpenDb("DBConfig.yml")
	if err != nil {
		log.Printf("Database error. err=%s.\n", err.Error())
		return
	}
	if dbconfig.UseSSH {
		defer Dialer.Close()
	}
	defer Db.Close()
	log.Printf("dbconfig=%+v.\n", dbconfig)

	dial := gorp.MySQLDialect{Engine: "InnoDB", Encoding: "utf8mb4"}
	Dbmap = &gorp.DbMap{Db: Db, Dialect: dial, ExpandSliceArgs: true}

	//  =============================================

	Dbmap.AddTableWithName(Aht10{}, "aht10").SetKeys(false, "Device", "Ts")
	Dbmap.AddTableWithName(Scd41{}, "scd41").SetKeys(false, "Device", "Ts")
	Dbmap.AddTableWithName(Accesslog{}, "accesslog").SetKeys(false, "Ts", "Eventid")

	//  =============================================
	// Server Configuration
	// Refer to the configuration file(ServerConfig.yml)
	svconfig := ServerConfig{}
	Serverconfig := &svconfig
	err = LoadConfig("ServerConfig.yml", Serverconfig)
	if err != nil {
		log.Printf("err=%s.\n", err.Error())
		os.Exit(1)
	}
	log.Printf("%+v\n", svconfig)

	rootPath := os.Getenv("SCRIPT_NAME")

	//  =============================================
	http.Handle("/", http.FileServer(http.Dir("public")))
	http.HandleFunc(rootPath+"/Measurements", HandlerMeasurements)
	http.HandleFunc(rootPath+"/Count", HandlerCount)
	//  =============================================
	go func() {
		no := 0
		Chimgfn = make(chan int)
		for {
			Chimgfn <- no
			no++
			if no >= 1000 {
				no = 0
			}
		}
	}()
	//  =============================================
	go func() {
		files, _ := os.ReadDir("YmlFiles")
		user := "User_00000.yml"
		for _, f := range files {
			fname := f.Name()
			if fname[0:5] == "User_" && fname > user {
				user = fname
			}
		}
		cfgfn, _ := strconv.Atoi(user[5:10])
		cfgfn++

		Chcfgfn = make(chan int)
		for {
			Chcfgfn <- cfgfn
			cfgfn++
		}
	}()
	//  =============================================
	// Start the Web server
	if svconfig.SSLcrt != "" {
		// Use SSL if you have a server certificate.
		log.Printf("           http.ListenAndServeTLS()\n")
		err := http.ListenAndServeTLS(":"+svconfig.HTTPport, svconfig.SSLcrt, svconfig.SSLkey, nil)
		if err != nil {
			log.Printf("%s\n", err.Error())
		}
	} else {
		// Start the web server locally (without SSL)
		log.Printf("           http.ListenAndServe()\n")
		err := http.ListenAndServe(":"+svconfig.HTTPport, nil)
		if err != nil {
			log.Printf("%s\n", err.Error())
		}
	}

}

var globalSessions *session.Manager

func init() {
	globalSessions, _ = session.NewManager("memory", "gosessionid", 60*60*24*7)
	go globalSessions.GC()
}
