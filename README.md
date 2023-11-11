1. create database `test`
2. create table `Server` with columns `id` and `status`
2. `go get github.com/go-sql-driver/mysql`
2. `go get -u github.com/gin-gonic/gin`
3. `go run main.go`

# Create Server
```shell
curl --location --request POST 'http://localhost:9000/servers'
```
# Short polling
Curl every second
```shell
curl --location 'http://localhost:9000/short/status/1'
```

# Long Polling
```shell
curl --location 'http://localhost:9000/long/status/1?status=TODO'
```
