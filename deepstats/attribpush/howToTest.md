##1.start attribution service

  - NSQ consumer
  - NSQ producer
  - http server that accept request: PUT /apps/callbackurl/:appid
  - push service, push attributed events to callback servers
   
```
go run cmd/deepstatsd/attribpush/main.go
```

##2.start a fake callback server

  - http server, accept events from attribution service

```
go run deepstats/attribution/hack/test_callback_server.go
```

##3.PUT http://127.0.0.1:8081/apps/callbackurl/1652E90881C1FAE8

```json
    {
      "url":"http://127.0.0.1:8082/callback"
    }
```

##4.start deepshare2 server

```
go run cmd/deepshared/main.go -nsq-url=127.0.0.1:4150 -log-level=debug
```

##5.run API test to generate events

```
go test ./test/testapi -server-addr=http://127.0.0.1:8080
```
