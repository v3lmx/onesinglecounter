package core

import (
	"encoding/binary"
	"slices"

	"github.com/charmbracelet/log"
	"github.com/robfig/cron/v3"
)

type CurrentBest struct {
	Minute  uint64
	Hour    uint64
	Day     uint64
	Week    uint64
	Month   uint64
	Year    uint64
	AllTime uint64
}

const (
	tagMinute  byte = 10
	tagHour    byte = 11
	tagDay     byte = 12
	tagWeek    byte = 13
	tagMonth   byte = 14
	tagYear    byte = 15
	tagAllTime byte = 16
)

func formatTimestamp(tag byte, value uint64) []byte {
	buf := make([]byte, binary.MaxVarintLen64)
	l := binary.PutUvarint(buf, value)
	m := Message{Tag: tag, Value: buf[:l]}
	return m.TLV()
}

func (best CurrentBest) Format() (Message, error) {
	tag := ServerMask & MessageTypeGetCurrent

	value := slices.Concat(
		formatTimestamp(tagMinute, best.Minute),
		formatTimestamp(tagHour, best.Hour),
		formatTimestamp(tagDay, best.Day),
		formatTimestamp(tagWeek, best.Week),
		formatTimestamp(tagMonth, best.Month),
		formatTimestamp(tagYear, best.Year),
		formatTimestamp(tagAllTime, best.AllTime),
	)

	return Message{Tag: tag, Value: value}, nil
}

func newBest() CurrentBest {
	return CurrentBest{}
}

func Best(count <-chan uint64, request <-chan struct{}, broadcast chan<- CurrentBest) {
	nextMinute := make(chan struct{})
	nextHour := make(chan struct{})
	nextDay := make(chan struct{})
	nextWeek := make(chan struct{})
	nextMonth := make(chan struct{})
	nextYear := make(chan struct{})
	// backup := make(chan struct{})
	lastCount := uint64(0)

	c := cron.New()
	_, err := c.AddFunc("* * * * *", func() { nextMinute <- struct{}{} })
	if err != nil {
		panic("Couldn't start cron")
	}
	_, err = c.AddFunc("@hourly", func() { nextHour <- struct{}{} })
	if err != nil {
		panic("Couldn't start cron")
	}
	_, err = c.AddFunc("@daily", func() { nextDay <- struct{}{} })
	if err != nil {
		panic("Couldn't start cron")
	}
	// @weekly starts sunday instead of monday
	_, err = c.AddFunc("0 0 * * 1", func() { nextWeek <- struct{}{} })
	if err != nil {
		panic("Couldn't start cron")
	}
	_, err = c.AddFunc("@monthly", func() { nextMonth <- struct{}{} })
	if err != nil {
		panic("Couldn't start cron")
	}
	_, err = c.AddFunc("@yearly", func() { nextYear <- struct{}{} })
	if err != nil {
		panic("Couldn't start cron")
	}
	c.Start()

	b := newBest()
	for {
		select {
		// case <-backup:
		// 	//backup
		case <-request:
			broadcast <- b
		case c := <-count:
			lastCount = c
			if c <= b.Minute {
				continue
			}
			b.Minute = c
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
			if c <= b.Year {
				continue
			}
			b.Year = c
			if c <= b.AllTime {
				continue
			}
			b.AllTime = c
		case <-nextMinute:
			b.Minute = lastCount
			// causes deadlock
			// broadcast <- b
		case <-nextHour:
			b.Hour = lastCount
		case <-nextDay:
			b.Day = lastCount
		case <-nextWeek:
			b.Week = lastCount
		case <-nextMonth:
			b.Month = lastCount
		case <-nextYear:
			b.Year = lastCount
		}
		log.Debugf("current best: %v", b)
	}
}
