package main

import (
	"encoding/json"
	"fmt"
	"github.com/goinggo/mapstructure"
	"io/ioutil"
	"reflect"
	"time"
	"strconv"
)

var _ = time.ANSIC
var _ = ioutil.Discard
var _ = reflect.Array
var _ = mapstructure.Error{}

func check(err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
}

type gasSct struct {
	Date     time.Time `json:"date"`
	CMeters  int64     `json:"cMeters"`
	Consumed int64     `json:"consumed"`
	CPrice   float64   `json:"cPrice"`
	Total    float64   `json:"total"`
}

type gasComponent struct {
	Snaps []gasSct
	SelectedSnapPos int64
	iSnap gasSct
}

func (gc *gasComponent) SaveDataToFile() string {
	b, err := json.Marshal(gc.Snaps)
	check(err)

	if err = ioutil.WriteFile("gas_data1.json", b, 0644); err != nil {
		return fmt.Sprintf("Err: %v \nCould not write to file", err.Error())
	}

	return "File saved ok"
}


func (gc *gasComponent) LoadJsonGas() {
	jsonGasData, err := ioutil.ReadFile("gas_data1.json")
	check(err)

	var gasData []map[string]interface{}

	err = json.Unmarshal(jsonGasData, &gasData)
	check(err)

	for _, gasSnap := range gasData {
		dateS, _ := gasSnap["date"].(string)
		date, _ := time.Parse("2006-01-02T15:04:05.000Z",dateS)

		gc.Snaps = append(gc.Snaps, gasSct{
			Date: date,
			CMeters: int64(gasSnap["cMeters"].(float64)),
			Consumed: int64(gasSnap["consumed"].(float64)),
			CPrice: gasSnap["cPrice"].(float64),
			Total: gasSnap["total"].(float64),
		})
	}

}

func (gc *gasComponent) IDate(dateS string) string {
	if date, err := time.Parse("2006-01-02T15:04:05.000Z",dateS); err != nil {
		return "invalid"
	} else {
		gc.iSnap.Date = date
		return "valid"
	}

}

func (gc *gasComponent) ICmeters(cMeters string) string {
	if m , err := strconv.ParseInt(cMeters, 10, 64); err == nil {
		gc.iSnap.CMeters = m
		return "valid"
	} else {
		return "invalid"
	}
}

func (gc *gasComponent) ICprice(price string) string {
	if v , err := strconv.ParseFloat(price,64); err == nil {
		gc.iSnap.CPrice = v
		return "valid"
	} else {
		return "invalid"
	}
}

func (gc *gasComponent) IaddGasSnap() {
	consumed := gc.iSnap.CMeters - gc.Snaps[len(gc.Snaps) -1].CMeters
	if consumed <= 0 {
		fmt.Printf("Consumed is %v", consumed)
		return
	}

	if gc.iSnap.CPrice <= 0 {
		fmt.Println("Cprice <= 0")
		return
	}

	gc.iSnap.Consumed = consumed
	total := float64(consumed) * gc.iSnap.CPrice
	gc.iSnap.Total, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", total), 64)
	gc.Snaps = append(gc.Snaps, gc.iSnap)

}

func SliceIndex(limit int, predicate func(v int) bool) int {
	for i := 0; i < limit; i++ {
        if predicate(i) {
            return i
        }
    }
    return -1
}

func (gComp *gasComponent) Select(v string) {
	cmeters, _ := strconv.ParseInt(v, 10, 64)
	gComp.SelectedSnapPos = int64(
		SliceIndex(len(gComp.Snaps), func(i int) bool {
			return gComp.Snaps[i].CMeters == cmeters
	}))

}

// You only need to delete the last entry or else all values after delete
// need to be recalculated
func (gComp *gasComponent) DeleteLast()  {
	gComp.Snaps = gComp.Snaps[:len(gComp.Snaps)-1]
}

func (gc *gasComponent) UcPrice(price string) string {
	if v , err := strconv.ParseFloat(price,64); err == nil {
		gc.Snaps[gc.SelectedSnapPos].CPrice = v
		return "valid"
	} else {
		return "invalid"
	}
}

func (gc *gasComponent) UcMeters(cMeters string) string {
	if m , err := strconv.ParseInt(cMeters, 10, 64); err == nil {
		gc.Snaps[gc.SelectedSnapPos].CMeters = m
		return "valid"
	} else {
		return "invalid"
	}
}

func (gc *gasComponent) UCalcSelectedSnap() {
	slc := &gc.Snaps[gc.SelectedSnapPos]

	if gc.SelectedSnapPos == 0 {
		slc.Consumed = 0
	} else {
		slc.Consumed = slc.CMeters - gc.Snaps[gc.SelectedSnapPos -1].CMeters
	}

	total := float64(slc.Consumed) * slc.CPrice
	slc.Total, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", total), 64)
	
	if gc.SelectedSnapPos != int64(len(gc.Snaps) - 1) {
		afterSelected := &gc.Snaps[gc.SelectedSnapPos + 1]

		afterSelected.Consumed = afterSelected.CMeters - slc.CMeters
		total = afterSelected.CPrice * float64(afterSelected.Consumed)
		afterSelected.Total, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", total), 64)
	}

}
	


func main() {
	gComp := gasComponent{}
	gComp.LoadJsonGas()
//	fmt.Println(gComp.Snaps)
//	fmt.Println(gComp.iSnap){"date":"2013-12-15T09:26:14.635","cMeters":22279,"consumed":0,"cPrice":0.0,"total":0.0}

//	fmt.Println(gComp.IDate("2015-11-15T16:55:14.635Z"))
//	gComp.ICmeters("28058")
//	gComp.ICprice("1.45")
//	gComp.IaddGasSnap()

//	fmt.Println(gComp.SaveDataToFile())

//	gComp.Select("27518")
//	fmt.Println(gComp.SelectedSnapPos)
//	gComp.UcMeters("27600")
//	gComp.UcPrice("10.0")
//	gComp.UCalcSelectedSnap()
//	gComp.DeleteLast()

//	for _, v := range gComp.Snaps[18:] {
//		fmt.Printf("%v\n", v)
//	}
//
	d1 := "2014-12-15T09:26:14.635Z"
	d2, err := time.Parse("2006-01-02T15:04:05.000Z", d1)
	check(err)
	v, err := d2.MarshalText();
	check(err)
	fmt.Println(string(v))


}

















