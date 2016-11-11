# Deepstatsd/Backup
## Description
Deepstatsd/Backup is used to backup NSQ message into our persistent storage.

## Usage of deepstatsd:
```
  // Obligatory   
  // Select a way to backup our message
  -backupsel string
      Specify the ways to backup our message, currenly only support localfs/mongo-compress-day/mongo-raw. (default "mongo-raw")
  -path string
        Local filesystem backup path(support localfs backup service)
  -mongo-addr string
      Specify the raw data mongo database URL(support mongo backup service)


  // Select a way to access NSQD, either nsqd or nsqlookupd
  -nsqsel string
        Specify the way to get nsq message, nsqlookupd/nsqd (default "")
  -nsqd-tcp-addr string
        Specify the nsqd adress (default "")
  -nsqlookupd-http-addr string
        Specify the nsqlookupd adress (default "no_nsqlookupd")
  
  // Optional
  -channel string
        Specify the NSQ channel for consumer (default "deepstats_backup")
  -topics string
      Specify the NSQ topic for consume, flag format should be topic1, topic2... (default "counter,match,sharelink,dsaction,genurl,inappdata")
  -mongocoll string
      Specify the Mongo collection (default "backup")
  -mongodb string
      Specify the Mongo database (default "deepstats")
```

## Run

``` 
./backup -backsel -path/-mongo-addr -nsqsel -nsqd-tcp-addr/-nsqlookupd-http-addr
```
