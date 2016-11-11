
## Frontend Plus

Here are some additional frontend APIs other than the basic ones.


### GET /v2/dsusages/:appID/:senderID

Get the numbers of install and open introduced by the specific sender.

If succeeded, returns HTTP 200 OK code and response:
```json
{
    "new_install":1,
    "new_open":2
}
```

### DELETE /v2/dsusages/:appID/:senderID

Clear the numbers of install and open introduced by the specific sender.

If succeeded, returns HTTP 200 OK code.
