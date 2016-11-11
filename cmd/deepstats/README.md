# Deepstats

Deepstats apply a HTTP handler handling dashboard client request for app/channel infos for their needs.

This version only support testing usage. 

## Usage of deepstatsd:
```
  // Obligatory 
  -mongo-addr string
        Specify the raw data mongo database URL

  // Optional
  -http-listen string
        HTTP/HTTPs Host and Port to listen on (default "0.0.0.0:16759")
  -mongocoll-appchannel string
        Specify the Mongo collection for app channel (default "appchannel")
  -mongocoll-appevent string
        Specify the Mongo collection for app event (default "appevent")
  -mongocoll-day-aggregate string
        Specify the Mongo collection for aggregate (default "day")
  -mongocoll-hour-aggregate string
        Specify the Mongo collection for aggregate (default "hour")
  -mongocoll-total-aggregate string
        Specify the Mongo collection for aggregate (default "total")

  -mongodb-appchannel string
        Specify the Mongo database for app channel (default "deepstats")
  -mongodb-appevent string
        Specify the Mongo database for app event (default "deepstats")
  -mongodb-day-aggregate string
        Specify the Mongo database for aggregate (default "deepstats")
  -mongodb-hour-aggregate string
        Specify the Mongo database for aggregate (default "deepstats")
  -mongodb-total-aggregate string
        Specify the Mongo database for aggregate (default "deepstats")
```

## Run
Support 5 service, need corresponding service running.
``` 
./deepstats -mongo-addr -http-listen
```
