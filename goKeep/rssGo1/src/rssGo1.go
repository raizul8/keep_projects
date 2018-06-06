package main

import (
	"encoding/xml"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
)

type Item struct {
	XMLName xml.Name `xml:"item"`
	Title   string   `xml:"title"`
	Link    string   `xml:"link"`
}

func (it Item) String() string {
	return fmt.Sprintf("%s %s", it.Title, it.Link)
}

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Title   string   `xml:"channel>title"`
	Items   []Item   `xml:"channel>item"`
}

//var mangaRssAddresses = []string{
//	"http://mangafox.me/rss/the_gamer.xml",
//	"http://mangafox.me/rss/ruler_of_the_land.xml",
//	"http://mangafox.me/rss/one_piece.xml",
//	"http://mangafox.me/rss/mob_psycho_100.xml",
//	"http://mangafox.me/rss/the_scholar_s_reincarnation.xml",
//}

var mangaRssAddresses = []string{
	"http://mangafox.me/rss/the_gamer.xml",
	"http://mangafox.me/rss/ruler_of_the_land.xml",
	"http://mangafox.me/rss/one_piece.xml",
	"http://mangafox.me/rss/mob_psycho_100.xml",
	// "http://mangafox.me/rss/girls_of_the_wild_s.xml",
	"http://mangafox.me/rss/the_scholar_s_reincarnation.xml",
	// "http://mangafox.me/rss/the_god_of_high_school.xml",
	// "http://mangafox.me/rss/tower_of_god.xml",
	"http://mangafox.me/rss/black_haze.xml",
	"http://mangafox.me/rss/silver_spoon.xml",
	"http://mangafox.me/rss/newman.xml",
	// "http://mangafox.me/rss/stretch.xml",
	"http://mangafox.me/rss/piano_no_mori.xml",
	// "http://mangafox.me/rss/where_tangents_meet.xml",
	"http://mangafox.me/rss/kingdom.xml",
	// "http://mangafox.me/rss/horimiya.xml",
	"http://mangafox.me/rss/red_storm.xml",
	"http://mangafox.me/rss/doupo_cangqiong.xml",
	// "http://mangafox.me/rss/untouchable_massstar.xml",
	"http://mangafox.me/rss/tamen_de_gushi.xml",
	"http://mangafox.me/rss/onepunch_man_one.xml",
	"http://mangafox.me/rss/relife.xml",
	"http://mangafox.me/rss/tomo_chan_wa_onnanoko.xml",
	"http://mangafox.me/rss/fuuka.xml",
	// "http://mangafox.me/rss/special_martial_arts_extreme_hell_private_high_school.xml",
	// "http://mangafox.me/rss/douluo_dalu_ii_jueshui_tangmen.xml",
	"http://mangafox.me/rss/the_breaker_new_waves.xml",
	"http://mangafox.me/rss/gantz_g.xml",
	"http://mangafox.me/rss/boku_no_hero_academia.xml",
}

type RssStruct struct {
	rss RSS
	nr  int
}

type nrAddr struct {
	nr   int
	addr string
}

func genFirst(done <-chan struct{}, nrOfAddresses int) <-chan nrAddr {
	out := make(chan nrAddr, nrOfAddresses)

	defer close(out)
	for i, v := range mangaRssAddresses {
		select {
		case out <- nrAddr{nr: i, addr: v}:
		case <-done:
			return out
		}
	}

	return out
}

func mangaWorker(done <-chan struct{}, nrAddress nrAddr, outRssStruct chan<- RssStruct, wg *sync.WaitGroup) {
	c := make(chan struct{}, 1)
	go func() {
		time.Sleep(time.Second * 4)
		c <- struct{}{}
	}()

	doneCalled := false

	// This function is used to call wg.Done()
	doneCall := func(where, info string) {
		if doneCalled {
			fmt.Printf("X_%v Done already called: %v for %v\n", string(info[4:6]), where, info)
		} else {
			doneCalled = true
			wg.Done()

			if strings.Contains(where, "Timed") {
				color.Red("%v Done called %v for %v\n", string(info[4:6]), where, info)
			} else {
				color.Green("%v Done called %v for %v\n", string(info[4:6]), where, info)
			}

		}
	}

	positionAndAddress := fmt.Sprintf("nr: %v adr: %v", nrAddress.nr, nrAddress.addr[23:len(nrAddress.addr)-4])

	resp, err := http.Get(nrAddress.addr)

	if err != nil || resp.StatusCode != 200 {
		log.Println("error for :", nrAddress.addr, "StatusCode:", resp.StatusCode)
		outRssStruct <- RssStruct{}
		doneCall("inside resp.StatusCode", positionAndAddress)
		return
	}

	XMLdata, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("ioutil.ReadAll(resp.Body) err:", err)

	}

	var rss RSS
	rssStruc1 := RssStruct{}

	err = xml.Unmarshal(XMLdata, &rss)

	if err != nil {
		log.Println("xml.Unmarshal:", err, "for", positionAndAddress)
		rssStruc1 = RssStruct{}
	} else {
		rss.Items = rss.Items[:5]
		rssStruc1 = RssStruct{rss, nrAddress.nr}
	}

	select {
	case outRssStruct <- rssStruc1:
		//fmt.Println("Done case", positionAndAddress)
		doneCall("In Select", positionAndAddress)
	case <-c:
		doneCall("Timed out", positionAndAddress)
	case <-done:
		return
	}

}

func mangaWorkerPool(done <-chan struct{}, nA <-chan nrAddr) <-chan RssStruct {
	var wg sync.WaitGroup
	outRssStruct := make(chan RssStruct)

	// as much workers as adresses
	wg.Add(len(mangaRssAddresses))

	for adrs := range nA {
		go mangaWorker(done, adrs, outRssStruct, &wg)
	}

	go func() {
		wg.Wait()
		close(outRssStruct)
	}()

	return outRssStruct
}

func init() {
	//log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetFlags(log.Lshortfile)
}

func main() {

	//fmt.Println("--------START------------")

	done := make(chan struct{})
	defer close(done)

	genCh := genFirst(done, len(mangaRssAddresses))
	ch1 := mangaWorkerPool(done, genCh)

	rssList := make([]RSS, len(mangaRssAddresses))

	for v := range ch1 {
		rssList[v.nr] = v.rss
	}

	dat, err := ioutil.ReadFile("template.html")
	if err != nil {

		log.Println("err with template.html:", err)
		return
	}

	tmpl, err := template.New("t1").Parse(string(dat))
	if err != nil {
		log.Println("Err with template Parse", err)
	}

	htmlFile, err := os.Create("res.html")
	if err != nil {
		log.Println("Err with os create:", err)
	}
	defer htmlFile.Close()

	tmpl.Execute(htmlFile, rssList)

}
