//package main

package gasLogic

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	rz "raz_utils"
	"strconv"
	"time"

	"io"

	svg "github.com/ajstarks/svgo"
)

type MyDate struct{ time.Time }

func (d *MyDate) UnmarshalJSON(data []byte) error {
	var s string

	err := json.Unmarshal(data, &s)

	rz.LogError(err)
	if err != nil {
		return err
	}

	t1, err1 := time.Parse("2006-01-02T15:04:05.000Z", s)
	t2, err2 := time.Parse("2006-01-02T15:04:05.000000000-07:00", s)

	if err1 != nil && err2 != nil {
		return fmt.Errorf("invalid date: %v \ninvalid date: %v", err1, err2)
	}

	if err1 == nil {
		d.Time = t1
	} else {
		d.Time = t2
	}
	return nil
}

func (d *MyDate) MarshalJSON() ([]byte, error) {
	dateFormat := d.Format(`"2006-01-02T15:04:05.000Z"`)
	return []byte(dateFormat), nil
}

type GasSct struct {
	Date     time.Time `json:"date"`
	CMeters  int64     `json:"cMeters"`
	Consumed int64     `json:"Consumed"`
	CPrice   float64   `json:"cPrice"`
	Total    float64   `json:"total"`
}

func (g *GasSct) UnmarshalJSON(data []byte) error {

	var aux struct {
		Date     MyDate
		CMeters  int64 `json:"cMeters"`
		Consumed int64
		CPrice   float64 `json:"cPrice"`
		Total    float64
	}

	dec := json.NewDecoder(bytes.NewReader(data))
	if err := dec.Decode(&aux); err != nil {
		return fmt.Errorf("decode err GasSct: %v", err)
	}

	// fmt.Println("__*GasSct UnmarshalJSON ", aux.Date.Time, aux.CMeters, aux.CPrice, aux.Total)
	g.Date = aux.Date.Time
	g.CMeters = aux.CMeters
	g.Consumed = aux.Consumed
	g.CPrice = aux.CPrice
	g.Total = aux.Total

	return nil
}

func (g *GasSct) MarshalJSON() ([]byte, error) {
	myDate := MyDate{}
	myDate.Time = g.Date
	buff := &bytes.Buffer{}

	// Assign g.Date to MyDate to encode in custom format already defined in MyDate type
	err := json.NewEncoder(buff).Encode(&myDate)
	rz.LogErrorFatal(err)

	//date := buff.String()
	dateS := buff.String()
	dateS = dateS[1 : len(dateS)-2] // date comes in json string need to remove surounding "" + \n from buffer

	jsonFormat := `{"date": "%v", "cMeters": %v, "consumed": %v, "cPrice": %v, "total": %v}`
	myGasSctJson := fmt.Sprintf(jsonFormat, dateS, g.CMeters, g.Consumed, g.CPrice, g.Total)

	return []byte(myGasSctJson), nil
}

// Returns a slice of GasSct from json file gas_data1.json
func GetGasSliceFromJsonFile(gasJsonFileAddr string) ([]GasSct, error) {
	data, err := ioutil.ReadFile(gasJsonFileAddr)
	if err != nil {
		return nil, err
	}

	var gasSctSliceAux []GasSct

	dec := json.NewDecoder(bytes.NewReader(data))
	err = dec.Decode(&gasSctSliceAux)
	rz.LogErrorFatal(err)

	return gasSctSliceAux, nil
}

var GasData1JsonPath = "data/gas_data1.json"
var GobFilePath string = "data/gobFile"
var GraphPath = "data/graph1.svg"
var CurrentDir string = ""
var DebugMode bool = true

func InitializePaths(currentDir string) error {

	CurrentDir = currentDir
	GasData1JsonPath = filepath.Join(CurrentDir, GasData1JsonPath)
	GobFilePath = filepath.Join(CurrentDir, GobFilePath)
	GraphPath = filepath.Join(CurrentDir, GraphPath)

	return nil
}

// Read the json file and return it as bytes
//func ReadGasDataFile() ([]byte, error) {
//
//	data, err := ioutil.ReadFile(GAS_DATA1_JSON_LOCATION)
//	return data, err
//}

// Read from Json file and write it to gobFile.
func InitialWriteToGobFile(gobPath string) error {
	gasSlice, err := GetGasSliceFromJsonFile(GasData1JsonPath)
	if err != nil {
		return err
	}

	network := new(bytes.Buffer)
	enc := gob.NewEncoder(network)

	err = enc.Encode(&gasSlice)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(gobPath, network.Bytes(), 0664)
	if err != nil {
		return err
	}

	return nil

}

func LoadFromGobFile(gobPath string) ([]GasSct, error) {
	data, err := ioutil.ReadFile(gobPath)
	if err != nil {
		return nil, err
	}

	buffDataFromFile := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buffDataFromFile)

	var mySlice []GasSct
	err = dec.Decode(&mySlice)
	if err != nil {
		return nil, err
	}

	return mySlice, nil

}

func SaveToGobFile(gobPath string, gasSlice []GasSct) error {
	buffer := new(bytes.Buffer)
	enc := gob.NewEncoder(buffer)

	err := enc.Encode(&gasSlice)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(gobPath, buffer.Bytes(), 0664)
	if err != nil {
		return err
	}

	return err
}

// TODO-raz Remember, add check btn to UI if saved to eongas site else don't save.
type GasComponent struct {
	Snaps []GasSct
	//SelectedSnapPos int64
	ISnap GasSct
}

func (gc *GasComponent) Idate(dateS string) error {
	date, err := time.Parse("2006/01/02 15:04", dateS)
	if err != nil {
		return errors.New("invalid date")
	}

	lenSnaps := len(gc.Snaps)

	if gc.Snaps[lenSnaps-1].Date.Year() > date.Year() {
		return fmt.Errorf("Cannot insert year less than %v", gc.Snaps[lenSnaps-1].Date.Year())
	}

	if gc.Snaps[lenSnaps-1].Date.Year() == date.Year() &&
		date.Month() <= gc.Snaps[lenSnaps-1].Date.Month() {

		return fmt.Errorf("Already inserted for %v", date.Month())
	}

	gc.ISnap.Date = date
	return nil

}

func (gc *GasComponent) IcMetersInt64(cMeters int64) error {
	lastCmeters := gc.Snaps[len(gc.Snaps)-1].CMeters

	if cMeters-lastCmeters <= 0 {
		return fmt.Errorf("invalid: Please enter value bigger than %v", lastCmeters)
	}

	gc.ISnap.CMeters = cMeters
	return nil
}

func (gc *GasComponent) IcMetersS(cMeters string) error {
	if m, err := strconv.ParseInt(cMeters, 10, 64); err == nil {
		lastCmeters := gc.Snaps[len(gc.Snaps)-1].CMeters

		if m-lastCmeters <= 0 {
			return fmt.Errorf("invalid: Please enter value bigger than %v", lastCmeters)
		}
		gc.ISnap.CMeters = m
		return nil
	} else {
		return fmt.Errorf("invalid: Could not format %s to int64", cMeters)
	}
}

func (gc *GasComponent) IcPrice(price string) error {
	if v, err := strconv.ParseFloat(price, 64); err == nil {
		if v <= 0 {
			return errors.New("invalid: Input cubic price less or equal to 0 !")
		}
		gc.ISnap.CPrice = v
		return nil
	} else {
		return fmt.Errorf("invalid: Could not parse to float val: %v", price)
	}
}

//TODO-raz After IaddGasSnap reset gc.iSnap to zero values and clear form fields
func (gc *GasComponent) IaddGasSnap() error {
	consumed := gc.ISnap.CMeters - gc.Snaps[len(gc.Snaps)-1].CMeters
	if consumed <= 0 {
		return fmt.Errorf("Enter a value for cMeters bigger than %v", gc.ISnap.CMeters)
	}

	gc.ISnap.Consumed = consumed
	total := float64(consumed) * gc.ISnap.CPrice

	total, err := strconv.ParseFloat(fmt.Sprintf("%.2f", total), 64)
	if err != nil {
		return fmt.Errorf("%v", "Could not parse total")
	}

	gc.ISnap.Total = total
	gc.Snaps = append(gc.Snaps, gc.ISnap)

	return nil

}

func (gComp *GasComponent) LoadFromGobFile(gobFilePath string) error {
	sl, err := LoadFromGobFile(gobFilePath)
	if err != nil {
		return err
	}

	gComp.Snaps = sl
	return nil
}

func (gComp *GasComponent) LoadFromJsonFile(jsonFilePath string) error {
	sl, err := GetGasSliceFromJsonFile(jsonFilePath)
	if err != nil {
		return err
	}

	gComp.Snaps = sl
	return nil

}

// Returns position as int if it finds the cMeters in GasComponent.Snaps.CMeters or -1.
func (gComp *GasComponent) Select(cMetersS string) int {
	cMeters, err := strconv.ParseInt(cMetersS, 10, 64)
	if err != nil {
		return -1
	}

	var pos = -1

	for i := len(gComp.Snaps) - 1; i >= 0; i-- {
		if gComp.Snaps[i].CMeters == cMeters {
			pos = i
			break
		} else {
			pos = -1
		}
	}

	if pos == -1 {
		return -1
	}

	return pos

}

func (gc *GasComponent) Update(cMeterstoUpdateS, cMetersNewValueS string) error {

	cMetersToUpdate, err := strconv.ParseInt(cMeterstoUpdateS, 10, 64)
	if err != nil {
		return fmt.Errorf("gc.Update: %v", cMetersToUpdate)
	}

	if cMetersToUpdate <= 215 {
		return fmt.Errorf("err, val to be updated %v is less or equal than 215", cMetersToUpdate)
	}

	snapsLen := len(gc.Snaps)
	cMetersNewValue, err := strconv.ParseInt(cMetersNewValueS, 10, 64)
	if err != nil {
		return fmt.Errorf("The val %s is not an int", cMetersNewValueS)

	}
	pos := gc.Select(cMeterstoUpdateS)
	slc := []GasSct{}

	for _, v := range gc.Snaps[pos:snapsLen] {
		slc = append(slc, v)
	}

	// delete from gc.Snaps all after pos
	gc.Snaps = gc.Snaps[:pos]

	// set the new value
	slc[0].CMeters = cMetersNewValue

	// reinsert the values from slc back into gc.Snaps
	for _, v := range slc {
		//err = gc.IcMetersS(strconv.FormatInt(v.CMeters, 10))
		err = gc.IcMetersInt64(v.CMeters)
		if err != nil {
			return fmt.Errorf("Could not input cMeters err:%v", err)
		}

		gc.ISnap.CPrice = v.CPrice
		gc.ISnap.Date = v.Date

		gc.IaddGasSnap()
	}

	return err
}

func (gComp *GasComponent) SaveGasSliceToJsonFile(gasJsonFileAddr string) error {
	snapsByte, err := json.Marshal(gComp.Snaps)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(gasJsonFileAddr, snapsByte, 0664)
	if err != nil {
		return err
	}

	return nil

}

func (gComp *GasComponent) SaveToGobFile(gobPath string) error {
	buffer := new(bytes.Buffer)
	enc := gob.NewEncoder(buffer)

	err := enc.Encode(&gComp.Snaps)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(gobPath, buffer.Bytes(), 0664)
	if err != nil {
		return err
	}

	return err
}

func (gComp *GasComponent) DeleteLast() error {
	if len(gComp.Snaps) < 1 {
		return fmt.Errorf("%s", "Could not delete snaps len is zero.")
	}
	gComp.Snaps = gComp.Snaps[:len(gComp.Snaps)-1]
	return nil
}

func GenerateSvgFromSnaps(sliceOfSlicesOfSnapsPerYear [][]GasSct) {

	graphFile, err := os.Create(GraphPath)
	rz.LogErrorFatal(err)

	width := 1250
	height := 1200
	canvas := svg.New(graphFile)
	canvas.Start(width, height)

	canvas.Rect(0, 0, width, height, "fill='lightgray'")

	graphYearHeight := 250
	graphYearLen := 100

	lineStyleGas := "fill:green;opacity:0.5"
	lineStyleMoney := "fill:red;opacity:0.5"

	// for each year draw graph
	for j, yearSlice := range sliceOfSlicesOfSnapsPerYear {
		baseY := graphYearHeight * (j + 1)
		var baseX int

		spathMoney := fmt.Sprintf("M%v %v ", graphYearLen, baseY)
		spathGas := fmt.Sprintf("M%v %v", graphYearLen, baseY)
		// Draw years for each slice
		canvas.Text(graphYearLen-90, baseY-50, fmt.Sprintf("%v", yearSlice[0].Date.Year()), "font-size:35px;")

		// For each month calc position
		for i, g := range yearSlice {
			gasY := baseY - int(0.14*float64(g.Consumed))
			moneY := baseY - int(0.132*g.Total)

			baseX = graphYearLen * (i + 1)
			spathMoney += fmt.Sprintf("L%v %v ", baseX, moneY)
			spathGas += fmt.Sprintf("L%v %v ", baseX, gasY)
			// Draw month name
			month := g.Date.Month().String()
			pretLeiLuna := fmt.Sprintf("%v lei", g.Total)
			consumedGas := fmt.Sprintf("%vm\u00B3", g.Consumed)
			cMeters := fmt.Sprintf("%vm\u00B3", g.CMeters)

			canvas.Text(baseX, baseY+20, month, "text-anchor:middle")
			canvas.Text(baseX, baseY+40, pretLeiLuna, "fill:red;text-anchor:middle")
			canvas.Text(baseX, baseY+60, consumedGas, "fill:green;text-anchor:middle")
			canvas.Text(baseX, baseY+80, cMeters, "fill:green;text-anchor:middle")

		}
		spathMoney += fmt.Sprintf("L%v %v Z", baseX, baseY)
		spathGas += fmt.Sprintf("L%v %v Z", baseX, baseY)
		//fmt.Println("spath: ", spath)
		canvas.Path(spathMoney, lineStyleMoney)
		canvas.Path(spathGas, lineStyleGas)
		//fmt.Printf("yearSlice: %v \n", yearSlice)
	}

	canvas.End()
	graphFile.Close()
}

func (gc *GasComponent) GenerateSvgFromSnaps(sliceOfSlicesOfSnapsPerYear [][]GasSct) {

	graphFile, err := os.Create(GraphPath)
	rz.LogErrorFatal(err)

	width := 1250
	height := 1200
	canvas := svg.New(graphFile)
	canvas.Start(width, height)

	canvas.Rect(0, 0, width, height, "fill='lightgray'")

	graphYearHeight := 250
	graphYearLen := 100

	lineStyleGas := "fill:green;opacity:0.5"
	lineStyleMoney := "fill:red;opacity:0.5"

	// for each year draw graph
	for j, yearSlice := range sliceOfSlicesOfSnapsPerYear {
		baseY := graphYearHeight * (j + 1)
		var baseX int

		spathMoney := fmt.Sprintf("M%v %v ", graphYearLen, baseY)
		spathGas := fmt.Sprintf("M%v %v", graphYearLen, baseY)
		// Draw years for each slice
		canvas.Text(graphYearLen-90, baseY-50, fmt.Sprintf("%v", yearSlice[0].Date.Year()), "font-size:35px;")

		// For each month calc position
		for i, g := range yearSlice {
			gasY := baseY - int(0.14*float64(g.Consumed))
			moneY := baseY - int(0.132*g.Total)

			baseX = graphYearLen * (i + 1)
			spathMoney += fmt.Sprintf("L%v %v ", baseX, moneY)
			spathGas += fmt.Sprintf("L%v %v ", baseX, gasY)
			// Draw month name
			month := g.Date.Month().String()
			pretLeiLuna := fmt.Sprintf("%v lei", g.Total)
			consumedGas := fmt.Sprintf("%vm\u00B3", g.Consumed)
			cMeters := fmt.Sprintf("%vm\u00B3", g.CMeters)

			canvas.Text(baseX, baseY+20, month, "text-anchor:middle")
			canvas.Text(baseX, baseY+40, pretLeiLuna, "fill:red;text-anchor:middle")
			canvas.Text(baseX, baseY+60, consumedGas, "fill:green;text-anchor:middle")
			canvas.Text(baseX, baseY+80, cMeters, "fill:green;text-anchor:middle")

		}
		spathMoney += fmt.Sprintf("L%v %v Z", baseX, baseY)
		spathGas += fmt.Sprintf("L%v %v Z", baseX, baseY)
		//fmt.Println("spath: ", spath)
		canvas.Path(spathMoney, lineStyleMoney)
		canvas.Path(spathGas, lineStyleGas)
		//fmt.Printf("yearSlice: %v \n", yearSlice)
	}

	canvas.End()
	graphFile.Close()
}

// LastThreeYearsSlice returns a slice of slices with gas data from last three years.
func (gc *GasComponent) LastThreeYearsSlice() [][]GasSct {
	var AllSvgSections map[int][]GasSct = make(map[int][]GasSct)
	snaps := gc.Snaps

	AllSvgSections = map[int][]GasSct{}
	for _, gas := range snaps {
		AllSvgSections[gas.Date.Year()] = append(AllSvgSections[gas.Date.Year()], gas)
	}

	year := time.Now().Year()
	var lastYearInSections int

	// check if current year exists in AllSvgSection
	if _, ok := AllSvgSections[year]; ok {
		lastYearInSections = year
	} else {
		lastYearInSections = year - 1
	}

	// Show only last 3 years
	last3yearsSlice := [][]GasSct{}
	for i := 0; i < 3; i++ {
		last3yearsSlice = append(last3yearsSlice, AllSvgSections[lastYearInSections-i])

	}

	return last3yearsSlice
}

func LastThreeYearsSlice(snaps []GasSct) [][]GasSct {
	var AllSvgSections map[int][]GasSct = make(map[int][]GasSct)

	AllSvgSections = map[int][]GasSct{}
	for _, gas := range snaps {
		AllSvgSections[gas.Date.Year()] = append(AllSvgSections[gas.Date.Year()], gas)
	}

	year := time.Now().Year()
	var lastYearInSections int

	// check if current year exists in AllSvgSection
	if _, ok := AllSvgSections[year]; ok {
		lastYearInSections = year
	} else {
		lastYearInSections = year - 1
	}

	// Show only last 3 years
	last3yearsSlice := [][]GasSct{}
	for i := 0; i < 3; i++ {
		last3yearsSlice = append(last3yearsSlice, AllSvgSections[lastYearInSections-i])

	}

	return last3yearsSlice
}

func (gc *GasComponent) GenerateSvgFromSnapsInMemory(last3YearsGasData [][]GasSct, writter io.Writer) {

	//f, err := os.Create(GraphPath)
	//rz.LogErrorFatal(err)

	width := 1250
	height := 1200
	canvas := svg.New(writter)
	canvas.Start(width, height)

	canvas.Rect(0, 0, width, height, "fill='lightgray'")

	graphYearHeight := 250
	graphYearLen := 100

	lineStyleGas := "fill:green;opacity:0.5"
	lineStyleMoney := "fill:red;opacity:0.5"

	// for each year draw graph
	for j, yearSlice := range last3YearsGasData {
		baseY := graphYearHeight * (j + 1)
		var baseX int

		spathMoney := fmt.Sprintf("M%v %v ", graphYearLen, baseY)
		spathGas := fmt.Sprintf("M%v %v", graphYearLen, baseY)
		// Draw years for each slice
		canvas.Text(graphYearLen-90, baseY-50, fmt.Sprintf("%v", yearSlice[0].Date.Year()), "font-size:35px;")

		// For each month calc position
		for i, g := range yearSlice {
			gasY := baseY - int(0.14*float64(g.Consumed))
			moneY := baseY - int(0.132*g.Total)

			baseX = graphYearLen * (i + 1)
			spathMoney += fmt.Sprintf("L%v %v ", baseX, moneY)
			spathGas += fmt.Sprintf("L%v %v ", baseX, gasY)
			// Draw month name
			month := g.Date.Month().String()
			pretLeiLuna := fmt.Sprintf("%v lei", g.Total)
			consumedGas := fmt.Sprintf("%vm\u00B3", g.Consumed)
			cMeters := fmt.Sprintf("%vm\u00B3", g.CMeters)

			canvas.Text(baseX, baseY+20, month, "text-anchor:middle")
			canvas.Text(baseX, baseY+40, pretLeiLuna, "fill:red;text-anchor:middle")
			canvas.Text(baseX, baseY+60, consumedGas, "fill:green;text-anchor:middle")
			canvas.Text(baseX, baseY+80, cMeters, "fill:green;text-anchor:middle")

		}
		spathMoney += fmt.Sprintf("L%v %v Z", baseX, baseY)
		spathGas += fmt.Sprintf("L%v %v Z", baseX, baseY)
		//fmt.Println("spath: ", spath)
		canvas.Path(spathMoney, lineStyleMoney)
		canvas.Path(spathGas, lineStyleGas)
		//fmt.Printf("yearSlice: %v \n", yearSlice)
	}

	canvas.End()
	//f.Close()
}
