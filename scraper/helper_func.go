package main

import (
	"math/rand"
	"time"
)

func RandDuration(low, high float64) time.Duration {
	return time.Duration((rand.Float64()*(high-low) + low) * float64(time.Second))
}
