package functiontemplates

import "testing"

func TestFormatAmountAR(t *testing.T) {
	cases := map[float64]string{
		9:      "9",
		9.5:    "9,5",
		9.05:   "9,05",
		1500.5: "1.500,5",
		1500:   "1.500",
		0:      "0",
		0.01:   "0,01",
		-9.5:   "-9,5",
	}
	for in, want := range cases {
		got := FormatAmountAR(in)
		if got != want {
			t.Fatalf("FormatAmountAR(%v) = %q, want %q", in, got, want)
		}
	}
	if got := FormatAmountAR(float32(9.05)); got != "9,05" {
		t.Fatalf("float32: got %q", got)
	}
}

func TestFormatAmountInput(t *testing.T) {
	cases := map[float64]string{
		9:         "9",
		9.5:       "9,5",
		9.05:      "9,05",
		80593:     "80593",
		80593.5:   "80593,5",
		80593.05:  "80593,05",
		1500:      "1500",
		1500.5:    "1500,5",
		0:         "0",
		-80593.05: "-80593,05",
	}
	for in, want := range cases {
		got := FormatAmountInput(in)
		if got != want {
			t.Fatalf("FormatAmountInput(%v) = %q, want %q", in, got, want)
		}
	}
}
