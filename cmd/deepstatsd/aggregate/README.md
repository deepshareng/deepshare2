# Deepstatsd/Aggregate
##Description
Deepstatsd/Aggregate is used to consume message from NSQ, aggregate channel's message form NSQ channel to counter persistent storage, in accordance with the aggregation rules.

## Usage of deepstatsd:
```
  // Obligatory
  -mongo-addr string
        Specify the raw data mongo database URL
  -agg-service string
      Specify which aggregate service is being used, "hour" means aggregating events by hour. Currently support hour/day/total. (default "day")

  // select a way to access NSQD, either nsqd or nsqlookupd
  -nsqsel string
        Specify the way to get nsq message, nsqlookupd/nsqd (default "nsqlookupd")
  -nsqd-tcp-addr string
        Specify the nsqd adress (default "")
  -nsqlookupd-http-addr string
        Specify the nsqlookupd adress (default "")
  
  // Optional
  -channel string
        Specify the NSQ channel for consumer (default "test1")
  -mongocoll string
        Specify the Mongo collection (default "counter")
  -mongodb string
        Specify the Mongo database (default "deepstats")    
  -topics string
      Specify the NSQ topic for consume, flag format should be topic1, topic2... (default "counter,sharelink,dsaction,genurl,inappdata,match")
```

## Run
In our general service support, we always need to run three services.
```
./aggregate -mongo-addr -nsqsel -nsqd-tcp-addr/-nsqlookupd-http-addr -mongocoll=counter -agg-service=day -nsqChannel=deepstats_day-aggregate
./aggregate -mongo-addr -nsqsel -nsqd-tcp-addr/-nsqlookupd-http-addr -mongocoll=total -agg-service=total -nsqChannel=deepstats_total-aggregate
./aggregate -mongo-addr -nsqsel -nsqd-tcp-addr/-nsqlookupd-http-addr -mongocoll=hour -agg-service=hour -nsqChanel=deepstats_hour-aggregate
```
