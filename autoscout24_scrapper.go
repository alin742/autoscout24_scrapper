package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	// "github.com/jezek/xgb/render"
)

type BodyType uint8
type Fuel string
type Year uint32
type Filter struct {
	make      string
	model     *string
	bodytypes *[]BodyType
	fuel      *[]Fuel
	fregfrom  *Year
	fregto    *Year
	pricefrom *uint32
	priceto   *uint32
}

const (
	NONE BodyType = iota
	HATCHBACK
	CABRIO
	COUPE
	SUV
	STATIONWAGON
	LIMOUSINE
	VAN
	TRANSPORTER
	OTHER
)

const (
	BENZIN        Fuel = "B"
	DESIL              = "D"
	ETHANOL            = "M"
	ELECTRIC           = "E"
	HYDROGEN           = "H"
	HYBRID_BENZIN      = "2"
	HYBRID_DESIL       = "3"
)

func bodyTypeArrayToString(a []BodyType, delim string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(a), " ", delim, -1), "[]")
	//return strings.Trim(strings.Join(strings.Split(fmt.Sprint(a), " "), delim), "[]")
	//return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(a)), delim), "[]")
}

func fuelArrayToString(a []Fuel, delim string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(a), " ", delim, -1), "[]")
}

func yearToString(y Year) string {
	current_year, _, _ := time.Now().Date()
	out := ""
	if y >= 1950 && int(y) < current_year {
		out += strconv.Itoa(int(y))
	}
	return out
}

func priceToString(p uint32) string {
	return strconv.Itoa(int(p))
}

func buildRequest(url string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/119.0")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	// req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Referer", "https://www.autoscout24.de/lst/toyota?atype=C&body=3%2C6&cy=D&damaged_listing=exclude&desc=0&fregto=2000&ocs_listing=include&powertype=kw&search_id=1ozcomz0eur&sort=standard&source=listpage_pagination&ustate=N%2CU")
	req.Header.Set("x-nextjs-data", "1")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Connection", "keep-alive")
	// req.Header.Set("Cookie", `as24Visitor=24817c73-c55d-4ad6-b5c7-7d0d27a9a320; euconsent-v2=CPwzWsAPwzWsAGNABCDEDSCgAAAAAELAAAAAAAAQhABAFS4gAKAgJCbQMIgEQIgrCAiAUAAAAkBBAAAECAgwRgEIsBkAAEAABAAABABBAACAAACAAAAAAAAgAAAAAAQAAAAAAAQAICAAIAAgAAAAAAUAQAAAAAAAAAIAIABABCgACAEkoEAAAAAHAAAAAACAQAAAAAAAAAQAAAAAAAAAAAAAAQERUAIAZYC5BkAIAZYC5B0DQABYAFQAMgAcgA-AEEAMgA0AB4AEQAJgATwAqwBcAF0AMQAZgA3gBzAD0AH6AQwBEgCWAE0AKMAUoAwwBogD2gH4AfoBAwCLAEdAJMASkAp4BcwC8gGKAOoAi8BIgCVAEyAKPAU2AtgBcgDBgGSAMnAZYA1gBxYDxyUBkABYAGQAOAAfAB4AEQAJgAVQAuABigEMARIAjgBRgD8AKeAXMAxQB1AEXgJEAUeAtgBk4DWAIQlIFAACwAKgAZAA5AB8AIAAaAA8ACIAEwAJ4AUgAqgBiADMAHMAP0AhgCJAFGAKUAaMA_AD9AIsAR0AlIBcwC8gGKAOoAi8BIgCmwFsALkAZIAycBlkDWANZAcEA8cCEIQAUABsAEgA0gBzgEHAJ2AWcAzQDFgGQhIFwACwAKgAZAA5AB4AIIAZABoADwAIgATAAngBVADeAHMAPQAfgBCQCGAIkARwAlgBNAClAGGAMsAe0A_AD9AIGARoAkwBKQCngFzAMUAaIBIgCjwFIgKbAWwAuQBgwDJAGTwNYA1kBwQDxwIQhgBAAiwBRgDnAOoApsBiwDWQHjiAA4AJABFgDSAHOAREBrIDxxwA4AEgAUABlgDnAHdAQcBCACIgE7ALOAZkBiwDIQGVAMzIgAgAAgBCKACoARABIAC0ARwAywBzgDuAIOATsA_4DFg0AMAZYBTwFyCIAYAywCngLkIQCgA9ACOAKeAXMAxQB1AEqALkAZOA8cgADAGWAOcAzJIAMAMsAdwBBwDMgMWAeOAA.YAAAAAAAA4CA; cconsent-v2=%7B%22purpose%22%3A%7B%22legitimateInterests%22%3A%5B25%5D%2C%22consents%22%3A%5B%5D%7D%2C%22vendor%22%3A%7B%22legitimateInterests%22%3A%5B10211%2C10218%2C10441%5D%2C%22consents%22%3A%5B%5D%7D%7D; addtl_consent=1~; as24-cmp-signature=L8vegNhGRFTQ%2BqXLOBpen8Bu90jryG2n4tx%2FRVynLQGdGMoM4Cd%2FKDQFip2VS8i%2BtxlqpLvVivt3qO5YA5dIy%2BNDxMP4bMxuO0q%2FHa41HHcPqAshije0oanh92JLhrIjR7ecqDfXbzYbaXbJ30R1A4id9V9QCH91UtMy7yCegoc%3D; optimizelyEndUserId=oeu1692492529898r0.333425464824997; last-seen-listings=[{"id":"f0c54914-e46c-4f9b-8eb5-7d96a0a1bed2","price":2600,"mileage":95200,"power":63,"firstRegistrationYear":1998,"makeId":70,"modelId":2052,"fuelId":"B","latitude":51.51176,"longitude":11.8815,"isSmyle":false,"isLeasingMarktPremium":null,"timestamp":1698487044},{"id":"fbbaa7b0-dfbf-4f92-95b1-50e7f46b4214","price":17500,"mileage":117000,"power":129,"firstRegistrationYear":1995,"makeId":70,"modelId":2058,"fuelId":"B","latitude":50.20149,"longitude":8.57491,"isSmyle":false,"isLeasingMarktPremium":null,"timestamp":1698487041},{"id":"e8f83a61-8b7d-4e82-beb7-81efdd90de66","price":5950,"mileage":120171,"power":72,"firstRegistrationYear":1991,"makeId":70,"modelId":2047,"fuelId":"B","latitude":53.55161,"longitude":10.05886,"isSmyle":false,"isLeasingMarktPremium":null,"timestamp":1698486620},{"id":"366a2239-d6de-4b66-a3f7-56710b18ec8f","price":5500,"mileage":347000,"power":72,"firstRegistrationYear":1991,"makeId":70,"modelId":2047,"fuelId":"B","latitude":50.0532,"longitude":8.68616,"isSmyle":false,"isLeasingMarktPremium":null,"timestamp":1698486600},{"id":"331ef794-892f-420d-9170-d17e80ffae0b","price":25000,"mileage":110000,"power":93,"firstRegistrationYear":1977,"makeId":70,"modelId":2050,"fuelId":"B","latitude":48.4082,"longitude":10.04195,"isSmyle":false,"isLeasingMarktPremium":null,"timestamp":1698486584},{"id":"02de9af8-48be-43f2-9227-9b027a836901","price":12900,"mileage":73000,"power":40,"firstRegistrationYear":1978,"makeId":70,"modelId":2052,"fuelId":"B","latitude":50.88569,"longitude":6.46309,"isSmyle":false,"isLeasingMarktPremium":null,"timestamp":1698486541},{"id":"0d06ddcb-9ca4-4218-8447-6f7e072d5367","price":3699,"mileage":100000,"power":55,"firstRegistrationYear":1992,"makeId":70,"modelId":2052,"fuelId":"B","latitude":48.21914,"longitude":7.77869,"isSmyle":false,"isLeasingMarktPremium":null,"timestamp":1698486504},{"id":"a171cbf1-1240-46b1-a4df-f92bcebe0262","price":3500,"mileage":239570,"power":77,"firstRegistrationYear":1992,"makeId":70,"modelId":2050,"fuelId":"B","latitude":48.34488,"longitude":8.4015,"isSmyle":false,"isLeasingMarktPremium":null,"timestamp":1698486480},{"id":"42c5ca23-2722-4559-b1a2-615b10688b2c","price":10900,"mileage":139965,"power":85,"firstRegistrationYear":1988,"makeId":70,"modelId":2058,"fuelId":"B","latitude":50.30288,"longitude":8.56934,"isSmyle":false,"isLeasingMarktPremium":null,"timestamp":1698486463},{"id":"299a65bc-eb89-47fe-ae23-6ee4542e9a13","price":1600,"mileage":99850,"power":71,"firstRegistrationYear":2000,"makeId":70,"modelId":2052,"fuelId":"B","latitude":49.08753,"longitude":11.21847,"isSmyle":false,"isLeasingMarktPremium":null,"timestamp":1698486446}]; culture=de-DE; ab_test_lp=%7B%7D; search-session=true; ab_test_dp=%7B%7D`)
	req.Header.Set("TE", "trailers")
	return client.Do(req)
}

func search(filter *Filter) {
	var data map[string]interface{}
	body_types_string := ""
	fuel_string := ""
	model := ""
	// page := 0
	if filter.model != nil {
		model = "/" + *filter.model
	}
	if filter.bodytypes != nil {
		body_types_string = "&body=" + bodyTypeArrayToString(*filter.bodytypes, "%2C")
	}
	if filter.fuel != nil {
		fuel_string = "&fuel=" + fuelArrayToString(*filter.fuel, "%2C")
	}
	if filter.fregfrom != nil {
		yearToString(*filter.fregfrom)
	}
	url := "https://www.autoscout24.de/_next/data/as24-search-funnel_main-4490/lst/" +
		filter.make + model + ".json?atype=C" + body_types_string + //"&cy=D&damaged_listing=exclude&desc=0" +
		fuel_string // + "&ocs_listing=include&powertype=kw&search_id=1ozcomz0eur&sort=standard" +
		// "&source=listpage_pagination" //+ "&ustate=N%2CU&page=2&slug=toyota"
	resp, err := buildRequest(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Printf("%s", bodyText)
	err = json.Unmarshal(bodyText, &data)
	if err != nil {
		log.Fatal(err)
	}
	pageProps := data["pageProps"].(map[string]interface{})
	listings := pageProps["listings"].([]interface{})
	for id, listing := range listings {
		listing_map := listing.(map[string]interface{})
		vehicle := listing_map["vehicle"].(map[string]interface{})
		fmt.Printf("%d: %s %s\n", id, vehicle["make"], vehicle["model"])
		fmt.Println(vehicle)
	}
}

func main() {
	search(&Filter{
		make: "toyota",
	})
}
