package game

import "time"

const INTENSITY_PER_MINUTE = .005555556
const MINIMUM_LIGHT_INTENSITY float32 = .11111112

type WorldClock struct {
	Current time.Time

	timeLastUpdate time.Time
}

func NewWorldClock() *WorldClock {
	return &WorldClock{
		Current: time.Date(1996, 01, 01, 00, 00, 00, 00, time.Local),
	}
}

func (w *WorldClock) Update() {
	if w.Current.Hour() == 0 {
		w.Current = w.Current.Add(time.Hour * 1)
	}
	// if time.Since(w.timeLastUpdate).Minutes() > 1 {
	w.Current = w.Current.Add(time.Second * 10)
	w.timeLastUpdate = time.Now()
	// }
}

func (w *WorldClock) IsDaylight() bool {
	return w.Current.After(
		time.Date(1996, 01, 01, 7, 0, 0, 0, time.Local),
	) && w.Current.Before(
		time.Date(1996, 01, 01, 18, 0, 0, 0, time.Local),
	)
}
