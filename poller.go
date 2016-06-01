package importer

import "time"

// Poller polls an importer.
type Poller struct {
	i Importer
	w Worker
}

// NewPoller returns a Poller to poll the given importer.
func NewPoller(i Importer, w Worker) *Poller {
	return &Poller{i, w}
}

// Start starts the poller polling in a separate goroutine.
func (p *Poller) Start(interval time.Duration) {
	go func() {
		for {
			p.i.Poll(p.w)
			time.Sleep(interval)
		}
	}()
}
