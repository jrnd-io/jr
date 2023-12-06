package emitter

import (
	"fmt"
	"regexp"
	"strconv"
)

type Throughput float64

const (
	bitMultiplier = 8
)

func parseThroughput(input string) (Throughput, error) {

	if input == "" {
		return -1, nil
	}
	re := regexp.MustCompile(`^((?:0|[1-9]\d*)(?:\.?\d*))([KkMmGgTt][Bb])/([smhd])$`)
	match := re.FindStringSubmatch(input)

	if len(match) != 4 {
		return 0, fmt.Errorf("invalid input format: %s", input)
	}

	valueStr := match[1]
	unitStr := match[2]
	timeStr := match[3]
	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse numeric value: %w", err)
	}

	switch timeStr {
	case "s":
		// nothing to do
	case "m":
		value = value / 60
	case "h":
		value = value / 3600
	case "d":
		value = value / 86400
	default:
		return 0, fmt.Errorf("unsupported time unit: %s", unitStr)
	}

	switch unitStr {
	case "b":
		return Throughput(value * bitMultiplier), nil
	case "B":
		return Throughput(value), nil
	case "kb", "Kb":
		return Throughput(value * 1024 * bitMultiplier), nil
	case "mb", "Mb":
		return Throughput(value * 1024 * 1024 * bitMultiplier), nil
	case "gb", "Gb":
		return Throughput(value * 1024 * 1024 * 1024 * bitMultiplier), nil
	case "tb", "Tb":
		return Throughput(value * 1024 * 1024 * 1024 * 1024 * bitMultiplier), nil
	case "kB", "KB":
		return Throughput(value * 1024), nil
	case "mB", "MB":
		return Throughput(value * 1024 * 1024), nil
	case "gB", "GB":
		return Throughput(value * 1024 * 1024 * 1024), nil
	case "tB", "TB":
		return Throughput(value * 1024 * 1024 * 1024 * 1024), nil
	default:
		return 0, fmt.Errorf("unsupported unit: %s", unitStr)
	}
}
