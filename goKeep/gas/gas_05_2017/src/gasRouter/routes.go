package gasRouter

import (
	"encoding/json"
	"fmt"
	gLogic "gasLogic"
	"io/ioutil"
	"net/http"
	"strings"
)

// Load from json and save to gob to generate correct svg
func (app *App) LoadGasFromJsonFile(w http.ResponseWriter, r *http.Request) {

	if err := gc.LoadFromJsonFile(app.JsonFile); err != nil {
		http.Error(w, "Could not load gas from json file", http.StatusNotFound)
		return
	}
	gc.SaveToGobFile(app.GobFile)
	gc.GenerateSvgFromSnaps(gc.LastThreeYearsSlice())
	w.Write([]byte("ok"))
}

func (app *App) GetGasSliceFromGobFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	gasSlice, err := gLogic.LoadFromGobFile(app.GobFile)
	if err != nil {
		w.Write([]byte("Could not load from GobFile: " + err.Error()))
	}

	if err := json.NewEncoder(w).Encode(gasSlice); err != nil {
		http.Error(w, "Could not encode to json", http.StatusForbidden)
		return
	}

}

func (app *App) GetGasSliceFromJsonFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	gasSlice, err := gLogic.GetGasSliceFromJsonFile(gLogic.GasData1JsonPath)
	if err != nil {
		http.Error(w, "Could not load json file", http.StatusForbidden)
		return
	}

	if err := json.NewEncoder(w).Encode(gasSlice); err != nil {
		http.Error(w, "Could not encode to json", http.StatusForbidden)
		return
	}
}

func (app *App) GetLast_cPrice(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	cPrice := fmt.Sprintf("%v", gc.Snaps[len(gc.Snaps)-1].CPrice)
	w.Write([]byte(cPrice))
}

func (app *App) InsertNewGasSnapFormString(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	data, _ := ioutil.ReadAll(r.Body)

	dataSlice := strings.Split(string(data), "|")
	if len(dataSlice) != 3 {
		http.Error(w, fmt.Sprintf("Invalid nr of args %v", dataSlice), http.StatusForbidden)
		return
	}

	err := gc.Idate(dataSlice[0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	err = gc.IcMetersS(dataSlice[1])
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	err = gc.IcPrice(dataSlice[2])
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	err = gc.IaddGasSnap()
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	//w.Write([]byte("ok"))
}

func (app *App) InsertNewGasSnapForm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	err := r.ParseForm()
	if err != nil {
		w.Write([]byte("Could not parse Form"))
	}

	formDate := r.PostFormValue("date") // "2017/04/15 09:26"
	form_cMeters := r.PostFormValue("cMeters")
	form_cPrice := r.PostFormValue("cPrice")

	err = gc.IcMetersS(form_cMeters)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	err = gc.Idate(formDate)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	err = gc.IcPrice(form_cPrice)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	if err == nil {
		err = gc.IaddGasSnap()
	}

	w.Write([]byte("ok"))

}

func (app *App) UpdateGas(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "err UpdateGas: could not read req body", http.StatusForbidden)
		return
	}

	dataSlice := strings.Split(string(data), " ")
	if len(dataSlice) != 2 {
		http.Error(w, fmt.Sprintf("Invalid data: %v", dataSlice), http.StatusForbidden)
		return
	}

	valToBeUpdated := dataSlice[0]
	newVal := dataSlice[1]

	err = gc.Update(valToBeUpdated, newVal)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
	}

}

// Generate the svg from saved file
func (app *App) GenSvgFromJsonFile(w http.ResponseWriter, r *http.Request) {

	slc, err := gLogic.GetGasSliceFromJsonFile(app.JsonFile)
	if err != nil {
		http.Error(w, "Could not load slices from json file", http.StatusForbidden)
		return
	}

	gLogic.GenerateSvgFromSnaps(gLogic.LastThreeYearsSlice(slc))

	svg, err := ioutil.ReadFile(app.GraphPath)
	if err != nil {
		w.Write([]byte("Error reading graphPath"))
		return
	}

	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(svg)
}

func (app *App) GetSvgFromMemory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Cache-Control", "max-age=0, must-revalidate, no-store")

	gc.GenerateSvgFromSnapsInMemory(gc.LastThreeYearsSlice(), w)
}

func (app *App) DeleteLastGas(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	err := gc.DeleteLast()
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Could not Delete last entry err:" + err.Error()))
		return
	}

	w.Write([]byte("ok"))
}

func (app *App) SaveToFiles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	err := gc.SaveToGobFile(app.GobFile)
	err = gc.SaveGasSliceToJsonFile(app.JsonFile)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Could not save to file gob or json" + err.Error()))
	}

	gc.GenerateSvgFromSnaps(gc.LastThreeYearsSlice())
	w.Write([]byte("ok"))
}

func (app *App) InititalWriteToGobFile(w http.ResponseWriter, r *http.Request) {
	err := gLogic.InitialWriteToGobFile(app.GobFile)
	if err != nil {
		w.Write([]byte("InitialWrite to gob file failed err: " + err.Error()))
	}

	gc.GenerateSvgFromSnaps(gc.LastThreeYearsSlice())
}
