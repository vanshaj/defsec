package functions

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

var pattern = regexp.MustCompile(`^P((?P<year>\d+)Y)?((?P<month>\d+)M)?((?P<week>\d+)W)?((?P<day>\d+)D)?(T((?P<hour>\d+)H)?((?P<minute>\d+)M)?((?P<second>\d+)S)?)?$`)

func DateTimeAdd(args ...interface{}) interface{} {
	if len(args) < 2 {
		return nil
	}

	base, ok := args[0].(string)
	if !ok {
		return nil
	}

	baseTime, err := time.Parse(time.RFC3339, base)
	if err != nil {
		return nil
	}

	duration, err := parseISO8601(args[1].(string))
	if err != nil {
		return nil
	}

	timeDuration := duration.timeDuration()
	baseTime = baseTime.Add(timeDuration)

	if len(args) == 3 {
		format, ok := args[2].(string)
		if ok {
			return baseTime.Format(format)
		}
	}

	return baseTime.Format(time.RFC3339)
}

type ISO8601Duration struct {
	Y int
	M int
	W int
	D int
	// Time Component
	TH int
	TM int
	TS int
}

func parseISO8601(from string) (ISO8601Duration, error) {
	var match []string
	var d ISO8601Duration

	if pattern.MatchString(from) {
		match = pattern.FindStringSubmatch(from)
	} else {
		return d, fmt.Errorf("could not parse duration string")
	}

	for i, name := range pattern.SubexpNames() {
		part := match[i]
		if i == 0 || name == "" || part == "" {
			continue
		}

		val, err := strconv.Atoi(part)
		if err != nil {
			return d, err
		}
		switch name {
		case "year":
			d.Y = val
		case "month":
			d.M = val
		case "week":
			d.W = val
		case "day":
			d.D = val
		case "hour":
			d.TH = val
		case "minute":
			d.TM = val
		case "second":
			d.TS = val
		default:
			return d, fmt.Errorf("unknown field %s", name)
		}
	}

	return d, nil
}

func (d ISO8601Duration) timeDuration() time.Duration {
	var dur time.Duration
	dur += time.Duration(d.TH) * time.Hour
	dur += time.Duration(d.TM) * time.Minute
	dur += time.Duration(d.TS) * time.Second
	dur += time.Duration(d.D) * 24 * time.Hour
	dur += time.Duration(d.W) * 7 * 24 * time.Hour
	dur += time.Duration(d.M) * 30 * 24 * time.Hour
	dur += time.Duration(d.Y) * 365 * 24 * time.Hour

	return dur
}
