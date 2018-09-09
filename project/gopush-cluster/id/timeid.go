

package id

import (
	"time"
)

/*
type TimeID struct {
	lastID int64
}

// NewTimeID create a new TimeID struct
func NewTimeID() *TimeID {
	return &TimeID{lastID: 0}
}

// ID generate a time ID
func (t *TimeID) ID() int64 {
	for {
		s := time.Now().UnixNano() / 100
		if t.lastID >= s {
			// if last time id > current time, may be who change the system id,
			// so sleep last time id minus current time
			panic("time delay!!!!!!")
		} else {
			// save the current time id
			t.lastID = s
			return s
		}
	}
	return 0
}
*/

// Get get a time id.
func Get() int64 {
	return time.Now().UnixNano() / 100
}
