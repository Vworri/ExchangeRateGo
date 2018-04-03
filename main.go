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


)

func main(){
	//Example query: http://localhost:8080/api?date=2018-03-28&dest=usd&orig=aud&amount=4
	GetEchangeRateFile() ;
	http.HandleFunc("/api",ApiResponse);
	go http.ListenAndServe(":8080", nil);
	DailyTask();


}

func ApiResponse(res http.ResponseWriter, req *http.Request){
	if req.Method != "GET"{
		http.Error(res, http.StatusText(405), 405)
		res.Write([]byte("This request method is not supported"))
		return
	}
	date :=  req.URL.Query()["date"][0];
	orig := req.URL.Query()["orig"][0];
	dest  := req.URL.Query()["dest"][0];
	amountStr :=  req.URL.Query()["amount"][0];
	amount, err  :=strconv.ParseFloat(amountStr, 64)
	if err != nil {
		log.Fatal(err.Error());
	}
	ex := RetrieveLocalData();
	rate, err := ex.At(date, orig, dest)
	if err != nil {
		log.Fatal(err.Error());
	}
	val := strconv.FormatFloat(amount*rate, 'f', 6, 64);
	res.Write([]byte(val));
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

func (exData *ExchangeData) At(t string,dest string, orig string)(float64, error){
	data := exData.FullData;
	num := float64(0);
	den := float64(0);
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
					return (num /den), nil
				}
			}
		}
	}
	return 0.0 , errors.New("destination or origin currency not found")
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
