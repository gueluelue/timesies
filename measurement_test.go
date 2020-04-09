package timesies

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func randomPoints(len int) []*point {
	points := make([]*point, len)

	r := rand.New(rand.NewSource(99))
	for i := 0; i < len; i++ {
		if r.Float32() > 0.5 {
			continue // generate some nil points
		}

		points[i] = &point{count: uint32(r.Int31n(10))}
	}

	return points
}

func TestMeasurement_Count(t *testing.T) {
	seconds := &series{
		start:  120,
		unit:   1,
		length: 60,
		points: randomPoints(60),
	}

	minutes := &series{
		start:  120,
		unit:   60,
		length: 60,
		points: randomPoints(60),
	}
	minutes.points[0] = &point{count: seconds.count(seconds.start, seconds.end())}

	hours := &series{
		start:  120,
		unit:   3600,
		length: 24 * 7, // one week
		points: randomPoints(24 * 7),
	}
	hours.points[0] = &point{count: minutes.count(minutes.start, minutes.end())}

	testmeasurement := &Measurement{
		series: []*series{seconds, minutes, hours},
	}

	cases := map[string]struct {
		measurement *Measurement
		from, until uint32
		count       uint32
	}{
		"measurement is nil": {
			from:  1,
			until: 42,
			count: 0,
		},
		"empty measurement": {
			measurement: &Measurement{},
			from:        1,
			until:       42,
			count:       0,
		},
		"before measurement": {
			measurement: testmeasurement,
			from:        0,
			until:       50,
			count:       0,
		},
		"until overlaps with start": {
			measurement: testmeasurement,
			from:        0,
			until:       seconds.start,
			count:       0,
		},
		"after": {
			measurement: testmeasurement,
			from:        hours.end() + 1,
			until:       hours.end() + 50,
			count:       0,
		},
		"within seconds": {
			measurement: testmeasurement,
			from:        seconds.start + 10,
			until:       seconds.end() - 10,
			count:       seconds.count(seconds.start+10, seconds.end()-10),
		},
		"within minutes": {
			measurement: testmeasurement,
			from:        minutes.start + 10,
			until:       minutes.end() - 10,
			count:       minutes.count(minutes.start+10, minutes.end()-10),
		},
		"within hours": {
			measurement: testmeasurement,
			from:        hours.start + 10,
			until:       hours.end() - 10,
			count:       hours.count(hours.start+10, hours.end()-10),
		},
		"range goes beyond end of hours": {
			measurement: testmeasurement,
			from:        hours.start + 10,
			until:       hours.end() + 10,
			count:       hours.count(hours.start+10, hours.end()),
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			actual := c.measurement.Count(c.from, c.until)
			assert.EqualValues(t, int(c.count), int(actual))
			if c.count != actual {
				fmt.Println(c.measurement.series[0])
			}
		})
	}
}
