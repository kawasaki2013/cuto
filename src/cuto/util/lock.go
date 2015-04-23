// Copyright 2015 unirita Inc.
// Created 2015/04/10 shanxia

package util

import (
	"fmt"
	"syscall"
	"unsafe"
)

// ミューテックスハンドルを保持する。
type MutexHandle struct {
	handle uintptr
}

const (
	wAIT_OBJECT_0  int = 0
	wAIT_ABANDONED int = 128
	wAIT_TIMEOUT   int = 258
)

// プロセス間で共通に使用する名前を指定する。
func InitMutex(name string) (*MutexHandle, error) {
	mutexName := fmt.Sprintf("Global\\%s", name)
	hMutex, _, err := procCreateMutexW.Call(
		0, 0, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(mutexName))))
	if hMutex == 0 {
		return nil, err
	}
	return &MutexHandle{hMutex}, nil
}

// ロックを開始する。
// 引数でタイムアウト時間（ミリ秒）を指定する。
func (m *MutexHandle) Lock(timeout_milisec int) (bool, error) {
	var success bool
	r1, _, _ := procWaitForSingleObject.Call(m.handle, uintptr(timeout_milisec))
	if int(r1) == wAIT_OBJECT_0 || int(r1) == wAIT_ABANDONED {
		success = true
	} else if int(r1) == wAIT_TIMEOUT {
		success = false
	}
	return success, nil
}

// ロック中であれば、解除する。
func (m *MutexHandle) Unlock() error {
	procReleaseMutex.Call(m.handle)
	return nil
}

// InitMutexで確保したミューテックスオブジェクトを破棄する。
func (m *MutexHandle) TermMutex() error {
	procCloseHandle.Call(m.handle)
	return nil
}