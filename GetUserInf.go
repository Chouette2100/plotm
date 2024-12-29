// Copyright © 2024 chouette2100@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php
package main

import (
	//	"SRCGI/ShowroomCGIlib"
	//	"bufio"
	// "bytes"
	//	"fmt"
	//	"html"
	"log"

	//	"math/rand"
	// "sort"
	// "strconv"
	"strings"
	"time"
	//	"os"

	"runtime"

	// "encoding/json"

	//	"html/template"
	"net/http"

	// "database/sql"

	"encoding/json"

	// _ "github.com/go-sql-driver/mysql"

	// "github.com/PuerkitoBio/goquery"

	//	svg "github.com/ajstarks/svgo/float"

	//	"github.com/dustin/go-humanize"

	//	"github.com/goark/sshql"
	//	"github.com/goark/sshql/mysqldrv"

	//	"github.com/Chouette2100/exsrapi"
	// "github.com/Chouette2100/srapi"
)

/*
ファンクション名とリモートアドレス、ユーザーエージェントを表示する。
*/
//	var Localhost bool
type KV struct {
	K string
	V []string
}

//      アクセスログ accesslog 2024-11-27 〜
type Accesslog struct {
        Handler       string
        Remoteaddress string
        Useragent     string
        Formvalues    string
        Eventid       string
        Roomid        int
        Ts            time.Time
}

func GetUserInf(r *http.Request) (
	ra string,
	ua string,
	isallow bool,
) {

	isallow = true

	pt, _, _, ok := runtime.Caller(1) //	スタックトレースへのポインターを得る。1は一つ上のファンクション。

	fn := ""
	if !ok {
		fn = "unknown"
	}

	fn = runtime.FuncForPC(pt).Name()
	fna := strings.Split(fn, ".")

	rap := r.RemoteAddr
	rapa := strings.Split(rap, ":")
	if rapa[0] != "[" {
		ra = rapa[0]
	} else {
		ra = "127.0.0.1"
	}
	ua = r.UserAgent()

	log.Printf("  *** %s() from %s by %s\n", fna[len(fna)-1], ra, ua)
	//	fmt.Printf("%s() from %s by %s\n", fna[len(fna)-1], ra, ua)

	if !IsAllowIp(ra) {
		log.Printf("%s is on the Blacklist(%s)", ra, ua)
		isallow = false
		return
	}

	//	パラメータを表示する
	if err := r.ParseForm(); err != nil {
		log.Printf("Error: %v\n", err)
		return
	}

	var al Accesslog
	al.Ts = time.Now().Truncate(time.Second)
	al.Handler = fna[len(fna)-1]
	al.Remoteaddress = ra
	al.Useragent = ua

	kvlist := make([]KV, len(r.Form))
	i := 0
	for kvlist[i].K, kvlist[i].V = range r.Form {
		log.Printf("%12v : %v\n", kvlist[i].K, kvlist[i].V)
		switch kvlist[i].K {
		case "fnc":
			al.Eventid = kvlist[i].V[0]
		default:
		}
		i++
	}
	jd, err := json.Marshal(kvlist)
	if err != nil {
		log.Printf(" GetUserInf(): %s\n", err.Error())
	}
	al.Formvalues = string(jd)

	err = Dbmap.Insert(&al)
	if err != nil {
		log.Printf(" GetUserInf(): %s\n", err.Error())
	}

	return
}
