package license

import "time"

type LicensePayload struct {
	Company string    `json:"company"`
	Exp        time.Time `json:"exp"`
}

type LicenseFile struct {
	Payload   LicensePayload `json:"payload"`
	Signature string         `json:"signature"`
}