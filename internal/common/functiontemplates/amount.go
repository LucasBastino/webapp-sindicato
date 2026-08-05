package functiontemplates

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// FormatAmountAR formats a monetary amount for Argentine display:
// thousands with ".", decimals with ",", and no trailing ",00".
// Examples: 9 -> "9"; 9.5 -> "9,5"; 9.05 -> "9,05"; 1500.5 -> "1.500,5"
func FormatAmountAR(amount any) string {
	return formatAmount(amount, true)
}

// FormatAmountInput formats an amount for editable form fields:
// no thousand separators, decimals with ",", and no trailing ",00".
// Examples: 80593 -> "80593"; 80593.5 -> "80593,5"; 80593.05 -> "80593,05"
func FormatAmountInput(amount any) string {
	return formatAmount(amount, false)
}

func formatAmount(amount any, withThousands bool) string {
	f, ok := toFloat64(amount)
	if !ok {
		return ""
	}

	negative := f < 0
	if negative {
		f = -f
	}

	cents := int64(math.Round(f * 100))
	whole := cents / 100
	frac := cents % 100

	intPart := strconv.FormatInt(whole, 10)
	if withThousands {
		intPart = formatThousandsAR(whole)
	}

	var out string
	switch {
	case frac == 0:
		out = intPart
	case frac%10 == 0:
		out = fmt.Sprintf("%s,%d", intPart, frac/10)
	default:
		out = fmt.Sprintf("%s,%02d", intPart, frac)
	}

	if negative {
		return "-" + out
	}
	return out
}

func toFloat64(amount any) (float64, bool) {
	switch v := amount.(type) {
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	case int32:
		return float64(v), true
	case string:
		normalized := strings.ReplaceAll(strings.TrimSpace(v), ",", ".")
		if normalized == "" {
			return 0, false
		}
		f, err := strconv.ParseFloat(normalized, 64)
		if err != nil {
			return 0, false
		}
		return f, true
	default:
		return 0, false
	}
}

func formatThousandsAR(n int64) string {
	s := strconv.FormatInt(n, 10)
	if len(s) <= 3 {
		return s
	}

	var b strings.Builder
	rem := len(s) % 3
	if rem > 0 {
		b.WriteString(s[:rem])
		if len(s) > rem {
			b.WriteByte('.')
		}
	}
	for i := rem; i < len(s); i += 3 {
		b.WriteString(s[i : i+3])
		if i+3 < len(s) {
			b.WriteByte('.')
		}
	}
	return b.String()
}
