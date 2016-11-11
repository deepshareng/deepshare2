// Package instrumentation defines behaviors to instrument the deepshare service

package instrumentation

import "time"

// We keep them separate for now, just in case they evolve in different direction.

type Shorturl interface {
	HTTPGetDuration(start time.Time)     // overall time performing GET HTTP Request
	HTTPPostDuration(start time.Time)    // overall time performing POST HTTP Request
	StorageGetDuration(start time.Time)  // time spent getting data from storage
	StorageSaveDuration(start time.Time) // time spent saving data to the storage
}

type InappData interface {
	HTTPPostDuration(start time.Time) // overall time performing POST HTTP Request
}

type JsApi interface {
	HTTPPostDuration(start time.Time) // overall time performing POST HTTP Request
}

type Sharelink interface {
	HTTPGetDuration(start time.Time) // overall time performing GET HTTP Request
}

type Token interface {
	HTTPGetDuration(start time.Time) // overall time performing GET HTTP Request
}

type Match interface {
	HTTPGetDuration(start time.Time)       // overall time performing GET HTTP Request
	HTTPPostDuration(start time.Time)      // overall time performing POST HTTP Request
	StorageGetDuration(start time.Time)    // time spent getting data from storage
	StorageSaveDuration(start time.Time)   // time spent saving data to the storage
	StorageDeleteDuration(start time.Time) // time spent deleting data from the storage
	StorageHSetDuration(start time.Time)   // time spent on storage HSet
	StorageHGetDuration(start time.Time)   // time spent on storage HGet
	StorageHDelDuration(start time.Time)   // time spent on storage HDel
}

type Counter interface {
	HTTPGETDuration(start time.Time)
	AggregateDuration(start time.Time)
}

type AppCookieDeviceInstrument interface {
	HTTPGetDuration(start time.Time)    // overall time performing GET HTTP Request
	HTTPPutDuration(start time.Time)    // overall time performing POST HTTP Request
	StorageGetDuration(start time.Time) // time spent getting data from storage
	StoragePutDuration(start time.Time) // time spent saving data to the storage
}

type AppInfoInstrument interface {
	HTTPGetDuration(start time.Time)    // overall time performing GET HTTP Request
	HTTPPutDuration(start time.Time)    // overall time performing POST HTTP Request
	StorageGetDuration(start time.Time) // time spent getting data from storage
	StoragePutDuration(start time.Time) // time spent saving data to the storage
}

type DSUsageInstrument interface {
	HTTPGetDuration(start time.Time)       // overall time performing GET HTTP Request
	HTTPDeleteDuration(start time.Time)    // overall time performing POST HTTP Request
	StorageGetDuration(start time.Time)    // time spent getting data from storage
	StorageDeleteDuration(start time.Time) // time spent deleting data from the storage
	StorageIncDuration(start time.Time)    // time spent inc data to the storage
}

type DeviceCookierInstrument interface {
	HTTPGetDuration(start time.Time)    // overall time performing GET HTTP Request
	HTTPPutDuration(start time.Time)    // overall time performing POST HTTP Request
	StorageGetDuration(start time.Time) // time spent getting data from storage
	StoragePutDuration(start time.Time) // time spent saving data to the storage
}

type BindDeviceToCookierInstrument interface {
	HTTPGetDuration(start time.Time) // overall time performing GET HTTP Request
}
