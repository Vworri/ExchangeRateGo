sudo docker stop currency_ex_service 
sudo docker rm currency_ex_service 
sudo docker build -t exchangeservice:test .
sudo docker run --name=currency_ex_service -d  -p 8080:8080 exchangeservice:test