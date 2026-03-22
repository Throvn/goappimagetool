package appimagetoolgo

import (
	"encoding/hex"

	"github.com/ProtonMail/gopenpgp/v2/helper"
)

func GeneratePGPPrivateKey(passphrase string) (string, error) {
	return helper.GenerateKey("My Key", "random@mail.com", []byte(passphrase), "rsa", 4096)
}

func SignSha256(hash []byte, privateKey string, passphrase string) (string, error) {
	hexlifiedHash := hex.EncodeToString(hash)
	// Keys initialization as before (omitted)

	return helper.SignCleartextMessageArmored(privateKey, []byte(passphrase), hexlifiedHash)
}
