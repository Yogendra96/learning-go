package otp

import (
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"

	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"fmt"
	"hash"
	"image"
	"net/url"
	"strings"
)

// Error when attempting to convert the secret from base32 to raw bytes.
var ErrValidateSecretInvalidBase32 = errors.New("Decoding of secret as base32 failed.")

// The user provided passcode length was not expected.
var ErrValidateInputInvalidLength = errors.New("Input length unexpected")

// When generating a Key, the Issuer must be set.
var ErrGenerateMissingIssuer = errors.New("Issuer must be set")

// When generating a Key, the Account Name must be set.
var ErrGenerateMissingAccountName = errors.New("AccountName must be set")

// Key represents an TOTP or HTOP key.
type Key struct {
	orig string
	url  *url.URL
}

func NewKeyFromURL(orig string) (*Key, error) {
	s := strings.TrimSpace(orig)

	u, err := url.Parse(s)

	if err != nil {
		return nil, err
	}

	return &Key{
		orig: s,
		url:  u,
	}, nil
}

func (k *Key) String() string {
	return k.orig
}

func (k *Key) Image(width int, height int) (image.Image, error) {
	b, err := qr.Encode(k.orig, qr.M, qr.Auto)
	if err != nil {
		return nil, err
	}

	b, err = barcode.Scale(b, width, height)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// Type returns "hotp" or "totp".
func (k *Key) Type() string {
	return k.url.Host
}

func (k *Key) Issuer() string {
	q := k.url.Query()

	issuer := q.Get("issuer")
	if issuer != "" {
		return issuer
	}

	p := strings.TrimPrefix(k.url.Path, "/")
	i := strings.Index(p, ":")

	if i == -1 {
		return ""
	}

	return p[:i]
}

// AccountName returns the name of the user's account.
func (k *Key) AccountName() string {
	p := strings.TrimPrefix(k.url.Path, "/")
	i := strings.Index(p, ":")

	if i == -1 {
		return p
	}

	return p[i+1:]
}

// Secret returns the opaque secret for this Key.
func (k *Key) Secret() string {
	q := k.url.Query()

	return q.Get("secret")
}

// URL returns the OTP URL as a string
func (k *Key) URL() string {
	return k.url.String()
}

// Algorithm represents the hashing function to use in the HMAC
// operation needed for OTPs.
type Algorithm int

const (
	AlgorithmSHA1 Algorithm = iota
	AlgorithmSHA256
	AlgorithmSHA512
	AlgorithmMD5
)

func (a Algorithm) String() string {
	switch a {
	case AlgorithmSHA1:
		return "SHA1"
	case AlgorithmSHA256:
		return "SHA256"
	case AlgorithmSHA512:
		return "SHA512"
	case AlgorithmMD5:
		return "MD5"
	}
	panic("unreached")
}

func (a Algorithm) Hash() hash.Hash {
	switch a {
	case AlgorithmSHA1:
		return sha1.New()
	case AlgorithmSHA256:
		return sha256.New()
	case AlgorithmSHA512:
		return sha512.New()
	case AlgorithmMD5:
		return md5.New()
	}
	panic("unreached")
}

type Digits int

const (
	DigitsSix   Digits = 6
	DigitsEight Digits = 8
)

// Format converts an integer into the zero-filled size for this Digits.
func (d Digits) Format(in int32) string {
	f := fmt.Sprintf("%%0%dd", d)
	return fmt.Sprintf(f, in)
}

// Length returns the number of characters for this Digits.
func (d Digits) Length() int {
	return int(d)
}

func (d Digits) String() string {
	return fmt.Sprintf("%d", d)
}
