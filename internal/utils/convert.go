package utils

import "strconv"

func ParseUintPtr(value string) *uint64 {
	if value == "" {
		return nil
	}

	id, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return nil
	}

	return &id
}

func Uint64OrZero(value string) uint64 {
	if value == "" {
		return 0
	}
	id, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0
	}
	return id
}

func parseUint(value string) uint64 {
	if value == "" {
		return 0
	}
	id, _ := strconv.ParseUint(value, 10, 64)
	return id
}
