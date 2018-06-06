package raz_utils

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

func RetErrInfo(err error) error {

	pc, fn, line, _ := runtime.Caller(1)
	return fmt.Errorf("[error] in %s[ %s:%d:0] %v", runtime.FuncForPC(pc).Name(), fn, line, err)

}

func LogErrors(err ...error) {
	// if err != nil {
	// 	// notice that we're using 1, so it will actually log where
	// 	// the error happened, 0 = this function, we don't want that.
	// 	_, fn, line, _ := runtime.Caller(1)
	// 	log.Printf("[error] %s:%d:0 %v", fn, line, err)
	// }

	for _, ve := range err {
		if err != nil {
			pc, fn, line, _ := runtime.Caller(1)
			log.Printf("[error] in %s[ %s:%d:0] %v", runtime.FuncForPC(pc).Name(), fn, line, ve)
		}
	}
}

func LogErrorsFatal(err ...error) {
	for _, ve := range err {
		if err != nil {
			pc, fn, line, _ := runtime.Caller(1)
			log.Fatalf("[error] in %s[ %s:%d:0] %v", runtime.FuncForPC(pc).Name(), fn, line, ve)
		}
	}
}

//this logs the function name as well.
func LogError(err error) {
	if err != nil {
		pc, fn, line, _ := runtime.Caller(1)
		log.Printf("[error] in %s[ %s:%d:0] %v", runtime.FuncForPC(pc).Name(), fn, line, err)
	}

}

func LogErrorFatal(err error) {
	if err != nil {
		pc, fn, line, _ := runtime.Caller(1)
		log.Fatalf("[error] in %s[ %s:%d:0] %v", runtime.FuncForPC(pc).Name(), fn, line, err)
	}

}

// Returns full path plus the name of file ex:
// /home/raz/IdeaProjects1/projects1/gas_03_2017/src/p1/p1.go
// Does not return path to binary use GetDirToBinary()
func GetRunningProgramPath() (string, error) {
	_, fn, _, ok := runtime.Caller(1)
	if ok {
		return fmt.Sprint(fn), nil
	}

	return "", errors.New("Could not get file name.")
}

// Return path to running program with the name of file stripped ex:
// /home/raz/IdeaProjects1/projects1/gas_03_2017/src/p1
// Does not return path to binary use GetDirToBinary()
func GetRunningProgramDir() (string, error) {
	_, fn, _, ok := runtime.Caller(1)
	if !ok {
		return "", errors.New("runtime.Caller it was not possible to recover the information")
	}

	return path.Dir(fn), nil

}

// Returns path to binary excluding name of file.
func GetDirToBinary() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}

	pathSlice := strings.Split(ex, "/")

	pathSlice[0] = "/"
	pathSlice = pathSlice[:len(pathSlice)-1]
	pathString := filepath.Join(pathSlice...)

	return pathString, err
}

func init() {
	// log.SetFlags(log.Ldate | log.Ltime)
	log.SetFlags(0)
	// fmt.Println("----raz_utils init -----")
}
