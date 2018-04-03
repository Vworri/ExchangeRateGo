package main
import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"errors"
	"os"
	"strings"
	"strconv"
	"github.com/jasonlvhit/gocron"
	"encoding/json"
	"time"


)
var layout = "2006-01-02";
func main(){
	//Example query: http://localhost:8080/api?date=2018-03-28&dest=usd&orig=aud&amount=4
	GetEchangeRateFile() ;
	http.HandleFunc("/api",ApiResponse);
	http.HandleFunc("/info",InfoResponse);
	go http.ListenAndServe(":8080", nil);
	DailyTask();
}

func ApiResponse(res http.ResponseWriter, req *http.Request){

	if req.Method != "GET"{
		http.Error(res, http.StatusText(405), 405)
		res.Write([]byte("This request method is not supported"))
		return
	}
	var err  error;
	var amountStr string;
	exchange := new(ExResp);
	if req.URL.Query()["date"] != nil{
	exchange.RateDate =  req.URL.Query()["date"][0];
	}else {exchange.RateDate = time.Now().Format(layout)}
	exchange.Orig = req.URL.Query()["orig"][0];
	exchange.Dest  = req.URL.Query()["dest"][0];
	if  req.URL.Query()["amount"] != nil{
		amountStr =  req.URL.Query()["amount"][0];
	}else{amountStr = "1"}
	exchange.OriginalAmount, err  = strconv.ParseFloat(amountStr, 64)
	if err != nil {
		log.Fatal(err.Error());
	}
	ex := RetrieveLocalData();
	exchange.RateDate, exchange.Rate, err = ex.At(exchange.RateDate, exchange.Orig, exchange.Dest)
	if err != nil {
		log.Fatal(err.Error());
	}
	exchange.ResultAmount = strconv.FormatFloat(exchange.OriginalAmount*exchange.Rate, 'f', 6, 64);
	json.NewEncoder(res).Encode(exchange)
	return
}

func InfoResponse(res http.ResponseWriter, req *http.Request){
	info := new(Info)
	ex := RetrieveLocalData();
	info.MaxDt, info.MinDt = ex.getMInMaxDates()
	json.NewEncoder(res).Encode(info)
	return
	}


func DailyTask(){
	gocron.Every(1).Day().At("16:45").Do(GetEchangeRateFile);
	<- gocron.Start();
	gocron.Remove(GetEchangeRateFile);
	gocron.Clear();
}

func GetEchangeRateFile() {
	resp, err := http.Get ("https://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist-90d.xml");
	log.Println("STARTED: Collecting exchange rates");
	if err != nil {
		log.Fatal(err.Error());
	}
	if resp.StatusCode != http.StatusOK {
		log.Fatal("response status from ECB not OK");
	}
	rss, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err.Error());
	}
	writeErr := ioutil.WriteFile("ExchangeData.xml", []byte(rss),0644);
	if writeErr != nil {
		log.Fatal(err.Error());
	}
	log.Println("FINISHED: Collecting exchange rates");
	
}



func RetrieveLocalData() *ExchangeData{
	dataFile, err := os.Open("ExchangeData.xml")
	if err != nil {
		log.Fatal(err.Error());
	}
	defer dataFile.Close()
	data, readError := ioutil.ReadAll(dataFile)
	if readError != nil {
		log.Fatal(err.Error());
	}
	newData := new(ExchangeData);
	err = xml.Unmarshal(data, &newData);
	if err != nil {
		log.Fatal(err.Error());
	}
	return newData;
}

func (exData *ExchangeData) At(t string,dest string, orig string)(string, float64, error){
	data := exData.FullData;
	num := float64(0);
	den := float64(0); 
	layout := "2006-01-02"
	date ,_ := time.Parse(layout, t);
	max, min := exData.getMInMaxDates()
	if date.After(max) || date.Before(min){
		t = max.Format(layout)
	}
	for _, entry := range data {
		if entry.Time == t {
			for _, currency := range entry.Data {
				if currency.Currency == strings.ToUpper(dest) {
					den = currency.Rate
				}
				if currency.Currency == strings.ToUpper(orig) {
					num = currency.Rate
				}
				if num != 0 && den != 0{
					return t, (num /den), nil
				}
			}
		}
	}
	return t, 0.0 , errors.New("destination or origin currency not found")
}


func  (exData *ExchangeData) getMInMaxDates()(time.Time, time.Time){
	data := exData.FullData;
	layout := "2006-01-02"
	max_date,err := time.Parse(layout, layout);
	if err != nil{
		log.Panic(err.Error())
	}
	min_date := time.Now()
	for _, entry := range data {
		thisTime, timeErr := time.Parse(layout, entry.Time)
		if timeErr != nil{
			log.Panic(timeErr.Error())
		}
		if ( thisTime.Before(min_date)) {
			min_date =  thisTime
		}
		if (thisTime.After(max_date)) {
			max_date = thisTime
		}

	}

	return  max_date, min_date
}

type ExchangeData struct {
	XMLNmae       xml.Name      `xml:"gesmes:Envelope"`
	FullData  []RateByTime `xml:"Cube>Cube"`
}

type RateByTime struct{
	Time 	string `xml:"time,attr"`
	Data []RateByCountry `xml:"Cube"`


}
type RateByCountry struct {
	Currency string `xml:"currency,attr"`
	Rate float64 `xml:"rate,attr"`
}

type ExResp struct {
	Dest string `json:"destination"`
	Orig string `json:"origin"`
	OriginalAmount float64 `json:"originalAmount"`
	Rate float64 `json:"rate"`
	RateDate string `json:"rateDate"`
	ResultAmount string `json:"resultAmount"`
}


type Info struct {
	MaxDt time.Time `json:"max"`
	MinDt time.Time `json:"min"`
}