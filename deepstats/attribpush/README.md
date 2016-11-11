## Register callback url

POST /apps/callbackurl/:appID

json```
{
    "url": "http://host:port/callback"
}
```

The url should be able to accept Attribution data in the following format:

json```
[
  {
    "sender_id": "testsender",
    "tag": "ds/open",
    "value": 100,
    "timestamp": 1448466001
  }
]
```

## Attribution flow

{event (sender_id maybe empty)}

-> attribution parser ->

{attributed event (with sender_id extracted)}

-> produce to nsq ->

{attributed event (with sender_id extracted)}

-> push service consumer ->

{attribution push info}

-> push to developer registered callback url