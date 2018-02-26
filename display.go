package main

import (
	"fmt"
	"io"
	"sort"
	"time"
)

func printText(ss *StationSet, w io.Writer) {

	loc, _ := time.LoadLocation("America/New_York")
	now := time.Now()
	nowUnix := now.Unix()

	header := `
<html><head>
<meta name="viewport" content="width=device-width, initial-scale=1">
<style>
body {
  font-size: 1.3em;
  font-family: monospace;
}

#time { margin: 0.5em; }

.stop {
  margin: 0.8em 0 0.2em 0.3em;
  font-weight: 500;
}

.row span {
  margin: 0.2em;
  display: table-cell;
  min-width: 1.8em;
  width: 1.8em;
}

.line { font-weight: 600; }
.arr { padding-right: .2em; }
</style></head><body>
`
	fmt.Fprintln(w, header)
	fmt.Fprintf(w, "<div id='time'>As of %s</div>\n", now.In(loc).Format("3:04:05"))

	// TODO sort or weight
	for _, s := range ss.Stations {
		fmt.Fprintf(w, "<div class='stop'>%s</div>\n", s.Name)

		for line, times := range s.NorthTimes {
			fmt.Fprintf(w, "<div class='row'><span class='dir'>N</span><span class='line'>%s</span>", line)
			// TODO sort
			for i, t := range unixToRelative(nowUnix, times) {
				fmt.Fprintf(w, "<span class='arr'>%d</span>", t)
				if i > 8 {
					break
				}
			}
			fmt.Fprintln(w, "</div>\n")

		}
		for line, times := range s.SouthTimes {
			fmt.Fprintf(w, "<div class='row'><span class='dir'>S</span><span class='line'>%s</span>", line)
			for i, t := range unixToRelative(nowUnix, times) {
				fmt.Fprintf(w, "<span class='arr'>%d</span>", t)
				if i > 8 {
					break
				}
			}
			fmt.Fprintln(w, "</div>\n")
		}
	}
	fmt.Fprintln(w, "</body></html>")
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
