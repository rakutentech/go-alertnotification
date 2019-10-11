package alertnotification

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/GitbookIO/diskache"
)

// Throttler struct storing disckage directory and Throttling duration
type Throttler struct {
	CacheOpt         string
	ThrottleDuration int
}

// ErrorOccurrence store error time and error
type ErrorOccurrence struct {
	StartTime time.Time
	ErrorType error
}

// NewThrottler constructs new Throttle struct and init diskcache directory
func NewThrottler() Throttler {

	t := Throttler{
		CacheOpt:         fmt.Sprintf("/tmp/cache/%v_throttler_disk_cache", os.Getenv("APP_NAME")),
		ThrottleDuration: 5, // default 5mn
	}
	if len(os.Getenv("THROTTLE_DURATION")) != 0 {
		duration, err := strconv.Atoi(os.Getenv("THROTTLE_DURATION"))
		if err != nil {
			return t
		}
		t.ThrottleDuration = duration
	}

	if len(os.Getenv("THROTTLE_DISKCACHE_DIR")) != 0 {
		t.CacheOpt = os.Getenv("THROTTLE_DISKCACHE_DIR")
	}

	return t
}

// IsThrottled checks if the error has been throttled. If not, throttle it
func (t *Throttler) IsThrottled(ocError error) bool {
	dc, err := t.getDiskCache()
	if err != nil {
		return false
	}
	cachedTime, throttled := dc.Get(ocError.Error())

	if throttled && !isOverThrottleDuration(string(cachedTime), t.ThrottleDuration) {
		// already throttled and not over throttling duration, do nothing
		return true
	}
	// if it has not throttled yet or over throttle duration, throttle it and return false to send notification
	// Rethrottler will also renew the timestamp in the throttler cache.
	if err = t.ThrottleError(ocError); err != nil {
		return false
	}
	return false
}

func isOverThrottleDuration(cachedTime string, throttleDuration int) bool {
	throttledTime, err := time.Parse(time.RFC3339, string(cachedTime))
	if err != nil {
		return false
	}
	now := time.Now()
	diff := int(now.Sub(throttledTime).Minutes())
	return diff > throttleDuration
}

// ThrottleError throttle the alert within the limited duration
func (t *Throttler) ThrottleError(errObj error) error {
	dc, err := t.getDiskCache()
	if err != nil {
		return err
	}
	now := time.Now().Format(time.RFC3339)
	err = dc.Set(errObj.Error(), []byte(now))

	return err
}

// CleanThrottlingCache clean all the diskcache in throttling cache directory
func (t *Throttler) CleanThrottlingCache() (err error) {
	dc, err := t.getDiskCache()
	if err != nil {
		return err
	}
	err = dc.Clean()
	return err
}

func (t *Throttler) getDiskCache() (*diskache.Diskache, error) {
	opts := diskache.Opts{
		Directory: t.CacheOpt,
	}
	return diskache.New(&opts)
}
