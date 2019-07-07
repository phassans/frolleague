package phantom

import (
	"fmt"
	"strconv"
	"strings"
)

func parseDateRangeForSchool(dateRange string) (FromYear, ToYear, error) {
	parts := strings.Split(dateRange, " – ")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid dateRange: %s", dateRange)
	}

	from, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, err
	}

	to, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, err
	}

	return FromYear(from), ToYear(to), nil
}

func parseDateRangeForCompanies(dateRange string) (FromYear, ToYear, error) {
	parts := strings.Split(dateRange, " – ")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid dateRange: %s", dateRange)
	}

	fromParts := strings.Split(parts[0], " ")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid fromParts: %s", parts[0])
	}

	from, err := strconv.Atoi(fromParts[1])
	if err != nil {
		return 0, 0, err
	}

	var to int
	if parts[1] != "Present" {
		toParts := strings.Split(parts[1], " ")
		if len(parts) != 2 {
			return 0, 0, fmt.Errorf("invalid toParts: %s", parts[1])
		}

		to, err = strconv.Atoi(toParts[1])
		if err != nil {
			return 0, 0, err
		}
	} else {
		to = 2018
	}

	return FromYear(from), ToYear(to), nil
}
