package validation

import "strings"

// digitsOnly returns only digit runes from s.
func digitsOnly(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// validCPF checks length (11 digits), rejects all-same-digit, and validates check digits.
func validCPF(cpf string) bool {
	digits := digitsOnly(cpf)
	if len(digits) != 11 {
		return false
	}
	// Reject all same digit (e.g. 111.111.111-11)
	first := digits[0]
	for i := 1; i < 11; i++ {
		if digits[i] != first {
			break
		}
		if i == 10 {
			return false
		}
	}
	// First check digit: sum(d[i]*(10-i)) for i=0..8, then (sum*10)%11
	var sum int
	for i := 0; i < 9; i++ {
		sum += int(digits[i]-'0') * (10 - i)
	}
	d10 := (sum * 10) % 11
	if d10 == 10 {
		d10 = 0
	}
	if int(digits[9]-'0') != d10 {
		return false
	}
	// Second check digit: sum(d[i]*(11-i)) for i=0..9, then (sum*10)%11
	sum = 0
	for i := 0; i < 10; i++ {
		sum += int(digits[i]-'0') * (11 - i)
	}
	d11 := (sum * 10) % 11
	if d11 == 10 { // when remainder is 10, check digit is 0 (no valid CPF has 11th digit 0 with sum2≡1 mod 11, but the rule is correct)
		d11 = 0
	}
	return int(digits[10]-'0') == d11
}
