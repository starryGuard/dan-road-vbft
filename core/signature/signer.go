package signature

import (
	"dan-road-vbft/crypto/keypair"
	"dan-road-vbft/crypto/signature"
)

// Signer is the abstract interface of user's information(Keys) for signing data.
type Signer interface {
	//get signer's private key
	PrivKey() keypair.PrivateKey

	//get signer's public key
	PubKey() keypair.PublicKey

	Scheme() signature.SignatureScheme
}
