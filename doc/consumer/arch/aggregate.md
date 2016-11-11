# Aggregate service
##Architecture
We implement aggregate service into following phases:
- Define the granularity by the implementation interface, converting the data time to the time user wants.
- Collect user data into temporary storage, convert the data time by our rules
- Commit temporary changes to persistent DB backend
- Response client request aggregated data by duration, eventFilter, and channel. 

On the basis of user concrete aggregate granularity, ConvertTimeToGranularity will regulate the aggregation data time.
E.g.
If user want to aggregate data by day,
For any time in a day, like
        20060102T10:00:00Z00:00
It might aggreate into this time
        20060102T00:00:00Z00:00

##Definition
```
type AggregateService interface {
    Insert(aggregate *CounterEvent) error

    QueryDuration(appid string, channel string, eventFilters []string, start time.Time, granularity time.Duration, limit int) ([]*api.AggregateResult, error)

    Aggregate() error

    ConvertTimeToGranularity(time.Time) time.Time
}
```