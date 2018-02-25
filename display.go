package main

import (
	"fmt"
	"io"
	"sort"
	"time"
)

func printText(ss *StationSet, w io.Writer) {

	now := time.Now()
	nowUnix := now.Unix()

	fmt.Fprintf(w, "As of %s\n\n", now.Format("15:04:05"))

	//fmt.Fprintf(w, "It is %d\n\n", now)
	_ = `
<html><head><style>
body { font-size: 1.5em; }
</style></head><body>
`
	//	fmt.Fprintln(w, header)

	// TODO sort or weight
	for _, s := range ss.Stations {
		fmt.Fprintln(w, s.Name)

		for line, times := range s.NorthTimes {
			fmt.Fprintf(w, "  N  %s ", line)
			// TODO sort
			for i, t := range unixToRelative(nowUnix, times) {
				fmt.Fprintf(w, " %d", t)
				if i > 8 {
					break
				}
			}
			fmt.Fprintln(w)
		}
		for line, times := range s.SouthTimes {
			fmt.Fprintf(w, "  S  %s ", line)
			for i, t := range unixToRelative(nowUnix, times) {
				fmt.Fprintf(w, " %d", t)
				if i > 8 {
					break
				}
			}
			fmt.Fprintln(w)
		}
		fmt.Fprintln(w)
	}
	//	fmt.Fprintln(w, "</body></html>")
}

func unixToRelative(now int64, ts []int64) []int {
	// TODO make this smaller?
	var rs []int

	for _, t := range ts {
		if t > now {
			// TODO round
			rs = append(rs, int((t-now)/60))
		}
	}

	sort.Ints(rs)
	return rs
}
