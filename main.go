package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/mqu/go-notify"
	daemon "github.com/sevlyar/go-daemon"
)

// To terminate the daemon use:
//
//	kill `cat sample.pid`
func main() {
	cntx := &daemon.Context{
		PidFileName: "NotifyDaemonLowBatteryByStas228",
		PidFilePerm: 0644,
		WorkDir:     "./",
	}
	d, err := cntx.Reborn()
	if err != nil {
		panic(err)
	}
	err = d.Release()
	if err != nil {
		panic("daemon release err")
	}
	log.Print("start")
	NotificatorApp()
}
func BatteryLevel() uint8 {
	const pattern = "/sys/class/power_supply/BAT*"
	bs, err := filepath.Glob(pattern)
	if err != nil || len(bs) == 0 {
		panic("no battery")
	}
	pathBatInfo := filepath.Join(bs[0], "capacity")
	f, err := os.Open(pathBatInfo)
	if err != nil {
		panic("error while opening")
	}
	var lvl uint8
	_, err = fmt.Fscanf(f, "%d", &lvl)
	if err != nil {
		panic("error while scanning file")
	}
	return lvl
}
func NotificatorApp() {
	timeRate := 30 * time.Second
	ch := time.Tick(timeRate)
	for range ch {
		if checkBattery() {
			notifyapp()
			time.Sleep(5 * time.Minute)
		} else {
			log.Print("not now")
		}
	}
}
func checkBattery() bool {
	lvl := BatteryLevel()
	return lvl <= 15
}
func notifyapp() {
	notify.Init("LowBattery")
	alert := notify.NotificationNew("Low Battery", "Battery at low level", "dialog-information")
	alert.Show()
}
