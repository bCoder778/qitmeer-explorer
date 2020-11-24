package main

import (
	"flag"
	"fmt"
	"github.com/bCoder778/log"
	"github.com/bCoder778/qitmeer-explorer/api"
	"github.com/bCoder778/qitmeer-explorer/conf"
	"github.com/bCoder778/qitmeer-explorer/version"
	"os"
	"runtime"
	"runtime/debug"
)

func main() {
	setSystemResource()
	dealCommand()
	runApi()
}

func dealCommand() {
	v := flag.Bool("v", false, "show bin info")
	flag.Parse()

	if *v {
		_, _ = fmt.Fprint(os.Stderr, version.StringifyMultiLine())
		os.Exit(1)
	}
}

func setSystemResource() {
	cpuNumber := runtime.NumCPU()
	gcPercent := 20
	if conf.Setting != nil {
		if conf.Setting.Resources.CPUNumber < cpuNumber {
			cpuNumber = conf.Setting.Resources.CPUNumber
		}
		if conf.Setting.Resources.GCPercent > 0 && conf.Setting.Resources.GCPercent < 100 {
			gcPercent = conf.Setting.Resources.GCPercent
		}

	}
	runtime.GOMAXPROCS(runtime.NumCPU())
	debug.SetGCPercent(gcPercent)
}

func runApi() {
	log.SetOption(&log.Option{
		LogLevel: conf.Setting.Log.Level,
		Mode:     conf.Setting.Log.Mode,
		Email:    &log.EMailOption{},
	})

	a, err := api.NewApi(conf.Setting)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err := a.Run(); err != nil {
		fmt.Println(err)
	}
}
