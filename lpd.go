package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/rakyll/launchpad"
)

func display(ctx context.Context, c *counters) {
	pad, err := launchpad.Open()
	if err != nil {
		log.Fatalf("Error initializing launchpad:", err)
	}
	defer pad.Close()

	// turn off all of the lights
	_ = pad.Clear()

	tick := time.NewTicker(3 * time.Second)
	defer tick.Stop()

	for {
		select {
		case <-tick.C:
			displayCounters(c, pad)
		case <-ctx.Done():
			return
		}
	}

}

func displayCounters(c *counters, pad *launchpad.Launchpad) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	maxValue := c.maxValue()
	if pad != nil {
		_ = pad.Clear()
	}

	var padIndex int
	for i := range c.counter {
		cntr := &c.counter[i]
		if err := displayCounter(cntr, pad, 8-padIndex, maxValue); err != nil && verbose {
			log.Println(err)
		}
		padIndex++
	}
}

func displayCounter(c *counter, pad *launchpad.Launchpad, index, maxValue int) error {
	if c == nil || maxValue == 0 || c.value() == 0 {
		return fmt.Errorf("no data on %d", index)
	}

	fInfo := float64(c.infos) / float64(maxValue)
	fWarn := float64(c.warns) / float64(maxValue)
	fError := float64(c.errors) / float64(maxValue)

	vInfo := int(math.Round(fInfo * 8))
	vWarn := int(math.Round(fWarn * 8))
	vError := int(math.Round(fError * 8))

	//Hilight error
	if c.errors > 0 && index < 8 {
		_ = pad.Light(index, 8, 0, 3)
	}

	//Display info
	if verbose {
		log.Println(index, c, "Info: ", vInfo, "Warn: ", vWarn, "Error: ", vError)
	}

	var y = 0
	for i := 0; i < vInfo; i++ {
		_ = pad.Light(index, 7-y, 3, 0)
		y++
	}

	//Display error
	for i := 0; i < vWarn; i++ {
		_ = pad.Light(index, 7-y, 3, 3)
		y++
	}

	//Display fatal
	for i := 0; i < vError; i++ {
		_ = pad.Light(index, 7-y, 0, 3)
		y++
	}

	return nil
}
