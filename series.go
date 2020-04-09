package timesies

import "fmt"

// A series is a collection of points.
type series struct {
	start  uint32 // unix timestamp
	length uint32 // theoretical length of the series in points
	unit   uint32 // number of seconds per point in the series.
	points []*point
}

// count returns the sum of all the points' count values with value between from
// (included) and until (excluded).
func (s *series) count(from, until uint32) uint32 {
	if s == nil || len(s.points) == 0 {
		return 0
	}

	start, end := s.index(from), s.index(until)

	var cnt uint32
	for i := start; i < end; i++ {
		cnt += s.points[i].getCount()
	}

	return cnt
}

// end returns the first second after the series timerange as a unix timestamp.
func (s *series) end() uint32 {
	return s.start + uint32(s.length*s.unit)
}

// index return the index in the series that holds this timestamp. If the
// timestamp is before the series, the returned value is zero. If it is after,
// then the length of the series.
func (s *series) index(ts uint32) int {
	if ts >= s.end() {
		return int(s.length)
	}

	if ts < s.start {
		return 0
	}

	return int((ts - s.start) / s.unit)
}

func (s *series) String() string {
	return fmt.Sprintf("%v", s.points)
}
