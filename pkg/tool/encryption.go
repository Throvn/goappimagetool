package goappimagetool

import (
	"encoding/hex"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/ProtonMail/gopenpgp/v2/helper"
)

// Creates a pgp key pair using RSA 4096 bit encryption.
// The armored public key and private key are written to two separate files called public.asc and private.asc.
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

// Signs the byte array using the supplied private key.
func SignSha256(hash []byte, pgp PGPMaterial) ([]byte, error) {
	hexlifiedHash := hex.EncodeToString(hash)
	signedMsg, err := helper.SignCleartextMessageArmored(pgp.PrivateKeyArmored, []byte(pgp.Passphrase), hexlifiedHash)
	if err != nil {
		return nil, err
	}

	rawSignedMsg, err := crypto.NewClearTextMessageFromArmored(signedMsg)
	return rawSignedMsg.GetBinarySignature(), err
}

// Updates the .AppImage elf binary by writing the public key into the .sig_key field.
func UpdateSigKey(path string, pgp PGPMaterial) error {
	privateKeyObj, err := crypto.NewKeyFromArmored(pgp.PrivateKeyArmored)
	if err != nil {
		return err
	}

	publicKey, err := privateKeyObj.GetPublicKey()
	return OverwriteSection(path, ".sig_key", publicKey)
}
