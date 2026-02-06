package graph

import "fmt"

const (
	DefaultLimit int = 5
)

func CheckLimit[T any](arr []T, limit *int, offset *int) error {
	if len(arr) == 0 {
		*limit = 0
		*offset = 0
	}

	if *limit > DefaultLimit {
		*limit = DefaultLimit
	}
	if *limit < 1 {
		*limit = 1
	}
	if *offset < 0 {
		*offset = 0
	}

	if (*offset + *limit) > len(arr) {
		*limit = len(arr) - *offset
	}

	if *offset >= len(arr) {
		return fmt.Errorf("Offser must be less than len of the arr")
	}
	return nil
}
