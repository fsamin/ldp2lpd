package main

import (
	"context"
	"log"
	"sort"
	"sync"
	"time"
)

func count(m message) {
	level := getLevel(int(m["level"].(float64)))
	ts := int64(m["timestamp"].(float64))
	globalCounters.add(time.Unix(ts, 0), level)
}

var globalCounters = &counters{}

type counter struct {
	infos  int
	warns  int
	errors int
}

type counters struct {
	mutex   sync.RWMutex
	counter [9]counter
}

func (c *counters) maxValue() (max int) {
	vals := []int{}
	for i := range c.counter {
		vals = append(vals, c.counter[i].value())
	}
	sort.Ints(vals)
	return vals[len(vals)-1]
}

func (c *counter) value() int {
	if c == nil {
		return 0
	}
	return c.infos + c.warns + c.errors
}

func rolling(ctx context.Context, c *counters) {
	tick := time.NewTicker(timeUnit)
	defer tick.Stop()

	for {
		select {
		case <-tick.C:
			c.mutex.Lock()
			for i := len(c.counter) - 1; i >= 1; i-- {
				if verbose {
					log.Printf("swapping %d with %d", i, (i - 1))
				}
				c.counter[i] = c.counter[i-1]
			}
			c.counter[0] = counter{}
			c.mutex.Unlock()
		case <-ctx.Done():
			return
		}
	}
}

func (c *counters) add(t time.Time, level level) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	d := time.Since(t)
	for i := 0; i < len(c.counter); i++ {
		cntr := &c.counter[i]
		if d >= time.Duration(i)*timeUnit && d < time.Duration(i+1)*timeUnit {
			cntr.add(level)
			break
		}
	}
}

func (c *counter) add(l level) {
	switch l {
	case debugLevel, infoLevel:
		c.infos++
	case warnLevel, errorLevel:
		c.warns++
	case fatalLevel, panicLevel:
		c.errors++
	}
}
