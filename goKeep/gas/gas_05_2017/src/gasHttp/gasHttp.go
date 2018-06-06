package main

import (
	"fmt"
	"gasRouter"

	gl "gasLogic"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"raz_utils"
	"strconv"
	"time"
)

var Graph1Path string
var GobWriteFile string

//var GobReadFile string
var JsonWriteFile string

// Sets debug mode from os.Args, if not set program is in debug mode.
func setDebugMode() {
	var debugMode bool = true

	if len(os.Args) >= 2 {
		debugM, err := strconv.ParseBool(os.Args[1])
		if err != nil {
			raz_utils.LogErrorFatal(err)
		}
		debugMode = debugM

		fmt.Printf("debugMode type: %T value: %t \n", debugMode, debugMode)
		gl.DebugMode = debugM
		return
	}

	gl.DebugMode = debugMode
	fmt.Printf("debugMode type: %T value: %t \n", debugMode, debugMode)

}

var gasCurrentDir string

func init() {
	// Set debugMode var to true or false
	setDebugMode()
	//checkPaths()

	if !gl.DebugMode {
		path, err := raz_utils.GetDirToBinary()
		raz_utils.LogErrorFatal(err)

		gasCurrentDir = path
	} else {
		path, err := raz_utils.GetRunningProgramDir()
		raz_utils.LogErrorFatal(err)

		gasCurrentDir = path
	}

	gl.InitializePaths(gasCurrentDir)

	Graph1Path = filepath.Join(gasCurrentDir, "data/graph1.svg")
	GobWriteFile = filepath.Join(gasCurrentDir, "data/gobFile")
	//GobReadFile = filepath.Join(gasCurrentDir, "data/gobFile")
	JsonWriteFile = filepath.Join(gasCurrentDir, "data/gas_data1.json")

}

func main() {
	address := "127.0.0.1:8000"
	fmt.Println(address)
	fmt.Println("Stop debug mode: start program with first os.Args false, ex: gasHttp false")
	dirStaticFiles := "data"

	app := new(gasRouter.App)

	app.GobFile = GobWriteFile
	app.GraphPath = Graph1Path
	app.JsonFile = JsonWriteFile

	r := app.NewRouter()
	// http://localhost:8000/static/<filename>
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(dirStaticFiles))))

	srv := &http.Server{
		Handler:      r,
		Addr:         address,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
