package configs

// Row interval (seconds): starts at 6s, floors at 1.5s, decreases by 0.1s per 10s elapsed.
const (
	InitialRowInterval   = 6.0
	MinRowInterval       = 1.5
	RowIntervalDecRate   = 0.1
	RowIntervalDecPeriod = 10.0
)

// RowInterval returns the current row spawn interval in seconds for the given elapsed time.
func RowInterval(elapsedSeconds float64) float64 {
	dec := (elapsedSeconds / RowIntervalDecPeriod) * RowIntervalDecRate
	interval := InitialRowInterval - dec
	if interval < MinRowInterval {
		return MinRowInterval
	}
	return interval
}

// Fall speed (rows per second): starts at 0.5, caps at 2.5, increases by 0.15 per 30s elapsed.
const (
	InitialFallSpeed   = 0.5
	MaxFallSpeed       = 2.5
	FallSpeedIncRate   = 0.15
	FallSpeedIncPeriod = 30.0
)

// FallSpeed returns the current mask fall speed in rows per second for the given elapsed time.
func FallSpeed(elapsedSeconds float64) float32 {
	inc := (elapsedSeconds / FallSpeedIncPeriod) * float64(FallSpeedIncRate)
	speed := InitialFallSpeed + inc
	if speed > MaxFallSpeed {
		return float32(MaxFallSpeed)
	}
	return float32(speed)
}
