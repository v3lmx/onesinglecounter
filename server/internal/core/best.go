package core

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/robfig/cron/v3"
)

type CurrentBest struct {
	sync.RWMutex
	Best Best
}

type Best struct {
	Minute  uint64
	Hour    uint64
	Day     uint64
	Week    uint64
	Month   uint64
	Year    uint64
	AllTime uint64
}

func (b *CurrentBest) Copy() Best {
	b.RLock()
	defer b.RUnlock()
	return Best{
		AllTime: b.Best.AllTime,
		Year:    b.Best.Year,
		Month:   b.Best.Month,
		Week:    b.Best.Week,
		Day:     b.Best.Day,
		Hour:    b.Best.Hour,
		Minute:  b.Best.Minute,
	}
}

func parseElement(elements []string, index int, key string) (uint64, error) {
	if elements[index] != key {
		return 0, fmt.Errorf("Could not parse best: %s key", key)
	}
	value, err := strconv.Atoi(elements[index+1])
	if err != nil {
		return 0, fmt.Errorf("Could not parse best: %s value", key)
	}
	return uint64(value), nil
}

func ParseBest(s string) (Best, error) {
	best := Best{}
	elements := strings.Split(s, ":")
	fmt.Printf("elements: %v\n", elements)

	alltime, err := parseElement(elements, 0, "alltime")
	if err != nil {
		return Best{}, err
	}
	best.AllTime = uint64(alltime)

	year, err := parseElement(elements, 2, "year")
	if err != nil {
		return Best{}, err
	}
	best.Year = uint64(year)

	month, err := parseElement(elements, 4, "month")
	if err != nil {
		return Best{}, err
	}
	best.Month = uint64(month)

	week, err := parseElement(elements, 6, "week")
	if err != nil {
		return Best{}, err
	}
	best.Week = uint64(week)

	day, err := parseElement(elements, 8, "day")
	if err != nil {
		return Best{}, err
	}
	best.Day = uint64(day)

	hour, err := parseElement(elements, 10, "hour")
	if err != nil {
		return Best{}, err
	}
	best.Hour = uint64(hour)

	minute, err := parseElement(elements, 12, "minute")
	if err != nil {
		return Best{}, err
	}
	best.Minute = uint64(minute)

	return best, nil
}

func (best Best) Format() string {
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

func (best *CurrentBest) Format() string {
	best.RLock()
	defer best.RUnlock()

	return best.Best.Format()
}

func BestLoop(count *atomic.Uint64, best *CurrentBest, tickBroadcast *Cond, bestBroadcast *Cond, bestTickTime time.Duration, backup Backup) {
	t := time.NewTicker(bestTickTime)

	defer t.Stop()

	bestChan := make(chan Best, 1)
	bestChan <- best.Copy()
	go func() {
		for range t.C {
			bestBroadcast.Broadcast()
			best := <-bestChan
			bestChan <- best
			err := backup.Backup(count.Load(), best)
			if err != nil {
				slog.Error("Could not backup", "error_msg", err)
			}
		}
	}()

	nextMinute := make(chan struct{})
	nextHour := make(chan struct{})
	nextDay := make(chan struct{})
	nextWeek := make(chan struct{})
	nextMonth := make(chan struct{})
	nextYear := make(chan struct{})

	go func() {
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

		for {
			select {
			case <-nextMinute:
				best.Lock()
				best.Best.Minute = count.Load()
				best.Unlock()
			case <-nextHour:
				best.Lock()
				best.Best.Hour = count.Load()
				best.Unlock()
			case <-nextDay:
				best.Lock()
				best.Best.Day = count.Load()
				best.Unlock()
			case <-nextWeek:
				best.Lock()
				best.Best.Week = count.Load()
				best.Unlock()
			case <-nextMonth:
				best.Lock()
				best.Best.Month = count.Load()
				best.Unlock()
			case <-nextYear:
				best.Lock()
				best.Best.Year = count.Load()
				best.Unlock()
			}
		}
	}()

	for {
		<-bestChan
		bestChan <- best.Copy()

		// wait for tick to update count
		tickBroadcast.L.Lock()
		tickBroadcast.Wait()

		c := count.Load()

		tickBroadcast.L.Unlock()

		best.Lock()
		if c <= best.Best.Minute {
			best.Unlock()
			continue
		}
		best.Best.Minute = c
		if c <= best.Best.Hour {
			best.Unlock()
			continue
		}
		best.Best.Hour = c
		if c <= best.Best.Day {
			best.Unlock()
			continue
		}
		best.Best.Day = c
		if c <= best.Best.Week {
			best.Unlock()
			continue
		}
		best.Best.Week = c
		if c <= best.Best.Month {
			best.Unlock()
			continue
		}
		best.Best.Month = c
		if c <= best.Best.Year {
			best.Unlock()
			continue
		}
		best.Best.Year = c
		if c <= best.Best.AllTime {
			best.Unlock()
			continue
		}
		best.Best.AllTime = c
		best.Unlock()
	}
}
