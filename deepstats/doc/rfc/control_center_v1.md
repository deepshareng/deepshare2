# DeepStats API V1

DeepStats provides a RESTful service to serve the data and power up dashboard.

## Definition

(TODO: Add description to all terms)

- Channel
- Counter Event
- App
- Account

## Endpoints

### GET /developers/:developer_id/channels
List aggregated counters associated with a developer ID.

If succeeded, returns HTTP 200 OK code and response:
```json
{
  "channels": [
    {"id": "...", "name": "...", "description": "..."},
    "..."
  ]
}
```

### GET /channels/:channel_id/counters
Get aggregated results of counters of a channel.
Usually filtered by events in query parameters.
- The results of one event are sorted in monotonically increasing order of time.
- "timestamp" is RFC3339 format.
- Granularity is day. We will support more in the future.


Get parameters:

`?event=install&event=...`

If none of events is given, it means "all" events.

If succeeded, returns HTTP 200 OK code and response:
```json
{
  "counters": [
    {"event": "install", "counts": [
        {"timestamp":"2006-01-02T00:00:00Z00:00", "count": 10},
        {"timestamp":"2006-01-01T00:00:00Z00:00", "count": 20},
    ]},
    "..."
  ]
}

```
