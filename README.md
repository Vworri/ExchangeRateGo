# Golang Currency Exchange API
This API is set up to return exchange rates and evaluate currency exchanges as well as keep a local store of data up-to-date.
##  Getting started API:
	1. pull from github: https://github.com/Vworri/ExchangeRateGo.git	2. make sure you have a recent version of Docker 
    3. All you have to do is run the deploy.sh script (remember to give execute permission)
    4. The container should be up and running, mapped to port 8080
## Exchange
There are only two required fields in the query:
origin: the initial currencry from which you would like to exchange
destination: the destination currency due for the conversion
The optional parameters follow:

| name  | format  |default |
|--|--|--|
|  date |YYYY-MM-DD  | max available date in data store|
|  amount | integer  | 1|

Example Urls:
http://127.0.0.1:8080/api?dest=usd&orig=aud
### Response

```javascript
    {
	"destination": "USD",
	"origin": "AUD",
	"originalAmount": 1,
	"rate": 0.7683337490646046,
	"rateDate": "2018-03-29",
	"resultAmount": "$0.77"} 
```


http://127.0.0.1:8080/api?date=2018-04-02&dest=usd&orig=aud&amount=4

``` javascript
    {
	"destination": "USD",
	"origin": "AUD",
	"originalAmount": 4.0,
	"rate": 0.7683337490646046,
	"rateDate": "2018-03-29",
	"resultAmount": "$3.07"}
```
## Data info
There is another endpoint I added to verify the data in store
http://127.0.0.1:8080/info

```javascript
    {
	"max": "2018-03-29",
	"min": "2018-01-02"}
```


## Example Code (generated by insomnia)
```ruby
require 'uri'
require 'net/http'

url = URI("http://localhost:8080/api?date=2017-02-08&dest=usd&orig=aud&amount=10")

http = Net::HTTP.new(url.host, url.port)

request = Net::HTTP::Get.new(url)

response = http.request(request)
puts response.read_body
```
