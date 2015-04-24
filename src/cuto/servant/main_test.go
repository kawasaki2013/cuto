package main

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"cuto/servant/config"
	"cuto/testutil"
)

func getTestDataDir() string {
	const s = os.PathSeparator
	return fmt.Sprintf("%s%c%s%c%s%c%s%c%s",
		os.Getenv("GOPATH"), s, "test", s, "cuto", s, "servant", s, "main")
}

func TestRealMain_バージョン確認ができる(t *testing.T) {
	c := testutil.NewStdoutCapturer()
	args := new(arguments)
	args.v = true

	c.Start()
	realMain(args)
	out := c.Stop()

	if !strings.Contains(out, Version) {
		t.Error("バージョンが出力されていない。")
	}
}

func TestRealMain_設定ファイルから設定がロードされた上で内容にエラーがあればリターンコードrc_errorを返す(t *testing.T) {
	const s = os.PathSeparator
	config.FilePath = fmt.Sprintf("%s%c%s",
		getTestDataDir(), s, "error.ini")

	args := new(arguments)
	rc := realMain(args)

	if config.Servant.Sys.BindPort != 65536 {
		t.Error("取得した設定値が想定と違っている。")
	}
	if rc != rc_error {
		t.Errorf("リターンコード[%d]が想定値と違っている。", rc)
	}
}

func TestRealMain_ロガー初期化でのエラー発生時にリターンコードrc_errorを返す(t *testing.T) {
	const s = os.PathSeparator
	config.FilePath = fmt.Sprintf("%s%c%s",
		getTestDataDir(), s, "logerror.ini")

	args := new(arguments)
	rc := realMain(args)

	if rc != rc_error {
		t.Errorf("リターンコード[%d]が想定値と違っている。", rc)
	}
}

func TestRealMain_Run関数でのエラー発生時にリターンコードrc_errorを返す(t *testing.T) {
	const s = os.PathSeparator
	config.FilePath = fmt.Sprintf("%s%c%s",
		getTestDataDir(), s, "binderror.ini")

	args := new(arguments)
	rc := realMain(args)

	if rc != rc_error {
		t.Errorf("リターンコード[%d]が想定値と違っている。", rc)
	}
}

func TestFetchArgs_実行時引数を取得できる(t *testing.T) {
	os.Args = append(os.Args, "-v")
	args := fetchArgs()
	if !args.v {
		t.Error("バージョン出力オプションが取得できていない。")
	}
}
