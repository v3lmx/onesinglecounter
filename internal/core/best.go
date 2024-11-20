package core

import "time"

type best struct {
	Hour    uint
	Day     uint
	Week    uint
	Month   uint
	AllTime uint
}

func Best(count <-chan uint) {
	b := best{}
	for {
		select {
		case c := <-count:
			if c <= b.Hour {
				continue
			}
			b.Hour = c
			if c <= b.Day {
				continue
			}
			b.Day = c
			if c <= b.Week {
				continue
			}
			b.Week = c
			if c <= b.Month {
				continue
			}
			b.Month = c
			if c <= b.AllTime {
				continue
			}
			b.AllTime = c
		// todo find cron style method to reset time
		case <-time.After(time.Second * 2):
			b.Hour = 0
		}
	}
}
