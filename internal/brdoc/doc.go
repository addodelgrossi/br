package brdoc

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strings"
)

func OnlyDigits(input string) string {
	var b strings.Builder
	for _, r := range input {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func ValidateCPF(input string) bool {
	digits := OnlyDigits(input)
	if len(digits) != 11 || allSame(digits) {
		return false
	}

	nums := toInts(digits)
	first := cpfCheckDigit(nums[:9])
	second := cpfCheckDigit(append(nums[:9], first))

	return nums[9] == first && nums[10] == second
}

func FormatCPF(input string) (string, error) {
	digits := OnlyDigits(input)
	if len(digits) != 11 {
		return "", errors.New("CPF deve conter 11 digitos")
	}
	return fmt.Sprintf("%s.%s.%s-%s", digits[0:3], digits[3:6], digits[6:9], digits[9:11]), nil
}

func GenerateCPF() (string, error) {
	for {
		nums, err := randomDigits(9)
		if err != nil {
			return "", err
		}

		nums = append(nums, cpfCheckDigit(nums))
		nums = append(nums, cpfCheckDigit(nums))
		digits := intsToString(nums)
		if !allSame(digits) {
			return digits, nil
		}
	}
}

func ValidateCNPJ(input string) bool {
	digits := OnlyDigits(input)
	if len(digits) != 14 || allSame(digits) {
		return false
	}

	nums := toInts(digits)
	first := cnpjCheckDigit(nums[:12])
	second := cnpjCheckDigit(append(nums[:12], first))

	return nums[12] == first && nums[13] == second
}

func FormatCNPJ(input string) (string, error) {
	digits := OnlyDigits(input)
	if len(digits) != 14 {
		return "", errors.New("CNPJ deve conter 14 digitos")
	}
	return fmt.Sprintf("%s.%s.%s/%s-%s", digits[0:2], digits[2:5], digits[5:8], digits[8:12], digits[12:14]), nil
}

func GenerateCNPJ() (string, error) {
	for {
		nums, err := randomDigits(12)
		if err != nil {
			return "", err
		}

		nums = append(nums, cnpjCheckDigit(nums))
		nums = append(nums, cnpjCheckDigit(nums))
		digits := intsToString(nums)
		if !allSame(digits) {
			return digits, nil
		}
	}
}

func cpfCheckDigit(nums []int) int {
	sum := 0
	for i, n := range nums {
		sum += n * (len(nums) + 1 - i)
	}

	remainder := sum % 11
	if remainder < 2 {
		return 0
	}
	return 11 - remainder
}

func cnpjCheckDigit(nums []int) int {
	weights := []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	offset := len(weights) - len(nums)

	sum := 0
	for i, n := range nums {
		sum += n * weights[offset+i]
	}

	remainder := sum % 11
	if remainder < 2 {
		return 0
	}
	return 11 - remainder
}

func randomDigits(length int) ([]int, error) {
	nums := make([]int, length)
	for i := range nums {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return nil, err
		}
		nums[i] = int(n.Int64())
	}
	return nums, nil
}

func allSame(input string) bool {
	for i := 1; i < len(input); i++ {
		if input[i] != input[0] {
			return false
		}
	}
	return true
}

func toInts(input string) []int {
	nums := make([]int, len(input))
	for i := range input {
		nums[i] = int(input[i] - '0')
	}
	return nums
}

func intsToString(nums []int) string {
	var b strings.Builder
	for _, n := range nums {
		b.WriteByte(byte('0' + n))
	}
	return b.String()
}
