// Copyright © 2024 chouette2100@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php
package main
//
//
/*
v1.0.1	Make the link description in English.
v1.0.2	Enable SSL in ServerConfig.yml.
v1.0.3	Create and update license.
*/

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-gorp/gorp"
)

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

var Ch chan int

func main() {
	LoadDenyIp("DenyIp.txt")
	logfilename := "plotm" + time.Now().Format("20060102") + ".txt"
	logfile, err := os.OpenFile(logfilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic("cannnot open logfile: " + logfilename + err.Error())
	}
	defer logfile.Close()
	log.SetOutput(logfile)
	// log.SetOutput(io.MultiWriter(logfile, os.Stdout))

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
	//  =============================================
	go func() {
		no := 0
		Ch = make(chan int)
		for {
			Ch <- no
			no++
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
