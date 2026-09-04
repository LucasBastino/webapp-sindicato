package license

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"
	"sync/atomic"
	"time"

	"github.com/LucasBastino/webapp-sindicato/internal/config"
	"github.com/LucasBastino/webapp-sindicato/internal/infra/logger"
)

type LicenseService struct {
    logger logger.Logger
	cfg config.LicenseConfig
}

func NewLicenseService(logger logger.Logger, cfg config.LicenseConfig) *LicenseService {
    return &LicenseService{
		logger: logger,
		cfg: cfg,
	}
}

var licenseValid atomic.Bool


func (s *LicenseService) StartChecker() {
	valid, err := s.Check()
	if err != nil {
		s.logger.Error("failed to check license", "err", err)
	}
	s.SetAtomic(valid)

	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		for range ticker.C {
			valid,  err := s.Check()
			if err != nil {
				s.logger.Error("failed to check license", "err", err)
			}
			s.SetAtomic(valid)
		}
	}()
}

func (s *LicenseService) IsValid() bool {
    return licenseValid.Load()
}


func (s *LicenseService) Check() (bool, error) {
	licensePath := s.cfg.Path
	if licensePath == "" {
		licensePath = "./license.json"
	}

	publicKeyPath := s.cfg.PublicKeyPath
	if publicKeyPath == "" {
		publicKeyPath = "./config/license/public.pem"
	}

	file, err := os.ReadFile(licensePath)
	if err != nil {
		return false, err
	}
	var licenseFile LicenseFile
	err = json.Unmarshal(file, &licenseFile)
	if err != nil {
		return false, err
	}

	payloadBytes, err := json.Marshal(licenseFile.Payload)
	if err != nil {
		return false, err
	}

	hash := sha256.Sum256(payloadBytes)

	signatureBytes, err := base64.StdEncoding.DecodeString(licenseFile.Signature)
	if err != nil {
		return false, err
	}
	publicKeyData, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return false, err
	}

	block, _ := pem.Decode(publicKeyData)
	if block == nil{
		return false, fmt.Errorf("invalid public key PEM")
	}

	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return false, fmt.Errorf("failed to parse public key: %w", err)
	}

	publicKey := publicKeyInterface.(*rsa.PublicKey)

	
	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash[:], signatureBytes)
	if err != nil {
		s.logger.Error("license rejected", "reason", "invalid signature", "err", err)
		return false, nil
	}

	if licenseFile.Payload.Company != "Sindicato" {
		s.logger.Error("license rejected", "reason", "invalid company")
		return false, nil
	}

	if time.Now().After(licenseFile.Payload.Exp){
		s.logger.Error("license rejected", "reason", "expired license")
		return false, nil
	}

	return true, nil
}


func (s *LicenseService) SetAtomic(valid bool) {
	licenseValid.Store(valid)
}
