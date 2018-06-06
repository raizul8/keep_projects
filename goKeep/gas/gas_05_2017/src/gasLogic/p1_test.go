package gasLogic_test

import (
	"fmt"
	"raz_utils"
	"strings"
	"testing"
	"time"

	"strconv"

	"gasLogic"

	"github.com/stretchr/testify/assert"
)

func TestGetGasSnapListFromJsonFile(t *testing.T) {
	gasSlice, err := gasLogic.GetGasSliceFromJsonFile(gasLogic.GasData1JsonPath)
	if err != nil {
		t.Fatal(err)
	}
	//fmt.Println("gasSlice", gasSlice[1], "path", p1.GasData1JsonPath)

	if len(gasSlice) == 0 || gasSlice == nil {
		t.Fatal("Slice either nil or empty")
	}
}

//func TestGetPathString(t *testing.T) {
//	path, err := p1.GetPathString()
//	fmt.Println("path:", path)
//	if err != nil {
//		t.Fatalf("path: %s err: %s", path, err)
//	}
//
//	if !strings.Contains(path, "src") {
//		t.Fatal("GetPathString did not contain src.")
//	}
//
//}

func TestInitialWriteToGobFile(t *testing.T) {
	t.SkipNow()

	if err := gasLogic.InitialWriteToGobFile(gasLogic.GobFilePath); err != nil {
		t.FailNow()
	}
}

func TestLoadFromGobFile(t *testing.T) {
	gasSlice, err := gasLogic.LoadFromGobFile(gasLogic.GobFilePath)
	if err != nil {
		t.Fatal(err)
	}
	if len(gasSlice) == 0 {
		t.Fatal("gasSlice from LoadFromGobFile is empty.")
	}

}

func TestSaveToGobFile(t *testing.T) {
	//t.SkipNow()
	gasSlice, _ := gasLogic.LoadFromGobFile("gobReadGas")
	bucharest, _ := time.LoadLocation("Europe/Bucharest")

	lastGasChanged := gasSlice[len(gasSlice)-1]
	//fmt.Println("Last gas before change\n", lastGasChanged)
	lastGasChanged.Date = time.Date(2017, 2, 15, 15, 0, 0, 0, bucharest)
	lastGasChanged.CMeters = 33100

	gobWriteGasPath := "gobWriteGas"
	gasSlice[len(gasSlice)-1] = lastGasChanged
	//fmt.Println("lastGasChanged:\n", lastGasChanged)

	err := gasLogic.SaveToGobFile(gobWriteGasPath, gasSlice)
	if err != nil {
		t.Fatal(err)
	}

	gasSlice, _ = gasLogic.LoadFromGobFile(gobWriteGasPath)
	lastGasReloaded := gasSlice[len(gasSlice)-1]
	//fmt.Println("last gas loaded from file after change:\n", lastGasReloaded)

	if lastGasChanged.Date.Hour() != lastGasReloaded.Date.Hour() && lastGasChanged.CMeters != lastGasReloaded.CMeters {
		t.Fatal("Could not save to file")
	}

}

func TestGasComponent(t *testing.T) {
	gc := gasLogic.GasComponent{}
	// initialize GasComponent with gas slice
	gsnaps, err := gasLogic.LoadFromGobFile("gobReadGas")
	if err != nil {
		t.Fatal(err)
	}

	if len(gsnaps) == 0 {
		t.Fatal("could not load gas slice from gobTest1")
	}

	gc.Snaps = gsnaps

	// insert data into GasComponent.ISnap: Date , cMeters, cPrice
	date1 := "2017/03/28 16:16"
	err = gc.Idate(date1)
	if err != nil {
		t.Fatal(err)
	}

	//fmt.Println("last gas", gc.Snaps[len(gc.Snaps)-1])
	err = gc.IcMetersS("3100")
	serr := fmt.Sprintf("%v", err)
	if !strings.Contains(serr, "invalid:") {
		t.Fatal(err)
	}

	err = gc.IcMetersS("33100")
	serr = fmt.Sprintf("%v", err)
	if strings.Contains(serr, "invalid:") {
		t.Fatal(err)
	}

	//fmt.Println("", gc.ISnap.CMeters)
	assert.Equal(t, int64(33100), gc.ISnap.CMeters, "cMeters not equal")

	gc.IcPrice("1.3")
	assert.Equal(t, gc.ISnap.CPrice, float64(1.3))

	// Insert the gc.ISnap to slice gc.snaps
	err = gc.IaddGasSnap()
	if err != nil {
		t.Fatal(err)
	}

	// Check if the insertion is the same with our data
	lastGas := gc.Snaps[len(gc.Snaps)-1]

	date2, _ := time.Parse("2006/01/02 15:04", date1)

	//fmt.Println("date2       ", date2, "\nlastGas.Date", lastGas.Date)
	assert.Equal(t, date2, lastGas.Date)
	assert.Equal(t, int64(33100), lastGas.CMeters)
	assert.Equal(t, lastGas.CPrice, float64(1.3))

	// Save our data to a test json file
	jFile := "saveGas.json"
	err = gc.SaveGasSliceToJsonFile(jFile)
	if err != nil {
		t.Fatal(err)
	}

	// Load our data from the test json file
	gc.Snaps = []gasLogic.GasSct{}
	err = gc.LoadFromJsonFile(jFile)
	if len(gc.Snaps) == 0 || err != nil {
		t.Fatal("Could not load gas slice from json file")
	}

	// Check if the data loaded from jFile is the same with what we wrote
	assert.Equal(t, date2, lastGas.Date)
	assert.Equal(t, int64(33100), lastGas.CMeters)
	assert.Equal(t, lastGas.CPrice, float64(1.3))
	//fmt.Println("", gc.Snaps[len(gc.Snaps)-1])

	// Save to gobWriteGas2 then load and check the data
	gobW2 := "gobWriteGas2"
	err = gc.SaveToGobFile(gobW2)
	if err != nil {
		t.Fatal(err)
	}

	// Load from gob file and check our input
	err = gc.LoadFromGobFile(gobW2)
	if err != nil || len(gc.Snaps) == 0 {
		t.Fatal("Problem loading from gob file")
	}

	assert.Equal(t, date2, lastGas.Date)
	assert.Equal(t, int64(33100), lastGas.CMeters)
	assert.Equal(t, lastGas.CPrice, float64(1.3))

	// test delete last
	len1 := len(gc.Snaps)
	lastVal := gc.Snaps[len1-1]

	if err := gc.DeleteLast(); err != nil {
		t.Fatalf("Delete error: %v", err)
	}

	// check if gc.Snaps len is actually smaller by one
	len2 := len(gc.Snaps)
	assert.Equal(t, len1-1, len2)
	// check that the previous last value is no longer last value
	assert.NotEqual(t, lastVal, gc.Snaps[len2-1])

}

func TestGasComponent_Select(t *testing.T) {
	gc := new(gasLogic.GasComponent)
	gc.LoadFromGobFile(gasLogic.GobFilePath)

	snapsLen := len(gc.Snaps)
	fmt.Println("len:", snapsLen)

	// value to be updatet
	cMetersStoUpdate := "32446"
	cMetersNewValueS := "32447"

	cMetersNewValue, err := strconv.ParseInt(cMetersNewValueS, 10, 64)
	if err != nil {
		// invalid string
	}

	pos := gc.Select(cMetersStoUpdate)
	fmt.Printf("update for cMeters: %s pos: %d with new val %s \n", cMetersStoUpdate, pos, cMetersNewValueS)

	slc := []gasLogic.GasSct{}

	for _, v := range gc.Snaps[pos:snapsLen] {
		slc = append(slc, v)
	}

	//for _, v := range slc {
	//	fmt.Println("slc: ", v)
	//}

	// delete from gc.Snaps all after pos
	gc.Snaps = gc.Snaps[:pos]
	fmt.Println("last snap befor inserting new slice:\n", gc.Snaps[len(gc.Snaps)-1])

	// set the new value
	slc[0].CMeters = cMetersNewValue

	// reinsert the values from slc back into gc.Snaps
	for _, v := range slc {
		//
		err = gc.IcMetersS(strconv.FormatInt(v.CMeters, 10))
		if err != nil {
			t.Fatal("Could not input cMeters.\n", err)
		}

		gc.ISnap.CPrice = v.CPrice
		gc.ISnap.Date = v.Date

		gc.IaddGasSnap()
	}

	for _, v := range gc.Snaps[pos:] {
		fmt.Println("New Vals: ", v)
	}

}

func TestGasComponent_Select2(t *testing.T) {
	gc := new(gasLogic.GasComponent)
	gc.LoadFromGobFile(gasLogic.GobFilePath)
	cMetersStoUpdate := "32446"
	cMetersNewValueS := "32447"
	gSnapLen := len(gc.Snaps)

	for i := gc.Select(cMetersStoUpdate) - 1; i < gSnapLen; i++ {
		fmt.Println(gc.Snaps[i])
	}

	fmt.Println("Update:", cMetersStoUpdate, "to", cMetersNewValueS)
	err := gc.Update(cMetersStoUpdate, cMetersNewValueS)

	if err != nil {
		t.Fatal(err)
	}

	for i := gc.Select(cMetersNewValueS) - 1; i < gSnapLen; i++ {
		fmt.Println(gc.Snaps[i])
	}
}

func init() {
	//fmt.Println("------p1_test_inside---------")

	path, err := raz_utils.GetRunningProgramDir()
	raz_utils.LogErrorFatal(err)

	err = gasLogic.InitializePaths(path)
	raz_utils.LogErrorFatal(err)

}
