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
	"path/filepath"
	"sync/atomic"
	"time"

	"github.com/LucasBastino/app-sindicato/internal/config"
	"github.com/LucasBastino/app-sindicato/internal/infra/logger"
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
	
	homeDir, err := os.UserHomeDir()
	if err!=nil{
		return false, err
	}
		
	// licensePath := s.cfg.Path
	// licensePath := os.Getenv("LICENSE_PATH")
	// if licensePath == "" {
	// 	licensePath = "Licencias/licencia.json" // fallback
	// }
	// todo: cambiar esto de filepath despues
	filePath := filepath.Join(homeDir, "Desktop", "Licencias", "license.json")
	fmt.Println(filePath)
	file, err := os.ReadFile(filePath)
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
	keyPath := filepath.Join(homeDir, "Desktop", "Licencias", "public.pem")
	fmt.Println(keyPath)
	// keyPath := s.cfg.PublicKeyPath
	publicKeyData, err := os.ReadFile(keyPath)
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
