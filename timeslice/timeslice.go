/*
timselice package provides a TimeSlice stuct with its methods.

TimeSlice represents a range of times bounded by two dates (time.Time) From and To. It accepts infinite boundaries (zero times) and can be chronological or anti-chronological.
*/
package timeslice

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/sunraylab/timeline/duration"
)

// Defines the chronological direction of a timeslice:
//   - AntiChronological
//   - Undefined
//   - Chronological
type Direction int

const (
	AntiChronological Direction = -1
	Undefined         Direction = 0
	Chronological     Direction = 1
)

// TimeSlice represents a range of times bounded by two dates (time.Time) From and To. Each boundary can be an infinite time.
type TimeSlice struct {
	From time.Time
	To   time.Time
}

// MakeTimeslice creates and returns a new timeslice with a defined d duration and a starting time.
//
//	If d == zero then the timeslice represents a single time.
//	If d > 0 then the given times represents the begining
//	If d < 0 then the given times represents the end
//
// panic if the given date is not defined (zero time)
func MakeTimeslice(dte time.Time, d time.Duration) TimeSlice {
	if dte.IsZero() {
		panic(dte)
	}
	ts := &TimeSlice{From: dte, To: dte.Add(d)}
	return *ts
}

// Moves the begining of the timeslice to the requested time.
// Postpone the end time if the request time exceeds it, or
// cap the date to the end of the timeslice, according to the direction.
// In case of capped or postponed move the timeslice become a single date.
func (pts *TimeSlice) MoveFrom(request time.Time, cap bool) {
	if !pts.To.IsZero() {
		if cap {
			if (pts.Direction() == Chronological && request.After(pts.To)) || (pts.Direction() == AntiChronological && request.Before(pts.To)) {
				request = pts.To
			}
		} else {
			if (pts.Direction() == Chronological && request.After(pts.To)) || (pts.Direction() == AntiChronological && request.Before(pts.To)) {
				pts.To = request
			}
		}
	}
	pts.From = request
}

// Moves the end of the timeslice to the requested time.
// Move back the begining time if the request time exceeds it, or
// cap the date to the begining of the timeslice, according to the direction.
// In case of capped or back move the timeslice become a single date.
func (pts *TimeSlice) MoveTo(request time.Time, cap bool) {
	if !pts.From.IsZero() {
		if cap {
			if (pts.Direction() == Chronological && request.Before(pts.From)) || (pts.Direction() == AntiChronological && request.After(pts.From)) {
				request = pts.From
			}
		} else {
			if (pts.Direction() == Chronological && request.Before(pts.From)) || (pts.Direction() == AntiChronological && request.After(pts.From)) {
				pts.From = request
			}
		}
	}
	pts.To = request
}

// ExtendTo add the duration at the end of the timeslice.
//
//	If the duration is negative then the end time moves backward.
//	If *pts.To is infinite, then To stays infinite
//
// The timeslice direction can change.
func (pts *TimeSlice) ExtendTo(dur duration.Duration) {
	if !pts.To.IsZero() {
		pts.To = pts.To.Add(time.Duration(dur))
	}
}

// ExtendTo add the duration at the begining of the timeslice.
//
//	If the duration is negative then the begining time moves backward.
//	If *pts.From is infinite, then From stays infinite
//
// The timeslice direction can change
func (pts *TimeSlice) ExtendFrom(dur duration.Duration) {
	if !pts.From.IsZero() {
		pts.To = pts.To.Add(time.Duration(dur))
	}
}

// String returns default formating: "{ from - to : duration }".
//
//	An infinite begining prints "past" and an infinite end prints "future".
//	if a boundary does not have any hours nor minutes nor seconds, then prints only the date.
//	if a boundary does not have any year nor month nor day, then prints only the time.
func (thists TimeSlice) String() string {
	var strfrom, strto, strdur string
	if thists.From.IsZero() {
		strfrom = "past"
	} else {
		if thists.From.Hour() == 0 && thists.From.Minute() == 0 && thists.From.Second() == 0 {
			strfrom = thists.From.Format("20060102 MST")
		} else {
			strfrom = thists.From.Format("20060102 15:04:05 MST")
		}
	}
	if thists.From.IsZero() {
		strto = "future"
	} else {
		if thists.To.Hour() == 0 && thists.To.Minute() == 0 && thists.To.Second() == 0 {
			strto = thists.To.Format("20060102 MST")
		} else if thists.From.Year() == thists.To.Year() && thists.From.Month() == thists.To.Month() && thists.From.Day() == thists.To.Day() {
			strto = thists.To.Format("15:04:05")
		} else {
			strto = thists.To.Format("20060102 15:04:05 MST")
		}
	}
	strdur = thists.Duration().FormatOrderOfMagnitude(3)
	return fmt.Sprintf("{ %s - %s : %s }", strfrom, strto, strdur)
}

// Duration returns the timeslice duration.
//
//	returns nil if one boundary is an infifinte time.
//	returns zero if timeslice boundaries have the exact same times.
func (ts TimeSlice) Duration() *duration.Duration {
	if ts.From.IsZero() || ts.To.IsZero() {
		return nil
	}
	d := duration.Duration(ts.To.Sub(ts.From))
	return &d
}

// Truncate returns the result of rounding t down to a multiple of dur (since the zero time).
// If dur <= 0, Truncate returns t stripped of any monotonic clock reading but otherwise unchanged.
func (ts TimeSlice) Truncate(dur time.Duration) TimeSlice {
	ts.From = ts.From.Truncate(dur)
	ts.To = ts.To.Truncate(dur)
	return ts
}

// Equal checks if 2 timeslices start and end at the same times, event if they're in a different timezone.
//
//	returns 1 if equal and in the same direction.
//	returns 0 if not equal.
//	returns -1 if equal but in the opposite direction.
func (one TimeSlice) Equal(another TimeSlice) int {
	if one.From.Equal(another.From) && one.To.Equal(another.To) {
		return 1
	} else if one.From.Equal(another.To) && one.To.Equal(another.From) {
		return -1
	}
	return 0
}

// Returns the direction of the timeslice.
//
//	returns 'Undefined' if both boundaries are infinite or if the timslice is a single date.
func (thists TimeSlice) Direction() Direction {
	if thists.From.IsZero() && thists.To.IsZero() {
		return Undefined
	}
	if thists.From.IsZero() {
		return AntiChronological
	}
	if thists.To.IsZero() {
		return Chronological
	}

	d := int(thists.To.Sub(thists.From))
	switch {
	case d < 0:
		return AntiChronological
	case d > 0:
		return Chronological
	default:
		return Undefined
	}
}

// Progress returns the progress rate of a given time within the timeslice, with the level of precision of the second.
//
// The progress is calculated from the begining of the timeslice, whatever its direction. The returned rate is always positive.
//
//	returns 0.5 if the timeslice has no duration
//	for a chronological timeslice:
//		- returns 0 if datetime is before the begining
//		- returns 1 if datetime is after the end
//	for an anti-chronological timeslice:
//		- returns 0 if datetime is after the begining
//		- returns 1 if datetime is before the end
func (ts TimeSlice) Progress(datetime time.Time) (rate float64) {
	pdur := ts.Duration()
	if pdur == nil {
		return 0.5
	}

	rate = datetime.Sub(ts.From).Seconds() / time.Duration(*pdur).Seconds()
	if *pdur < 0 {
		rate = -rate
	}

	// bound it between 0 an 1
	if rate < 0 {
		rate = 0
	} else if rate > 1 {
		rate = 1
	}
	return rate
}

// WhatTime returns the datetime at a certain rate within the timeslice.
//
// The progress is calculated from the begining of the timeslice, whatever its direction. The returned date is always within the time slice
//
//	returns a zero time if the timeslice has an infinite duration.
//	if the timeslice is a single date then returns it.
func (ts TimeSlice) WhatTime(rate float64) time.Time {
	pdur := ts.Duration()
	if pdur == nil {
		return time.Time{}
	}

	var t time.Time
	dprog := float64(*pdur) * rate
	t = ts.From.Add(time.Duration(dprog))

	// bount it within the timeslice boundaries
	if *pdur > 0 && t.After(ts.To) || *pdur < 0 && t.Before(ts.To) {
		t = ts.To
	}
	if *pdur > 0 && t.Before(ts.From) || *pdur < 0 && t.After(ts.From) {
		t = ts.From
	}
	return t
}

// Split a timeslice in multiple timeslices of a d duration.
//
// The end of a slice is the exact time of the begining of the next one.
// The last slice duration can be lower than d duration if thists duration is not a multiple of d.
//
// returns an error if a boundary is infinite.
//
// panic if d is <= 0
func (ts TimeSlice) Split(d time.Duration) ([]TimeSlice, error) {
	if d <= 0 {
		log.Fatalf("TimeSlice.Split with invalid duration: %v", d)
	}

	// check duration of ts
	pdur := ts.Duration()
	if pdur == nil {
		return []TimeSlice{}, errors.New("unable to split an infinite timeslice")
	}
	if *pdur < 0 {
		d = -d
	}

	slices := make([]TimeSlice, 0)
	for {
		split := MakeTimeslice(ts.From, d)
		if *pdur > 0 && split.To.After(ts.To) || *pdur < 0 && split.To.Before(ts.To) {
			split.To = ts.To
			if *split.Duration() != 0 {
				slices = append(slices, split)
			}
			break
		}
		slices = append(slices, split)
		ts.From = split.To
	}
	return slices, nil
}

// Scan returns next time, within the timeslice boundaries, matching mask.
//
// Scan always starts by the From date of the timeslice. If the From date is infinite then returns a zero date and the cursor is reset to nil.
//
// Scan use to cursor and look for the next time matching the mask and returns it. The cursor moves to this time.
// If the matching time is over the timslice boundary then returns a zero time and reset the cursor.
//
// Use fBoundaries if you want the scanner returns the boundaries even if they do not match the mask.
//
// if the timeslice has an infinite end boundary, then the scan will never returns a nil cursor.
//
// panic if mask <= 0
func (ts *TimeSlice) Scan(cursor *time.Time, mask time.Duration, fBoundaries bool) time.Time {
	if mask <= 0 {
		log.Fatalf("invalid scan mask: %d", mask)
	}
	if ts.From.IsZero() {
		return time.Time{}
	}

	var newcursor time.Time

	// init the cursor with the first scan
	first := false
	if cursor.IsZero() {
		first = true
		newcursor = ts.From
	} else {
		newcursor = *cursor
	}

	// calculate the next cursor according to the mask, and the direction
	if ts.Direction() == AntiChronological {
		// returns the begining of the timeslice
		if first && (fBoundaries || newcursor.Truncate(mask).Equal(ts.From) || newcursor.Truncate(mask).Add(mask).Equal(ts.From)) {
			newcursor = ts.From
			*cursor = newcursor
			return newcursor
		}

		// move the cursor before
		if newcursor.Equal(newcursor.Truncate(mask)) {
			newcursor = newcursor.Truncate(mask).Add(-mask)
		} else {
			newcursor = newcursor.Truncate(mask)
		}

		// check over end boundaries
		if newcursor.Before(ts.To) {
			if fBoundaries && !cursor.Equal(ts.To) {
				newcursor = ts.To
			} else {
				newcursor = time.Time{}
			}
		}
	} else {
		// returns the begining of the timeslice
		if first && (fBoundaries || newcursor.Truncate(mask).Equal(ts.From)) {
			newcursor = ts.From
			*cursor = newcursor
			return newcursor
		}

		// move the cursor after
		newcursor = newcursor.Truncate(mask).Add(mask)

		// check over end boundaries
		if newcursor.After(ts.To) {
			if fBoundaries && !cursor.Equal(ts.To) {
				newcursor = ts.To
			} else {
				newcursor = time.Time{}
			}
		}
	}

	*cursor = newcursor
	return newcursor
}
