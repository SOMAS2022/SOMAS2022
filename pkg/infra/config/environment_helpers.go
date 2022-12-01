package config

import (
	"fmt"
	"os"
	"strconv"

	"infra/logging"
)

func EnvToUint(key string, def uint) uint {
	levels, err := strconv.ParseUint(os.Getenv(key), 10, 0)
	if err != nil {
		logging.Log(logging.Warn, nil, fmt.Sprintf("%s unset, defaulting to %d\n", key, def))

		return def
	}

	return uint(levels)
}

func EnvToFloat(key string, def float32) float32 {
	levels, err := strconv.ParseFloat(os.Getenv(key), 32)
	if err != nil {
		logging.Log(logging.Warn, nil, fmt.Sprintf("%s unset, defaulting to %f\n", key, def))

		return def
	}

	return float32(levels)
}

func EnvToString(key string, def string) string {
	s := os.Getenv(key)
	if s == "" {
		logging.Log(logging.Warn, nil, fmt.Sprintf("%s unset, defaulting to %s\n", key, def))

		return def
	}

	return s
}
