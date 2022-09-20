package timeslice

import (
	"fmt"
	"time"

	"github.com/sunraylab/timeline/duration"
)

func ExampleTimeSlice_MoveFrom_one() {
	// take a date and build a timeslice staring at this date and ending 7 days after
	from := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	ts := MakeTimeslice(from, duration.Week)
	fmt.Println(ts)

	// Move forward the begining by 4 days
	ts.MoveFrom(ts.From.Add(duration.Day*4), false)
	fmt.Println(ts)

	// Move fotrward again the begining by 4 days
	ts.MoveFrom(ts.From.Add(duration.Day*4), false)
	fmt.Println(ts)

	// Output:
	// { 20220101 UTC - 20220108 UTC : 7d }
	// { 20220105 UTC - 20220108 UTC : 3d }
	// { 20220109 UTC - 20220109 UTC : 0 }
}

func ExampleTimeSlice_MoveTo_one() {
	// take a date and build a timeslice staring at this date and ending 7 days after
	from := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	ts := MakeTimeslice(from, duration.Week)
	fmt.Println(ts)

	// Move backward the ending by 4 days
	ts.MoveTo(ts.To.Add(-duration.Day*4), false)
	fmt.Println(ts)

	// Move backward again the ending by 4 days
	ts.MoveTo(ts.To.Add(-duration.Day*4), false)
	fmt.Println(ts)

	// Output:
	// { 20220101 UTC - 20220108 UTC : 7d }
	// { 20220101 UTC - 20220104 UTC : 3d }
	// { 20211231 UTC - 20211231 UTC : 0 }
}

func ExampleTimeSlice_MoveFrom_two() {
	// take a date and build a timeslice staring at this date and ending 7 days after
	from := time.Date(2022, 1, 6, 8, 0, 0, 0, time.UTC)
	ts := MakeTimeslice(from, duration.Day)
	fmt.Println(ts)
	ts.MoveFrom(time.Date(2022, 1, 6, 9, 0, 0, 0, time.UTC), true)
	fmt.Println(ts)
	ts.MoveFrom(time.Date(2022, 1, 7, 9, 0, 0, 0, time.UTC), true)
	fmt.Println(ts)
	// Output:
	// { 20220106 08:00:00 UTC - 20220107 08:00:00 UTC : 1d }
	// { 20220106 09:00:00 UTC - 20220107 08:00:00 UTC : 23h }
	// { 20220107 08:00:00 UTC - 08:00:00 : 0 }
}

func ExampleTimeSlice_MoveTo_two() {
	// take a date and build a timeslice staring at this date and ending 7 days after
	from := time.Date(2022, 1, 6, 8, 0, 0, 0, time.UTC)
	ts := MakeTimeslice(from, duration.Day)
	fmt.Println(ts)
	ts.MoveTo(time.Date(2022, 1, 6, 9, 0, 0, 0, time.UTC), true)
	fmt.Println(ts)
	ts.MoveTo(time.Date(2022, 1, 6, 7, 0, 0, 0, time.UTC), true)
	fmt.Println(ts)
	// Output:
	// { 20220106 08:00:00 UTC - 20220107 08:00:00 UTC : 1d }
	// { 20220106 08:00:00 UTC - 09:00:00 : 1h }
	// { 20220106 08:00:00 UTC - 08:00:00 : 0 }
}

func ExampleTimeSlice_String() {
	ts := MakeTimeslice(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC), duration.Week+time.Hour*31)
	fmt.Println(ts)
	// Output: { 20220101 UTC - 20220109 07:00:00 UTC : 8d7h }
}

func ExampleTimeSlice_Progress() {
	// take a 24 hours timeslice, starting the 2022,1,1 at midnight
	ts := MakeTimeslice(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC), duration.Day)

	// get the corresponding progress date the same day at 6 AM
	rate := ts.Progress(time.Date(2022, 1, 1, 6, 0, 0, 0, time.UTC))
	fmt.Println(rate)
	// Output: 0.25
}

func ExampleTimeSlice_WhatTime() {
	// take a 10 days timeslice, starting the 2022,1,1 at 8AM
	ts := MakeTimeslice(time.Date(2022, 1, 1, 8, 0, 0, 0, time.UTC), duration.Day*10)
	fmt.Println(ts)

	for rate := 0.0; rate <= 1.0; rate += 0.2 {
		t := ts.WhatTime(rate)
		fmt.Println(t.Format("20060102 15:04:05 MST"))
	}

	// Output:
	// { 20220101 08:00:00 UTC - 20220111 08:00:00 UTC : 10d }
	// 20220101 08:00:00 UTC
	// 20220103 08:00:00 UTC
	// 20220105 08:00:00 UTC
	// 20220107 08:00:00 UTC
	// 20220109 08:00:00 UTC
	// 20220111 08:00:00 UTC
}

func ExampleTimeSlice() {

	// take a 3 days timeslice, starting the 2022,1,6 at 7:30AM
	ts := MakeTimeslice(time.Date(2022, 1, 6, 7, 30, 0, 0, time.UTC), duration.Day*3)
	fmt.Printf("A timeslice: %s\n", ts)

	// get a scan mask to handle 10 steps max
	mask := ts.GetScanMask(10)
	fmt.Printf("scan mask = %s\n", mask.String())

	// scan to build a grid with dates matching the mask inside this time slice, includes boundaries any time
	var xgridtime time.Time
	for ts.Scan(&xgridtime, mask, true); !xgridtime.IsZero(); ts.Scan(&xgridtime, mask, true) {
		progress := ts.Progress(xgridtime)
		fmt.Printf("%s ==> progress: %3.1f%%\n", xgridtime.Format("20060102 15:04:05"), progress*100)
	}

	// What is the time at the middle of this timeslice ?
	middle := ts.WhatTime(0.5)
	fmt.Printf("the middle of this timeslice is: %v\n", middle)

	// Apply a mask to get the Quarter corresponding to this date
	quarter, _ := MASK_QUARTER.Apply(middle)
	fmt.Printf("the corresponding quarter starts: %v\n", quarter)

	// Output:
	// A timeslice: { 20220106 07:30:00 UTC - 20220109 07:30:00 UTC : 3d }
	// scan mask = half-day
	// 20220106 07:30:00 ==> progress: 0.0%
	// 20220106 12:00:00 ==> progress: 6.2%
	// 20220107 00:00:00 ==> progress: 22.9%
	// 20220107 12:00:00 ==> progress: 39.6%
	// 20220108 00:00:00 ==> progress: 56.2%
	// 20220108 12:00:00 ==> progress: 72.9%
	// 20220109 00:00:00 ==> progress: 89.6%
	// 20220109 07:30:00 ==> progress: 100.0%
	// the middle of this timeslice is: 2022-01-07 19:30:00 +0000 UTC
	// the corresponding quarter starts: 2022-01-01 00:00:00 +0000 UTC
}

func ExampleTimeSlice_GetScanMask() {

	ts := MakeTimeslice(time.Date(2008, 10, 31, 21, 0, 0, 0, time.UTC), duration.Month*3)

	for i := 10; i > 0; i-- {
		mask := ts.GetScanMask(12)
		fmt.Printf("best scan mask:%12s <== Timeslice: %s\n", mask.String(), ts)

		ts.ExtendTo(-ts.Duration().Adjust(0.7))
	}

	// Output:
	// best scan mask:       month <== Timeslice: { 20081031 21:00:00 UTC - 20090131 04:30:00 UTC : 3M }
	// best scan mask:       month <== Timeslice: { 20081031 21:00:00 UTC - 20081128 06:27:00 UTC : 27d9h27m }
	// best scan mask:         day <== Timeslice: { 20081031 21:00:00 UTC - 20081109 02:14:06 UTC : 8d5h14m~ }
	// best scan mask:    half-day <== Timeslice: { 20081031 21:00:00 UTC - 20081103 08:10:13 UTC : 2d11h10m~ }
	// best scan mask:     4 hours <== Timeslice: { 20081031 21:00:00 UTC - 20081101 14:45:04 UTC : 17h45m4s }
	// best scan mask:   half-hour <== Timeslice: { 20081031 21:00:00 UTC - 20081101 02:19:31 UTC : 5h19m31s }
	// best scan mask:  15 minutes <== Timeslice: { 20081031 21:00:00 UTC - 22:35:51 : 1h35m51s }
	// best scan mask:  15 minutes <== Timeslice: { 20081031 21:00:00 UTC - 21:28:45 : 28m45s }
	// best scan mask:      minute <== Timeslice: { 20081031 21:00:00 UTC - 21:08:37 : 8m37s }
	// best scan mask:      minute <== Timeslice: { 20081031 21:00:00 UTC - 21:02:35 : 2m35s }
}

func ExampleTimeMask_GetTimeFormat() {

	// according to a choosen date
	t1 := time.Date(2008, 10, 30, 21, 12, 59, 0, time.UTC)
	fmt.Printf("Choosen time t1=%s\n", t1.Format("2006-01-02 15:04:05"))

	// format this date according to the mask
	for mask := MASK_min; mask <= MASK_max; mask++ {
		strfmt := mask.GetTimeFormat(t1, t1)
		strt := t1.Format(strfmt)
		fmt.Printf("with mask:%12s, renders: %s\n", mask.String(), strt)
	}

	// Now renders the same time, with a hour level mask, but comparing with another time
	// GetTimeFormat decides if another time component needs to be printed to make
	// the output more comprehensive.
	// Usefull if you scan times thru a timeline and want to streamline the output
	t2 := t1.Add(1 * time.Hour * 24 * 31)
	fmt.Printf("Next time t2=%s\n", t2.Format("2006-01-02 15:04:05"))

	fmt.Printf("Streamlined output for t2 renders: %s\n", t2.Format(MASK_HOUR.GetTimeFormat(t2, t1)))

	// Output:
	// Choosen time t1=2008-10-30 21:12:59
	// with mask:      minute, renders: 21:12
	// with mask:  15 minutes, renders: 21:12
	// with mask:   half-hour, renders: 21:12
	// with mask:        hour, renders: 21:12
	// with mask:     4 hours, renders: 21:12
	// with mask:    half-day, renders: Thu, 30 21:12
	// with mask:         day, renders: Thu, 30
	// with mask:       month, renders: Oct
	// with mask:     quarter, renders: Oct
	// with mask:        year, renders: 2008
	// Next time t2=2008-11-30 21:12:59
	// Streamlined output for t2 renders: Sun, Nov 30 21:12
}