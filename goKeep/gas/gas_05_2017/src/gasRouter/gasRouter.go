package gasRouter

import (
	"net/http"

	"raz_utils"

	"github.com/gorilla/mux"

	gl "gasLogic"
)

type App struct {
	GraphPath string
	GobFile   string
	JsonFile  string
}

var r *mux.Router = mux.NewRouter()
var gc = new(gl.GasComponent)

//var AllSvgSections map[int][]p1.GasSct = make(map[int][]p1.GasSct)
//
//func (app *App) generateSvgFromSnaps(last3yearsSnaps [][]p1.GasSct) {
//	f, err := os.Create(app.GraphPath)
//	raz_utils.LogErrorFatal(err)
//
//	width := 1250
//	height := 640
//	canvas := svg.New(f)
//	canvas.Start(width, height)
//
//	canvas.Rect(0, 0, width, height, "fill='lightgray'")
//
//	graphYearHeight := 200
//	graphYearLen := 100
//
//	lineStyleGas := "fill:red;opacity:0.5"
//	lineStyleMoney := "fill:green;opacity:0.5"
//
//	// for each year draw graph
//	for j, yearSlice := range last3yearsSnaps {
//		baseY := graphYearHeight * (j + 1)
//		var baseX int
//
//		spathMoney := fmt.Sprintf("M%v %v ", graphYearLen, baseY)
//		spathGas := fmt.Sprintf("M%v %v", graphYearLen, baseY)
//		// Draw years for each slice
//		canvas.Text(graphYearLen-90, baseY-50, fmt.Sprintf("%v", yearSlice[0].Date.Year()), "font-size:35px;")
//
//		// For each month calc position
//		for i, g := range yearSlice {
//			gasY := baseY - int(0.14*float64(g.Consumed))
//			moneY := baseY - int(0.132*g.Total)
//
//			baseX = graphYearLen * (i + 1)
//			spathMoney += fmt.Sprintf("L%v %v ", baseX, moneY)
//			spathGas += fmt.Sprintf("L%v %v ", baseX, gasY)
//			// Draw month name
//			canvas.Text(baseX, baseY+20, g.Date.Month().String(), "text-anchor:middle")
//
//		}
//		spathMoney += fmt.Sprintf("L%v %v Z", baseX, baseY)
//		spathGas += fmt.Sprintf("L%v %v Z", baseX, baseY)
//		//fmt.Println("spath: ", spath)
//		canvas.Path(spathMoney, lineStyleMoney)
//		canvas.Path(spathGas, lineStyleGas)
//		//fmt.Printf("yearSlice: %v \n", yearSlice)
//	}
//
//	canvas.End()
//	f.Close()
//}
//
//func (app *App) lastThreeYearsSlice(snaps []p1.GasSct) [][]p1.GasSct {
//	AllSvgSections = map[int][]p1.GasSct{}
//	for _, gas := range snaps {
//		AllSvgSections[gas.Date.Year()] = append(AllSvgSections[gas.Date.Year()], gas)
//	}
//
//	year := time.Now().Year()
//	var lastYearInSections int
//
//	// check if current year exists in AllSvgSection
//	if _, ok := AllSvgSections[year]; ok {
//		lastYearInSections = year
//	} else {
//		lastYearInSections = year - 1
//	}
//
//	// Show only last 3 years
//	last3yearsSlice := [][]p1.GasSct{}
//	for i := 0; i < 3; i++ {
//		last3yearsSlice = append(last3yearsSlice, AllSvgSections[lastYearInSections-i])
//
//	}
//
//	return last3yearsSlice
//}

func (app *App) NewRouter() *mux.Router {
	//if err := gc.LoadFromJsonFile(app.JsonFile); err != nil {
	//    raz_utils.LogErrorFatal(err)
	//}

	if err := gc.LoadFromGobFile(app.GobFile); err != nil {
		raz_utils.LogErrorFatal(err)
	}

	r.Methods("OPTIONS").HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if origin := r.Header.Get("Origin"); origin != "" {
				//fmt.Println("Origin:", origin)
				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
				w.Header().Set("Access-Control-Allow-Headers",
					"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
			}
		})

	r.HandleFunc("/LoadGasFromJsonFile", app.LoadGasFromJsonFile).Methods("GET")
	r.HandleFunc("/GetGasSliceFromJsonFile", app.GetGasSliceFromJsonFile).Methods("GET")
	r.HandleFunc("/GetGasSliceFromGobFile", app.GetGasSliceFromGobFile).Methods("GET")
	r.HandleFunc("/GetGasSliceFromMemory", app.GetSvgFromMemory).Methods("GET")

	//r.HandleFunc("/InsertNewGasSnap", app.InsertNewGasSnap).Methods("POST")
	r.HandleFunc("/InsertNewGasSnapForm", app.InsertNewGasSnapForm).Methods("POST")
	r.HandleFunc("/InsertNewGasSnapFormString", app.InsertNewGasSnapFormString).Methods("POST")
	r.HandleFunc("/DeleteLastGas", app.DeleteLastGas).Methods("GET")
	r.HandleFunc("/SaveToFiles", app.SaveToFiles).Methods("GET")
	r.HandleFunc("/InititalWriteToGobFile", app.InititalWriteToGobFile).Methods("GET")
	r.HandleFunc("/GetSvg", app.GenSvgFromJsonFile).Methods("GET")
	r.HandleFunc("/GetSvgFromMemory", app.GetSvgFromMemory).Methods("GET")
	r.HandleFunc("/GetLast_cPrice", app.GetLast_cPrice).Methods("GET")
	r.HandleFunc("/UpdateGas", app.UpdateGas).Methods("POST")

	return r
}

//func (app *App) InsertNewGasSnap(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("Content-Type", "text/plain;charset=utf-8")
//	body, err := ioutil.ReadAll(r.Body)
//
//	var auxGasSnap struct {
//		DateS   string
//		CMeters string
//		CPrice  string
//	}
//
//	err = json.Unmarshal(body, &auxGasSnap)
//	if err != nil {
//		w.Write([]byte("Could not unmarshal input gasSnap"))
//		return
//	}
//}
