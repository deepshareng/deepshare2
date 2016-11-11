# Deepstatsd/Retention
## Description
Deepstatsd/Retention is used to consume message from NSQ, resend retention message if match 

## Usage of deepstatsd/Retention:
```
  // Obligatory   
  -nsq-url string
      The nsqd url to produce message

  // Select a way to access NSQD, either nsqd or nsqlookupd
  -nsqsel string
      Specify the way to get nsq message, nsqlookupd/nsqd (default "nsqlookupd")
  -nsqd-tcp-addr string
      Specify the nsqd adress (default "")
  -nsqlookupd-http-addr string
      Specify the nsqlookupd adress (default "")

  // redis
  -redis-addr string
      The redis url to use as DB
  -redis-password string
      The redis password to use as DB
 
  // Optional
  -channel string
      Specify the NSQ channel for consumer (default "deepstats_retention")
  -topics string
      Specify the NSQ topic for consume, flag format should be topic1, topic2... (default "counter,sharelink,dsaction,genurl,inappdata,match")


  //retention service config
  -retention-day int
      Specify the duration of retention calculated (default 3)
  -retention-service-name string
      Specify the unique retention service (default "3-day-retention")
  -retention-topic string
      Specify the NSQ Topic for retention service produce its message (default "retention")
```

## Run

``` 
./retention -mongo-addr -nsqsel -nsqd-tcp-addr/-nsqlookupd-http-addr
```
