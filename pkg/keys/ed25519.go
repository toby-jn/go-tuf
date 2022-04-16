package keys

import (
	"crypto"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/json"
	"errors"

	"github.com/theupdateframework/go-tuf/data"
)

func init() {
	SignerMap.Store(data.KeyTypeEd25519, newEd25519Signer)
	VerifierMap.Store(data.KeyTypeEd25519, newEd25519Verifier)
}

func newEd25519Signer() Signer {
	return &ed25519Signer{}
}

func newEd25519Verifier() Verifier {
	return &ed25519Verifier{}
}

type ed25519Verifier struct {
	PublicKey data.HexBytes `json:"public"`
	key       *data.PublicKey
}

func (e *ed25519Verifier) Public() string {
	return string(e.PublicKey)
}

func (e *ed25519Verifier) Verify(msg, sig []byte) error {
	if !ed25519.Verify([]byte(e.PublicKey), msg, sig) {
		return errors.New("tuf: ed25519 signature verification failed")
	}
	return nil
}

func (e *ed25519Verifier) MarshalPublicKey() *data.PublicKey {
	return e.key
}

func (e *ed25519Verifier) UnmarshalPublicKey(key *data.PublicKey) error {
	e.key = key
	if err := json.Unmarshal(key.Value, e); err != nil {
		return err
	}
	if len(e.PublicKey) != ed25519.PublicKeySize {
		return errors.New("tuf: unexpected public key length for ed25519 key")
	}
	return nil
}

type Ed25519PrivateKeyValue struct {
	Public  data.HexBytes `json:"public"`
	Private data.HexBytes `json:"private"`
}

type ed25519Signer struct {
	ed25519.PrivateKey
}

func GenerateEd25519Key() (*ed25519Signer, error) {
	_, private, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return &ed25519Signer{
		PrivateKey: ed25519.PrivateKey(data.HexBytes(private)),
	}, nil
}

func NewEd25519Signer(keyValue Ed25519PrivateKeyValue) *ed25519Signer {
	return &ed25519Signer{
		PrivateKey: ed25519.PrivateKey(data.HexBytes(keyValue.Private)),
	}
}

func (e *ed25519Signer) SignMessage(message []byte) ([]byte, error) {
	return e.Sign(rand.Reader, message, crypto.Hash(0))
}

func (e *ed25519Signer) MarshalPrivateKey() (*data.PrivateKey, error) {
	valueBytes, err := json.Marshal(Ed25519PrivateKeyValue{
		Public:  data.HexBytes([]byte(e.PrivateKey.Public().(ed25519.PublicKey))),
		Private: data.HexBytes(e.PrivateKey),
	})
	if err != nil {
		return nil, err
	}
	return &data.PrivateKey{
		Type:       data.KeyTypeEd25519,
		Scheme:     data.KeySchemeEd25519,
		Algorithms: data.HashAlgorithms,
		Value:      valueBytes,
	}, nil
}

func (e *ed25519Signer) UnmarshalPrivateKey(key *data.PrivateKey) error {
	keyValue := &Ed25519PrivateKeyValue{}
	if err := json.Unmarshal(key.Value, keyValue); err != nil {
		return err
	}
	*e = ed25519Signer{
		PrivateKey: ed25519.PrivateKey(data.HexBytes(keyValue.Private)),
	}
	return nil
}

func (e *ed25519Signer) PublicData() *data.PublicKey {
	keyValBytes, _ := json.Marshal(ed25519Verifier{PublicKey: []byte(e.PrivateKey.Public().(ed25519.PublicKey))})
	return &data.PublicKey{
		Type:       data.KeyTypeEd25519,
		Scheme:     data.KeySchemeEd25519,
		Algorithms: data.HashAlgorithms,
		Value:      keyValBytes,
	}
}
