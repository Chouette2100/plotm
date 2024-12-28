package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type DenyIp struct {
	Bip int
	Eip int
}

var DenyIpList []DenyIp

func LoadDenyIp(filename string) (
	err error,
) {
	DenyIpList = make([]DenyIp, 0)

	f, err := os.OpenFile(filename, os.O_RDONLY, 0666)
	if err != nil {
		err = fmt.Errorf("LoadDenyIp: %w", err)
		return err
	}

	var line string
	for {
		_, err = fmt.Fscanf(f, "%s%\n", &line)
		if err != nil {
			if err.Error() != "unexpected newline" {
				break
			}
		}
		linea := strings.Split(line, "/")

		atoi := func(s string) int {
			i, _ := strconv.Atoi(s)
			return i
		}

		sipa := strings.Split(linea[0], ".")

		bip := 0
		for i := 0; i < 4; i++ {
			bip = bip*256 + atoi(sipa[i])
		}

		m := 1
		for i := 0; i < 32-atoi(linea[1]); i++ {
			m = m * 2
		}

		DenyIpList = append(DenyIpList, DenyIp{bip, bip + m - 1})

	}
	log.Printf("DenyIpList: %v\n", DenyIpList)
	return nil

}

func IsAllowIp(sip string) bool {

	atoi := func(s string) int {
		i, _ := strconv.Atoi(s)
		return i
	}

	//	change IP address to integer
	sipa := strings.Split(sip, ".")
	ip := 0
	for i := 0; i < 4; i++ {
		ip = ip*256 + atoi(sipa[i])
	}

	for _, v := range DenyIpList {
		if ip >= v.Bip && ip <= v.Eip {
			return false
		}
	}

	return true
}
