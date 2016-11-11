package aggregate

// sort a list of AggregateCount by timestamp
// - Recent year comes first, then
// - Recent month comes first, then
// - Recent day comes first
type ByTimestamp []*AggregateCountWithTimestamp

func (a ByTimestamp) Len() int           { return len(a) }
func (a ByTimestamp) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTimestamp) Less(i, j int) bool { return a[j].Timestamp.After(a[i].Timestamp) }
