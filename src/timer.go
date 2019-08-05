package chip8

import (
	"sync"
	"time"
)

//Timer decrements the timer value at a rate of 60Hz.
type Timer struct {
	timer  uint
	mux    *sync.RWMutex
	ticker *time.Ticker
}

//NewTimer returns a pointer to a timer struct.
//Sets the ticker to 60Hz and initialises the timer at 0.
//It also kicks off the ticker goroutine.
func NewTimer() *Timer {
	ticker := time.NewTicker(16666666 * time.Nanosecond)
	t := &Timer{
		timer:  0,
		mux:    &sync.RWMutex{},
		ticker: ticker,
	}
	go t.Tick()
	return t
}

//Tick is a goroutine that will decrement 1 from the timer at a rate of 60Hz.
func (t *Timer) Tick() {
	for range t.ticker.C {
		if t.timer != 0 {
			t.mux.Lock()
			t.timer--
			t.mux.Unlock()
		}
	}
}

//Set the timer to a new value.
func (t *Timer) Set(n uint) error {
	t.mux.Lock()
	defer t.mux.Unlock()
	t.timer = n
	return nil
}

//Get the current timer value.
func (t *Timer) Get() (uint, error) {
	t.mux.RLock()
	defer t.mux.RUnlock()
	v := t.timer
	return v, nil
}

//Stop the ticker.
func (t *Timer) Stop() {
	t.ticker.Stop()
}
