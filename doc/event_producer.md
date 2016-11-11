When deepshare backend APIs are called, we produce some logs to the message queue, later consumed and analysed by deepstats, we call the logs as **events**.

## Topic
Every event belongs to a topic.
Different topics can be consumed by different consumers.

Now we have the following topics:

- counter
- match

We can also use appID as topic when we need to scale the message queue by apps.

## Content
The content of event. Should be json data, with the following fields:

- **AppID** 
- **Event**
- **Channels**: required for match/\*, optional for other events.
- **SenderID**: required for match/\*, optional for other events.
- **UniqueID**: required for match/install, match/open, counters/* and dsactions/* when called from SDK.
- **HardwareID**: receiver hardware ID, the field is only set for match/install.
- **Count**
- **UAInfo**:User agent and ip address and other information that included in http header. 
- **KVs**: Extra data, Different events can have different type of Kvs, typically field of request or response when the corresponding API is called.


###Exmple events are listed as follows:

- **/v2/matches/bind**: produced when PUT matches is called
- **/v2/matches/install**: produced when GET matches is called with "tracking=install"
- **/v2/matches/open**: produced when GET matches is called with "tracking=open"
- **/v2/matches/unknown**: produced when GET matches is called without tracking parameter or tracking is invalid
- **/v2/dasctions/close**: when app closed, SDK should push a close event through dsactions.
- **/v2/counters/***: * can be defined by user, for example: buy,reward,etc.

As a convention, please use web API path prefix as event prefix(or use API path prefix as event name)

### /v2/matches/bind
```json
{
    "AppID": "<the_app_id>",
    "Event": "match/bind",
    "Channels": ["channel_x", "channel_y", "channel_z"],
    "SenderID": "Who is sharing the link",
    "Count": 1,
    "UAInfo": {
        "ua": "ua string",
        "ip": "client ip",
        "os": "ios",
        "os_version": "9.1",
        "browser": "chrome 45.0"
    },
    "KVs": {
        "ua": "UA of client",
        "cookie_id": "some_cookie_id",
        "inapp_data": "app_data_put_by_client",
    }
}
```
#### /v2/matches/install
```json
{
    "AppID": "<the_app_id>",
    "Event": "match/install",
    "Channels": ["channel_x", "channel_y", "channel_z"],
    "SenderID": "who shared the link and caused this install",
    "UniqueID": "identity of receiver device",
    "Count": 1,
    "UAInfo": {
      "ua": "ua string",
      "ip": "client ip",
      "os": "ios",
      "os_version": "9.1"
    },
    "KVs": {
        "ua": "UA used for match",
        "cookie_id": "cookieID used for match",
        "inapp_data": "in app data saved by match/bind"
    }
}
```

#### /v2/matches/open
same as match/install

#### /v2/matches/unknown
same as match/install

### /v2/dasctions/close
```json
{
    "AppID": "<the_app_id>",
    "Event": "app/close",
    "UniqueID": "identity of receiver device",
    "Count": 1,
    "UAInfo": {
      "ua": "ua string",
      "ip": "client ip",
      "os": "ios",
      "os_version": "9.1"
    }
}
```

### /v2/counters/*
```json
{
    "AppID": "<the_app_id>",
    "Event": "counters/buy",
    "UniqueID": "identity of receiver device",
    "Count": 100,
    "UAInfo": {
      "ua": "ua string",
      "ip": "client ip",
      "os": "ios",
      "os_version": "9.1"
    }
}
```
