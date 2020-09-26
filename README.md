# golang-api
Implement Data Caching Service using micro service architecture 

This golang api microservice is using below library:<br>
github.com/gorilla/mux <br>
github.com/Shopify/sarama <br>
github.com/go-redis/redis/v7 <br>
github.com/go-sql-driver/mysql <br>
github.com/rs/zerolog/log <br>

Api will be running on localhost:8080 <br>
<h3>Endpoints:</h3>
GET /getbyid/{id}+{limit} <br>
POST /addData 
<h5>RequestBody:</h5>
{
Id: 1,
Name: "Jim"
}
