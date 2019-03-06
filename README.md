# SafeBoda PromoCode
This application uses `go` and `mongodb`. To get started you'd need to install go on your system and also get the mongodb server running. `gin-gonic` is used as the router and micro framework for handling request and configuring the middleware layer of the application.


# Set project flags
Running `go run main.go --help` from the root of the directory will output as is below. There are already default flags set, but just incase.

```sh
Options:

  -h, --help                                                      display help information
  -p, --port[=5000]                                               Application is running on this port
      --db-host[=mongodb://localhost:27017]                       mongoDB host
      --db-name[=safeboda]                                        mongoDB name
```
You can provide the your flag, if you dont want to use the default specified in the project.

To run the code on port :5000 run the command
  ```bash
    go run main.go
  ```

# Endpoints
1. Create An event

    Required params: `name`, `address`, `latitude`, `longitude`

    ```bash
    curl -i -X POST -H 'Content-Type: application/json' -d '{"name": "Grand Global Hotel", "address": "Grand Global Hotel, Kampala, Uganda", "latitude" : 0.3316466, "longitude": 32.5641206 }' http://localhost:5000/promo/event
    ```

2. Create A Promo

    Required params: `radius`, `amount`, `expiration_date`, `event_id`

    ```bash 
      curl -i -X POST -H 'Content-Type: application/json' -d '{"radius": 100,"amount": 1000,"expiration_date": "2019-03-19T11:45:26.371Z", "event_id" : "5c7d792763f44c82858f55ac"}' http://localhost:5000/promo/new
    ```

3. Deactivate a promo code 

    Required params: `code` - Promo code to be deactivated.
    ```bash
      curl -i -X POST  -H 'Content-Type: application/json' -d '{"code": "SAFE-4b71aa34e96d"}' http://localhost:5000/promo/deactivate
    ```
4. Validation 

    Required params: `code`, `destination`, `origin`

    ```bash
      curl -i -X POST -H 'Content-Type: application/json' -d '{"code": "SAFE-4b71aa34e96d","origin": "Grand Global Hotel, Uganda","destination" : "Serena Musa, Uganda"}' http://localhost:5000/promo/validate
    ```

5. Retrieve all promos

    ```bash
      curl -i -X GET http://localhost:5000/promo/all
    ```

6. Retrieve all active promos

    ```bash
      curl -i -X GET http://localhost:5000/promo/active
    ```

# Test
To run the test cases for the implemented services run the command from the root of the project directory
```sh
go test -count=1 -v ./...
```