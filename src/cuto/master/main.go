// Copyright 2015 unirita Inc.
// Created 2015/04/10 honda

package main

import (
	"flag"
	"fmt"
	"os"

	"cuto/console"
	"cuto/log"

	"cuto/master/config"
	"cuto/master/jobnet"
)

// 実行時引数のオプション
type arguments struct {
	versionFlag bool   // バージョン情報表示フラグ
	networkName string // ジョブネットワーク名
	startFlag   bool   // 実行フラグ
	configPath  string // 設定ファイルのパス
}

// masterの戻り値
const (
	rc_OK    = 0
	rc_ERROR = 1
)

// フラグ系実行時引数のON/OFF
const (
	flag_ON  = true
	flag_OFF = false
)

const defaultConfig = `master.ini`

func main() {
	console.Display("CTM001I", os.Getpid())

	args := fetchArgs()
	rc := realMain(args)
	console.Display("CTM002I", rc)
	os.Exit(rc)
}

func realMain(args *arguments) int {
	if args.versionFlag == flag_ON {
		showVersion()
		return rc_OK
	}

	if args.networkName == "" {
		showUsage()
		return rc_ERROR
	}

	if args.configPath == "" {
		args.configPath = defaultConfig
	}

	if err := config.Load(args.configPath); err != nil {
		console.Display("CTM004E", args.configPath)
		return rc_ERROR
	}

	if err := config.DetectError(); err != nil {
		console.Display("CTM005E", err)
		return rc_ERROR
	}

	if err := log.Init(config.Dir.LogDir,
		"master",
		config.Log.OutputLevel,
		config.Log.MaxSizeKB,
		config.Log.MaxGeneration); err != nil {
		console.Display("CTM021E", err)
		return rc_ERROR
	}
	defer log.Term()

	nwk := jobnet.LoadNetwork(args.networkName)
	if nwk == nil {
		return rc_ERROR
	}
	defer nwk.Terminate()

	err := nwk.DetectFlowError()
	if err != nil {
		console.Display("CTM011E", nwk.MasterPath, err)
		return rc_ERROR
	}

	if args.startFlag == flag_OFF {
		console.Display("CTM020I", nwk.MasterPath)
		return rc_OK
	}

	err = nwk.LoadJobEx()
	if err != nil {
		console.Display("CTM004E", nwk.JobExPath)
		log.Error(err)
		return rc_ERROR
	}

	// @Todo: ネットワークへの仮ID付与を修正
	//	nwk.ID = 12345678

	err = nwk.Run()
	if err != nil {
		console.Display("CTM013I", nwk.Name, nwk.ID, "ABNORMAL")
		log.Error(err)
		return rc_ERROR
	}
	console.Display("CTM013I", nwk.Name, nwk.ID, "NORMAL")
	return rc_OK
}

// コマンドライン引数を解析し、arguments構造体を返す。
func fetchArgs() *arguments {
	args := new(arguments)
	flag.Usage = showUsage
	flag.BoolVar(&args.versionFlag, "v", false, "version option")
	flag.StringVar(&args.networkName, "n", "", "network name option")
	flag.BoolVar(&args.startFlag, "s", false, "start option")
	flag.StringVar(&args.configPath, "c", "", "config file option")
	flag.Parse()
	return args
}

// バージョンを表示する。
func showVersion() {
	fmt.Printf("cuto master version %s\n", Version)
}

// オンラインヘルプを表示する。
func showUsage() {
	console.Display("CTM003E")
	fmt.Print(console.USAGE)
}