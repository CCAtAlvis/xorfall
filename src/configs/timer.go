package configs

import rl "github.com/gen2brain/raylib-go/raylib"

type GameTimeManager struct {
	StartTime float64
	RealDelta float32
	Delta     float32
	TimeScale float32
}

var GameTime = &GameTimeManager{
	StartTime: rl.GetTime(),
	TimeScale: 1.0,
}

func (g *GameTimeManager) Update() {
	g.RealDelta = rl.GetFrameTime()
	g.Delta = g.RealDelta * g.TimeScale
}

type Timer struct {
	duration        float64
	elapsedDuration float64
	callback        func()
}

var timers = make([]*Timer, 0)

func GetAllTimers() []*Timer {
	return timers
}

func NewTimer(duration float64, callback func()) *Timer {
	timer := &Timer{
		duration: duration,
	}
	if callback != nil {
		timer.callback = callback
	}
	timers = append(timers, timer)
	return timer
}

func (t *Timer) Update(dt float64) {
	t.elapsedDuration += dt
	if t.IsDone() {
		t.callback()
	}
}

func (t *Timer) IsDone() bool {
	return t.elapsedDuration >= t.duration
}
