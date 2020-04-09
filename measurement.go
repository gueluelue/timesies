package timesies

// Measurement is represents a set of series that share the same key and labels.
type Measurement struct {
	Key    string
	Labels []string

	series []*series // ordered list of series by ascending granularity.
}

// Count returns the number of points within the defined timerange.
func (m *Measurement) Count(from, until uint32) uint32 {
	if m == nil {
		return 0
	}

	var s *series
	for _, s = range m.series {
		// because the series are sorted, we reach this point only if we're looking
		// at the future relative to the series timerange.
		if until < s.start {
			return 0
		}

		// we haven't found the most accurate timeseries if the end is outside of
		// time series' timerange.
		if until > s.end() {
			continue
		}

		// We're in the right timeseries
		return s.count(from, until)
	}

	// if no other time series matches, use the largest one
	return s.count(from, until)
}
