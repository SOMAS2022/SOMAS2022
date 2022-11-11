package config

import (
	"log"
	"os"
	"strconv"
)

func EnvIsSet(key string) bool {
	_, err := strconv.ParseUint(os.Getenv(key), 10, 0)
	return err == nil
}

func EnvToUint(key string, def uint) uint {
	levels, err := strconv.ParseUint(os.Getenv(key), 10, 0)
	if err != nil {
		log.Printf("%s unset, defaulting to %d\n", key, def)
		return def
	}
	return uint(levels)
}

func EnvToFloat(key string, def float32) float32 {
	levels, err := strconv.ParseFloat(os.Getenv(key), 32)
	if err != nil {
		log.Printf("%s unset, defaulting to %f\n", key, def)
		return def
	}
	return float32(levels)
}
