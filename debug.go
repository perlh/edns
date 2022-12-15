package main

import (
	"log"
	"runtime"
)

func ednsConsole(message string) {
	// 获取当前函数信息
	pc, _, _, _ := runtime.Caller(0)
	// 获取当前函数名
	funcName := runtime.FuncForPC(pc).Name()
	// 打印当前函数名
	//fmt.Println(funcName)
	log.Println(funcName, message)
}
