package config

import (
	"strings"
	"testing"
)

func TestValidate_RejectsEmptyJWTSecret(t *testing.T) {
	cfg := validConfig()
	cfg.Auth.JWTSecret = ""
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for empty JWT_SECRET")
	}
}

func TestValidate_RejectsShortJWTSecret(t *testing.T) {
	cfg := validConfig()
	cfg.Auth.JWTSecret = "too-short"
	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected error for short JWT_SECRET")
	}
	if !strings.Contains(err.Error(), "at least") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidate_OK(t *testing.T) {
	cfg := validConfig()
	if err := cfg.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseCookieSecure(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"", true},
		{"true", true},
		{"1", true},
		{"yes", true},
		{"false", false},
		{"0", false},
		{"no", false},
		{"FALSE", false},
	}
	for _, tc := range cases {
		if got := parseCookieSecure(tc.in); got != tc.want {
			t.Fatalf("parseCookieSecure(%q)=%v want %v", tc.in, got, tc.want)
		}
	}
}

func TestParseDBTLS(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"", ""},
		{"false", ""},
		{"0", ""},
		{"no", ""},
		{"FALSE", ""},
		{"true", "true"},
		{"1", "true"},
		{"yes", "true"},
		{"TRUE", "true"},
		{"skip-verify", "skip-verify"},
		{"SKIP-VERIFY", "skip-verify"},
		{"tru", "tru"},
	}
	for _, tc := range cases {
		if got := parseDBTLS(tc.in); got != tc.want {
			t.Fatalf("parseDBTLS(%q)=%q want %q", tc.in, got, tc.want)
		}
	}
}

func TestValidate_RejectsInvalidDBTLS(t *testing.T) {
	cfg := validConfig()
	cfg.Database.TLS = "tru"
	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected error for invalid DB_TLS")
	}
	if !strings.Contains(err.Error(), "DB_TLS") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidate_AcceptsDBTLSModes(t *testing.T) {
	for _, mode := range []string{"", "true", "skip-verify"} {
		cfg := validConfig()
		cfg.Database.TLS = mode
		if err := cfg.Validate(); err != nil {
			t.Fatalf("mode %q: unexpected error: %v", mode, err)
		}
	}
}

func validConfig() Config {
	return Config{
		Database: DatabaseConfig{
			User: "u",
			Host: "localhost",
			Port: "3306",
			Name: "db",
		},
		Auth: AuthConfig{
			JWTSecret: strings.Repeat("x", minJWTSecretLen),
		},
		License: LicenseConfig{
			Path:          "./license.json",
			PublicKeyPath: "./public.pem",
		},
		Server: ServerConfig{
			Port: "8080",
			Host: "0.0.0.0",
		},
	}
}
