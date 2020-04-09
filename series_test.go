package timesies

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSeries_end(t *testing.T) {
	s := &series{
		start:  100, // start is at 100
		unit:   10,  // each point is 10 seconds long
		length: 5,   // end is at 1 + 5 * 10 = 150
	}

	assert.Equal(t, 150, int(s.end()))
}

func TestSeries_count(t *testing.T) {
	testseries := &series{
		start:  100, // start is at 100
		unit:   10,  // each point is 10 seconds long
		length: 5,   // end is at 1 + 5 * 10 = 150
		points: []*point{
			&point{count: 1}, // [100, 110)
			&point{count: 2}, // [110, 120)
			&point{count: 3}, // [120, 130)
			nil,              // [130, 140)
			&point{count: 5}, // [140, 150)
		},
	}

	cases := map[string]struct {
		series      *series
		from, until uint32
		count       uint32
	}{
		"series is nil": {
			from:  1,
			until: 42,
			count: 0,
		},
		"empty series": {
			series: &series{},
			from:   1,
			until:  42,
			count:  0,
		},
		"before series": {
			series: testseries,
			from:   0,
			until:  50,
			count:  0,
		},
		"until overlaps with start": {
			series: testseries,
			from:   0,
			until:  100,
			count:  0,
		},
		"after": {
			series: testseries,
			from:   200,
			until:  250,
			count:  0,
		},
		"from overlaps with end": {
			series: testseries,
			from:   150,
			until:  200,
			count:  0,
		},
		"range overlaps with start": {
			series: testseries,
			from:   42,
			until:  122,
			count:  1 + 2,
		},
		"range overlaps with end": {
			series: testseries,
			from:   125,
			until:  165,
			count:  3 + 5,
		},
		"range within series": {
			series: testseries,
			from:   113,
			until:  137,
			count:  2 + 3,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			actual := c.series.count(c.from, c.until)
			assert.Equal(t, c.count, actual)
		})
	}
}
