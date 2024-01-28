package luna

import (
	"errors"
	"fmt"
	"strconv"
)

var (
	ErrWrongFormat = errors.New("wrong format")
)

// Check - Проверка корректности номера используя алгоритм Луна
func Check(num string) (isCorrect bool, err error) {
	if len(num) <= 1 {
		return false, ErrWrongFormat
	}

	controlDigit, err := strconv.Atoi(string(num[len(num)-1]))
	if err != nil {
		return false, fmt.Errorf("luna - Check - strconv.Atoi: %w", ErrWrongFormat)
	}
	sum, err := Sum(num[:len(num)-1])
	if err != nil {
		return false, fmt.Errorf("luna - Check - Sum: %w", err)
	}

	if (sum+controlDigit)%10 == 0 {
		return true, nil
	}
	return false, nil
}

func Sum(num string) (int, error) {
	if len(num) == 0 {
		return 0, ErrWrongFormat
	}

	sum := 0
	parity := len(num) % 2
	for i, b := range num {
		n, err := strconv.Atoi(string(b))
		if err != nil {
			return 0, ErrWrongFormat
		}

		if i%2 != parity {
			n = n * 2
		}
		if n > 9 {
			n -= 9
		}

		sum += n
	}
	return sum, nil
}
