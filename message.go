package main

type message map[string]interface{}

type level uint8

const (
	panicLevel level = iota
	fatalLevel
	errorLevel
	warnLevel
	infoLevel
	debugLevel
)

func getLevel(v int) level {
	return level(v - 2)
}
