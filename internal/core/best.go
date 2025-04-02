package core

import (
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/charmbracelet/log"
	"github.com/robfig/cron/v3"
)

type CurrentBest struct {
	sync.RWMutex
	Minute  uint64
	Hour    uint64
	Day     uint64
	Week    uint64
	Month   uint64
	Year    uint64
	AllTime uint64
}

func (b *CurrentBest) Copy() CurrentBest {
	return CurrentBest{
		AllTime: b.AllTime,
		Year:    b.Year,
		Month:   b.Month,
		Week:    b.Week,
		Day:     b.Day,
		Hour:    b.Hour,
		Minute:  b.Minute,
	}
}

func (best *CurrentBest) Format() string {
	var sb strings.Builder

	sb.WriteString("alltime:")
	sb.WriteString(strconv.Itoa(int(best.AllTime)))
	sb.WriteString(":year:")
	sb.WriteString(strconv.Itoa(int(best.Year)))
	sb.WriteString(":month:")
	sb.WriteString(strconv.Itoa(int(best.Month)))
	sb.WriteString(":week:")
	sb.WriteString(strconv.Itoa(int(best.Week)))
	sb.WriteString(":day:")
	sb.WriteString(strconv.Itoa(int(best.Day)))
	sb.WriteString(":hour:")
	sb.WriteString(strconv.Itoa(int(best.Hour)))
	sb.WriteString(":minute:")
	sb.WriteString(strconv.Itoa(int(best.Minute)))

	return sb.String()
}

func newBest() CurrentBest {
	return CurrentBest{}
}

func Best(count *atomic.Uint64, best *CurrentBest, tickClock *sync.Cond, bestClock *sync.Cond) {
	nextMinute := make(chan struct{})
	nextHour := make(chan struct{})
	nextDay := make(chan struct{})
	nextWeek := make(chan struct{})
	nextMonth := make(chan struct{})
	nextYear := make(chan struct{})
	// backup := make(chan struct{})

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

		// wait for next minute
		// bestClock.L.Lock()
		// bestClock.Wait()
		//
		// best.RLock()
		// b := best.Copy()
		// best.RUnlock()
		//
		// bestClock.L.Unlock()

		// wait for tick to update count
		tickClock.L.Lock()
		tickClock.Wait()

		c := count.Load()

		tickClock.L.Unlock()

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

		select {
		case <-nextMinute:
			b.Minute = c
			bestClock.Broadcast()
		case <-nextHour:
			b.Hour = c
		case <-nextDay:
			b.Day = c
		case <-nextWeek:
			b.Week = c
		case <-nextMonth:
			b.Month = c
		case <-nextYear:
			b.Year = c
		}
		log.Debugf("current best: %s", b.Format())
	}
}
