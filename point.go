package timesies

import "fmt"

// A point represents a single data record.
type point struct {
	count uint32 // number of events
}

func (p *point) getCount() uint32 {
	if p == nil {
		return 0
	}

	return p.count
}

func (p *point) String() string {
	return fmt.Sprintf("%d", p.count)
}
