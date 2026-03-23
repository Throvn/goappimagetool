package appimagetoolgo

import (
	"encoding/hex"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/ProtonMail/gopenpgp/v2/helper"
)

func GenerateSigningKey(name string, email string, passphrase string) (secretKey string, publicKey string, err error) {
	privateKey, err := helper.GenerateKey(name, email, []byte(passphrase), "rsa", 4096)
	if err != nil {
		return "", "", err
	}

	privateKeyObj, err := crypto.NewKeyFromArmored(privateKey)
	if err != nil {
		return privateKey, "", err
	}

	publicKey, err = privateKeyObj.GetArmoredPublicKey()
	return privateKey, publicKey, err
}

func SignSha256(hash []byte, privateKey string, passphrase string) (string, error) {
	hexlifiedHash := hex.EncodeToString(hash)
	return helper.SignCleartextMessageArmored(privateKey, []byte(passphrase), hexlifiedHash)
}
