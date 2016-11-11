# Deepstatsd/Appevent
## Description
Deepstatsd/Appevent is used to consume message from NSQ, distribute message to appevent service.

## Usage of deepstatsd:
```
  // Obligatory   
  -mongo-addr string
        Specify the raw data mongo database URL

  // Select a way to access NSQD, either nsqd or nsqlookupd
  -nsqsel string
      Specify the way to get nsq message, nsqlookupd/nsqd (default "nsqlookupd")
  -nsqd-tcp-addr string
      Specify the nsqd adress (default "")
  -nsqlookupd-http-addr string
      Specify the nsqlookupd adress (default "")
  
  // Optional
  -channel string
      Specify the NSQ channel for consumer (default "deepstats_appevent")
  -mongocoll string
      Specify the Mongo collection (default "appevent")
  -mongodb string
      Specify the Mongo database (default "deepstats")
  -topics string
      Specify the NSQ topic for consume, flag format should be topic1, topic2... (default "counter,sharelink,dsaction,genurl,inappdata,match")
```


## Run

``` 
./appevent -mongo-addr -nsqsel -nsqd-tcp-addr/-nsqlookupd-http-addr
```
