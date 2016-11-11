# User/Developer API Guide V2

## Backend Endpoints

### PUT /v2/matches/:appID/:cookieID
Bind in-app data and sender info (tracking) with current user session.

This is called by the JS server.

PUT body example:
```json
{
  "sender_info": {
    "sender_id": "123",
    "channels": ["x","y","z"]
  },
  "inapp_data": "bytes_of_user_data",
  "client_ip": "ip of receiver device",
  "client_ua": "ua of receiver device(browser)"
}
```

If succeeded, returns HTTP 200 OK code.


### GET /v2/matches/:appID/:cookieID
Get in-app data based on cookieID, this performs exact match.

GET Parameters (for tracking):

`?receiver_info={unique_id:ddd,other_info:...}
&tracking=install/open`

The following params are required:

- **unique_id** field in receiver_info: receiver uniqueID.
- **tracking** for tracking, could be install or open.

If succeeded, returns HTTP 200 OK code and response:
```json
{
  "inapp_data": "bytes_of_user_data"
}
```

### GET /v2/matches/:appID
Get in-app data for current user session, based on user agent and ip, non-exact match.

GET Parameters (for tracking):


`?receiver_info={unique_id:uuu,hardware_id:hhh,other_info:...}
&tracking=install/open
&client_ip=ip_of_receiver
&client_ua=userAgentOfReveiver`

The following params are required:

- **unique_id** field in receiver_info
- **tracking**
- **client_ip**
- **client_ua**
- **hardware_id** field in receiver_info is required when "tracking" is "install"


If succeeded, returns HTTP 200 OK code and response:
```json
{
  "inapp_data": "bytes_of_user_data"
}
```

### PUT /v2/devicecookie/cookies/:deviceID
request
```json
{
  "cookie_id": "a_uuid_generated_by_js_server"
}
```
If succeeded, returns HTTP 200 OK code.

### GET /v2/devicecookie/cookies/:deviceID
If succeeded, returns HTTP 200 OK code and response:
```json
{
  "cookie_id": "a_uuid_generated_by_frontend_or_js_server"
}
```

### PUT /v2/devicecookie/devices/:cookieID
request
```json
{
  "unique_id": "device uniqueID"
}
```
If succeeded, returns HTTP 200 OK code.

### GET /v2/devicecookie/devices/:cookieID
If succeeded, returns HTTP 200 OK code and response:
```json
{
  "unique_id": "device uniqueID"
}
```

### GET /v2/appcookiedevice/:appID/:cookieID
If succeeded, returns HTTP 200 OK code and response:
```json
{
  "unique_id": "device unique id"
}
```

### PUT /v2/appcookiedevice/:appID/:cookieID
PUT body example:
```json
{
  "unique_id": "device unique id"
}
```
If succeeded, returns HTTP 200 OK code.

### POST /v2/dsactions/:appID
SDK or JS side of DeepShare can push pre-defined actions by calling this endpoint.
The service will simply push the action to message queue in form of event for future data analyse.

POST body example:
```json
{
  "receiver_info": {
    "unique_id": "receiver unique ID"
  },
  "action": "<action_name>",
  "kvs": {
    "customised_params":"..."
  }
}
```
#### Supported actions:
- app/* sent from SDK
  - app/close

    `lapse` holds the time spent on different API calls, within one session (from app open/install to app close), in milli-seconds.
    ```json
    {
      "receiver_info": {
        "unique_id": "receiver unique ID"
      },
      "action": "app/close",
      "kvs": {
        "lapse": {
            "inappdata": [1,2,1],
            "genurl": [],
            "attribute": [],
            "other actions...": []
        }
      }
    }
    ```
- js/* sent from JS
  - js/dst
  - js/deeplink

### POST /v2/counters/:appID
Developer can push custom actions by calling this endpoint. Receiver unique_id is attached for tracking. For attribution tracking, please make sure the unique_id here is the same as the unique_id in GET /matches.

**Event** are user defined actions, e.g. “buy”, "reward", etc.
All events are aggregated in our backend pipeline system.

POST body example:
```json
{
  "receiver_info": {
    "unique_id": "receiver unique ID"
  },
  "counters": [
    {"event":"buy",   "count": 5},
    {"event":"reward","count": 10}
  ]
}
```
