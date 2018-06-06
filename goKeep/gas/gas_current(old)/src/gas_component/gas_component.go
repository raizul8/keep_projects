package gas_component

import (
	"encoding/json"
	"fmt"
	//	"github.com/goinggo/mapstructure"
	"io/ioutil"
	"strconv"
	"time"
)

func check(err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
}

type GasSct struct {
	Date     time.Time `json:"date"`
	CMeters  int64     `json:"cMeters"`
	Consumed int64     `json:"consumed"`
	CPrice   float64   `json:"cPrice"`
	Total    float64   `json:"total"`
}

type GasComponent struct {
	Snaps           []GasSct
	SelectedSnapPos int64
	iSnap           GasSct
}

func (gc *GasComponent) SaveDataToFile() string {
	b, err := json.Marshal(gc.Snaps)
	check(err)

	if err = ioutil.WriteFile("gas_data1.json", b, 0644); err != nil {
		return fmt.Sprintf("Err: %v \nCould not write to file", err.Error())
	}

	return "valid: File saved ok"
}

func (gc *GasComponent) LoadFromMemoryGas() []byte {
	snaps_json, err := json.Marshal(gc.Snaps)
	if err != nil {
		return []byte(fmt.Sprintf("Error LoadFromMemoryGas: %v", err))
	}
	return snaps_json
}

func (gc *GasComponent) LoadJsonGas() {
	jsonGasData, err := ioutil.ReadFile("gas_data1.json")
	check(err)

	var gasData []map[string]interface{}

	err = json.Unmarshal(jsonGasData, &gasData)
	check(err)

	for _, gasSnap := range gasData {
		dateS, _ := gasSnap["date"].(string)
		date, _ := time.Parse("2006-01-02T15:04:05.000Z", dateS)

		gc.Snaps = append(gc.Snaps, GasSct{
			Date:     date,
			CMeters:  int64(gasSnap["cMeters"].(float64)),
			Consumed: int64(gasSnap["consumed"].(float64)),
			CPrice:   gasSnap["cPrice"].(float64),
			Total:    gasSnap["total"].(float64),
		})
	}

}

func (gc *GasComponent) IDate(dateS string) string {
	if date, err := time.Parse("2006-01-02T15:04:05.000Z", dateS); err != nil {
		return "invalid"
	} else {
		gc.iSnap.Date = date
		return "valid"
	}

}

func (gc *GasComponent) ICmeters(cMeters string) string {
	if m, err := strconv.ParseInt(cMeters, 10, 64); err == nil {
		lastCmeters := gc.Snaps[len(gc.Snaps)-1].CMeters
		if m-lastCmeters <= 0 {
			return fmt.Sprintf("invalid: Please enter value bigger than %v", lastCmeters)
		}
		gc.iSnap.CMeters = m
		return "valid"
	} else {
		return "invalid"
	}
}

func (gc *GasComponent) ICprice(price string) string {
	if v, err := strconv.ParseFloat(price, 64); err == nil {
		if v <= 0 {
			return "invalid: Input cubic price less or equal to 0 !"
		}
		gc.iSnap.CPrice = v
		return "valid"
	} else {
		return "invalid"
	}
}

func (gc *GasComponent) IaddGasSnap() string {
	consumed := gc.iSnap.CMeters - gc.Snaps[len(gc.Snaps)-1].CMeters
	if consumed == 0 {
		return fmt.Sprintf("Enter a value for cMeters bigger than %v", gc.iSnap.CMeters)
	} else if consumed < 0 {
		return fmt.Sprintf("Consumed < 0", consumed)
	}

	if gc.iSnap.CPrice <= 0 {
		return "Cprice <= 0"
	}

	gc.iSnap.Consumed = consumed
	total := float64(consumed) * gc.iSnap.CPrice
	var err error
	gc.iSnap.Total, err = strconv.ParseFloat(fmt.Sprintf("%.2f", total), 64)
	if err != nil {
		return "Could not parse IaddGasSnap "
	}
	gc.Snaps = append(gc.Snaps, gc.iSnap)

	return fmt.Sprintf("valid: Added:\n%v", gc.iSnap)

}

func SliceIndex(limit int64, predicate func(i int64) bool) int64 {
	var i int64
	for i = 0; i < limit; i++ {
		if predicate(i) {
			return i
		}
	}
	return -1
}

func (gComp *GasComponent) Select(v string) string {
	cmeters, err := strconv.ParseInt(v, 10, 64)
	pos := SliceIndex(int64(len(gComp.Snaps)), func(i int64) bool {
		return gComp.Snaps[i].CMeters == cmeters
	})

	if pos != -1 {
		gComp.SelectedSnapPos = pos
	} else {
		return "invalid Position could not be selected"
	}

	if err == nil {
		return "valid"
	} else {
		return "invalid"
	}

}

// func (gComp *GasComponent) Select(v string) string {
// 	cmeters, err := strconv.ParseInt(v, 10, 64)
// 	gComp.SelectedSnapPos = SliceIndex(int64(len(gComp.Snaps)), func(i int64) bool {
// 		return gComp.Snaps[i].CMeters == cmeters
// 	})

// 	if err == nil {
// 		return "valid"
// 	} else {
// 		return "invalid"
// 	}

// }

// You only need to delete the last entry or else all values after delete
// need to be recalculated
func (gComp *GasComponent) DeleteLast() {
	if len(gComp.Snaps) < 1 {
		return
	}
	gComp.Snaps = gComp.Snaps[:len(gComp.Snaps)-1]
}

var valid_cPrice bool = false

func (gc *GasComponent) UcPrice(price string) string {
	if v, err := strconv.ParseFloat(price, 64); err == nil {
		if v <= 0 {
			valid_cPrice = false
			return "invalid Value is less or equal to zero"
		}
		gc.Snaps[gc.SelectedSnapPos].CPrice = v
		valid_cPrice = true
		return "valid"
	} else {
		valid_cPrice = false
		return "invalid"
	}
}

var valid_cMeters bool = false

func (gc *GasComponent) UcMeters(cMeters string) string {
	m, errU := strconv.ParseInt(cMeters, 10, 64)

	if gc.SelectedSnapPos == -1 {
		valid_cMeters = false
		return "invalid Not selected"
	}
	// current_snap := gc.Snaps[gc.SelectedSnapPos]
	previousValue := gc.Snaps[gc.SelectedSnapPos-1].CMeters

	if int(gc.SelectedSnapPos+1) == len(gc.Snaps) && errU == nil {
		if m < previousValue {
			valid_cMeters = false
			return "invalid Current snap value smaller than than last month value"
		}
		gc.Snaps[gc.SelectedSnapPos].CMeters = m
		valid_cMeters = true
		return "valid"

	} else if errU == nil {
		nextValue := gc.Snaps[gc.SelectedSnapPos+1].CMeters

		if m < previousValue || m > nextValue {
			valid_cMeters = false
			return fmt.Sprintf("invalid cMeters need to be between %v %v", previousValue, nextValue)
		}

		gc.Snaps[gc.SelectedSnapPos].CMeters = m
		valid_cMeters = true

		return "valid"
	} else {
		valid_cMeters = false
		return "invalid"
	}
}

func (gc *GasComponent) UCalcSelectedSnap() string {
	if valid_cMeters == false || valid_cPrice == false {
		return "Some update values are false"
	}
	if gc.SelectedSnapPos < 0 {
		return "GasSnap position not selected"
	}
	slc := gc.Snaps[gc.SelectedSnapPos]

	if gc.SelectedSnapPos == 0 {
		slc.Consumed = 0
	} else {
		slc.Consumed = slc.CMeters - gc.Snaps[gc.SelectedSnapPos-1].CMeters
	}

	total := float64(slc.Consumed) * slc.CPrice
	slc.Total, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", total), 64)

	if total < 0 {
		return "invalid total < 0"
	}

	if gc.SelectedSnapPos != int64(len(gc.Snaps)-1) {
		afterSelected := &gc.Snaps[gc.SelectedSnapPos+1]

		afterSelected.Consumed = afterSelected.CMeters - slc.CMeters
		total = afterSelected.CPrice * float64(afterSelected.Consumed)
		afterSelected.Total, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", total), 64)
	}

	slc2 := &gc.Snaps[gc.SelectedSnapPos]

	slc2.Date = slc.Date
	slc2.CMeters = slc.CMeters
	slc2.CPrice = slc.CPrice
	slc2.Consumed = slc.Consumed
	slc2.Total = slc.Total

	return fmt.Sprintf("valid: slc.Total=%v", slc2.Total)
}
