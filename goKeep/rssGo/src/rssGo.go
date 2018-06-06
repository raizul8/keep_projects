package main

import (
	"encoding/xml"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

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

var mangaRssAddresses = []string{
	"http://mangafox.me/rss/the_gamer.xml",
	"http://mangafox.me/rss/ruler_of_the_land.xml",
	// "http://mangafox.me/rss/one_piece.xml",
	"http://mangafox.me/rss/mob_psycho_100.xml",
	"http://mangafox.me/rss/girls_of_the_wild_s.xml",
	"http://mangafox.me/rss/the_scholar_s_reincarnation.xml",
	// "http://mangafox.me/rss/the_god_of_high_school.xml",
	// "http://mangafox.me/rss/tower_of_god.xml",
	"http://mangafox.me/rss/black_haze.xml",
	"http://mangafox.me/rss/silver_spoon.xml",
	"http://mangafox.me/rss/newman.xml",
	// "http://mangafox.me/rss/stretch.xml",
	"http://mangafox.me/rss/piano_no_mori.xml",
	// "http://mangafox.me/rss/where_tangents_meet.xml",
	"http://mangafox.me/rss/boku_no_hero_academia.xml",
	// "http://mangafox.me/rss/kingdom.xml",
	// "http://mangafox.me/rss/horimiya.xml",
	"http://mangafox.me/rss/red_storm.xml",
	"http://mangafox.me/rss/doupo_cangqiong.xml",
	// "http://mangafox.me/rss/untouchable_massstar.xml",
	// "http://mangafox.me/rss/tamen_de_gushi.xml",
	"http://mangafox.me/rss/onepunch_man_one.xml",
	"http://mangafox.me/rss/relife.xml",
	"http://mangafox.me/rss/tomo_chan_wa_onnanoko.xml",
	"http://mangafox.me/rss/fuuka.xml",
	// "http://mangafox.me/rss/special_martial_arts_extreme_hell_private_high_school.xml",
	// "http://mangafox.me/rss/douluo_dalu_ii_jueshui_tangmen.xml",
	"http://mangafox.me/rss/the_breaker_new_waves.xml",
	// "http://mangafox.me/rss/gantz_g.xml",
}

type RssStruct struct {
	rss RSS
	nr  int
}

func getRss(addrr string, nr int, ch chan RssStruct) {
	// you have to check if resp gets somethin else than ok
	resp, err := http.Get(addrr)
	if err != nil {
		return
	}


	XMLdata, _ := ioutil.ReadAll(resp.Body)

	var rss RSS
	xml.Unmarshal(XMLdata, &rss)

	rss.Items = rss.Items[:5]
	ch <- RssStruct{rss, nr}
}

func main() {
	ch1 := make(chan RssStruct)

	for i := range mangaRssAddresses {
		go getRss(mangaRssAddresses[i], i, ch1)
	}

	ln := len(mangaRssAddresses)
	rssList := make([]RSS, ln)

	for i := 0; i < ln; i++ {
		rStruct := <-ch1
		rssList[rStruct.nr] = rStruct.rss
	}

	dat, err := ioutil.ReadFile("template.html")
	check(err)
	tmpl, err := template.New("t1").Parse(string(dat))
	check(err)

	htmlFile, err := os.Create("res.html")
	check(err)
	defer htmlFile.Close()

	tmpl.Execute(htmlFile, rssList)

}
