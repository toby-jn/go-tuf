package keys

import (
	"github.com/theupdateframework/go-tuf/data"
	. "gopkg.in/check.v1"
)

type RsaSuite struct{}

var _ = Suite(&RsaSuite{})

func (RsaSuite) TestSignVerify(c *C) {
	signer, err := GenerateRsaKey()
	c.Assert(err, IsNil)
	msg := []byte("foo")
	sig, err := signer.SignMessage(msg)
	c.Assert(err, IsNil)
	publicData := signer.PublicData()
	pubKey, err := GetVerifier(publicData)
	c.Assert(err, IsNil)
	c.Assert(pubKey.Verify(msg, sig), IsNil)
}

func (RsaSuite) TestMarshalUnmarshal(c *C) {
	signer, err := GenerateRsaKey()
	c.Assert(err, IsNil)
	publicData := signer.PublicData()
	pubKey, err := GetVerifier(publicData)
	c.Assert(err, IsNil)
	c.Assert(pubKey.MarshalPublicKey(), DeepEquals, publicData)
}

func (RsaSuite) TestMarshalUnmarshalPrivateKey(c *C) {
	signer, err := GenerateRsaKey()
	c.Assert(err, IsNil)
	privateData, err := signer.MarshalPrivateKey()
	c.Assert(err, IsNil)
	c.Assert(privateData.Type, Equals, data.KeyTypeRSASSA_PSS_SHA256)
	c.Assert(privateData.Scheme, Equals, data.KeySchemeRSASSA_PSS_SHA256)
	c.Assert(privateData.Algorithms, DeepEquals, data.HashAlgorithms)
	s, err := GetSigner(privateData)
	c.Assert(err, IsNil)
	c.Assert(s, DeepEquals, signer)
}
