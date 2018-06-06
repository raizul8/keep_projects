//package main
//
//import (
//	"encoding/json"
//	"fmt"
//	"github.com/goinggo/mapstructure"
//	"io/ioutil"
////	"reflect"
//	"time"
//)
//
//var _ = time.ANSIC
//var _ = ioutil.Discard
////var _ = reflect.Array
//var _ = mapstructure.Error{}
//
//func check(err error) {
//	if err != nil {
//		fmt.Println(err.Error())
//	}
//}
////time.Time
//type gasSnap struct {
//	Date   	 string    `jpath:"date"`
//	CMeters  int       `jpath:"cMeters"`
//	Consumed int       `jpath:"consumed"`
//	CPrice   float64   `jpath:"cPrice"`
//	Total    float64   `jpath:"total"`
//}
//
//type Profession struct {
//	Job string `jpath:"job"`
//	Exp int `jpath:"exp"`
//}
//
//type NameDoc struct {
//	Name string `jpath:"name"`
//	Age  int 	`jpath:"age"`
//	Profession Profession `jpath:"profession"`
//}
//
//var document = `[{
//	"name":"bill",
//	"age":20,
//	"profession": {
//		"job": "chemist",
//		"exp": 200
//	}
//}]`
//
//func main() {
////	sliceScript := []byte(document)
////	sliceScript, err := ioutil.ReadFile("gas_data1.json")
//	sliceScript := []byte(`[{"date":"2013-12-15T09:26:14.635","cMeters":22279,"consumed":0,"cPrice":0.0,"total":0.0}]`)
////	check(err)
////	fmt.Println(string(sliceScript) )
//	var sliceMap []map[string]interface{}
//	json.Unmarshal(sliceScript, &sliceMap)
//
//	var myslice []gasSnap
//	err := mapstructure.DecodeSlicePath(sliceMap, &myslice)
//	check(err)
//
//	gas1 := myslice[0]
////	fmt.Printf("%#v \n%v \n%v  \n%v \n%v",gas1.date, gas1.cMeters, gas1.consumed, gas1.cPrice, gas1.total)
//	fmt.Println(gas1)
//}
