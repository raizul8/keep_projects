package main

import (
	"fmt"
	"gas_component"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"

	"os"

	"github.com/gorilla/mux"
)

var _ = gas_component.GasSct{}

//var _ = time.ANSIC
var _ = ioutil.Discard
var _ = reflect.Array

var gComp = gas_component.GasComponent{}

type MyServer struct {
	r *mux.Router
}

func (s *MyServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if origin := req.Header.Get("Origin"); origin != "" {
		rw.Header().Set("Access-Control-Allow-Origin", origin)
		rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		rw.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	}
	// Stop here if its Preflighted OPTIONS request
	if req.Method == "OPTIONS" {
		return
	}
	// Lets Gorilla work
	s.r.ServeHTTP(rw, req)
}

func SaveGas(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain;charset=utf-8")
	w.Write([]byte(gComp.SaveDataToFile()))
}

func GetGas(w http.ResponseWriter, r *http.Request) {
	gas_data, err := ioutil.ReadFile("gas_data1.json")
	if err != nil {
		log.Println("Err at GetGas gas.json:", err)
		os.Exit(1)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(gas_data)
}

func LoadFromMemoryGas(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(gComp.LoadFromMemoryGas())
}

func IDate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}

	b, _ := ioutil.ReadAll(r.Body)
	w.Header().Set("Content-Type", "text/plain;charset=utf-8")
	w.Write([]byte(gComp.IDate(string(b))))
}

func ICmeters(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}

	b, _ := ioutil.ReadAll(r.Body)
	w.Header().Set("Content-Type", "text/plain;charset=utf-8")
	w.Write([]byte(gComp.ICmeters(string(b))))
}

func ICprice(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}

	b, _ := ioutil.ReadAll(r.Body)
	w.Header().Set("Content-Type", "text/plain;charset=utf-8")
	w.Write([]byte(gComp.ICprice(string(b))))
}

func IAddGasSnap(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain;charset=utf-8")
	w.Write([]byte(gComp.IaddGasSnap()))
}

func Select(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}
	b, _ := ioutil.ReadAll(r.Body)
	w.Header().Set("Content-Type", "text/plain;charset=utf-8")
	w.Write([]byte(gComp.Select(string(b))))
}

func DeleteLast(w http.ResponseWriter, r *http.Request) {
	gComp.DeleteLast()
}

func UcPrice(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}
	b, _ := ioutil.ReadAll(r.Body)
	w.Header().Set("Content-Type", "text/plain;charset=utf-8")
	w.Write([]byte(gComp.UcPrice(string(b))))
}

func UcMeters(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}
	b, _ := ioutil.ReadAll(r.Body)
	w.Header().Set("Content-Type", "text/plain;charset=utf-8")
	w.Write([]byte(gComp.UcMeters(string(b))))
}

func UcalcSelectedSnap(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain;charset=utf-8")
	w.Write([]byte(gComp.UCalcSelectedSnap()))
}

func main() {
	gComp.LoadJsonGas()
	gComp.SelectedSnapPos = -1

	//	for _, v := range gComp.Snaps[18:] {
	//		fmt.Printf("%v\n", v)
	//	}
	fmt.Println("Server Start :8080")
	r := mux.NewRouter()

	r.HandleFunc("/loadFromMemoryGas", LoadFromMemoryGas)
	r.HandleFunc("/getGas", GetGas)
	r.HandleFunc("/iDate", IDate)
	r.HandleFunc("/iCmeters", ICmeters)
	r.HandleFunc("/iCprice", ICprice)
	r.HandleFunc("/iAddGasSnap", IAddGasSnap)
	r.HandleFunc("/select", Select)
	r.HandleFunc("/deleteLast", DeleteLast)
	r.HandleFunc("/uCprice", UcPrice)
	r.HandleFunc("/uCmeters", UcMeters)
	r.HandleFunc("/uCalcSelectedSnap", UcalcSelectedSnap)
	r.HandleFunc("/saveGas", SaveGas)
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("../web/web"))))

	http.Handle("/", &MyServer{r})
	err := http.ListenAndServe(":8080", nil)
	log.Println("err at main:", err)
	fmt.Println("Server stoped")
}
