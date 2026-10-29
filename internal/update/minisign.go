package update

import (
	"bytes"
	"crypto/ed25519"
	"encoding/base64"
	"errors"
	"strings"
)

// SigningPublicKey is the minisign public key release SHA256SUMS files are
// signed with: the shared release key DBC46B23C3AF0EF3, the same one zwriter
// and zedit use. Not secret. The matching secret key lives only on the
// maintainer's machine and in this repository's MINISIGN_SECRET_KEY /
// MINISIGN_PASSWORD Actions secrets.
const SigningPublicKey = "RWTzDq/DI2vE2xMeF01QWTxP2lwwlxGL9M1UxQE4EFV5kVmTF84wwyi2"

// VerifyMinisign checks a minisign signature file (the contents of
// SHA256SUMS.minisig) over msg against a base64 public key line. Both the
// prehashed ("ED", minisign's default) and legacy ("Ed") forms are accepted,
// and the trusted comment's global signature must verify too.
func VerifyMinisign(publicKeyB64 string, msg, sigFile []byte) error {
	pk, err := base64.StdEncoding.DecodeString(strings.TrimSpace(publicKeyB64))
	if err != nil || len(pk) != 42 || pk[0] != 'E' || pk[1] != 'd' {
		return errors.New("invalid minisign public key")
	}
	lines := strings.Split(strings.ReplaceAll(string(sigFile), "\r\n", "\n"), "\n")
	if len(lines) < 4 || !strings.HasPrefix(lines[0], "untrusted comment:") {
		return errors.New("malformed signature file")
	}
	sig, err := base64.StdEncoding.DecodeString(strings.TrimSpace(lines[1]))
	if err != nil || len(sig) != 74 {
		return errors.New("malformed signature")
	}
	const tcPrefix = "trusted comment: "
	if !strings.HasPrefix(lines[2], tcPrefix) {
		return errors.New("signature has no trusted comment")
	}
	trusted := []byte(strings.TrimPrefix(lines[2], tcPrefix))
	global, err := base64.StdEncoding.DecodeString(strings.TrimSpace(lines[3]))
	if err != nil || len(global) != 64 {
		return errors.New("malformed global signature")
	}
	if !bytes.Equal(sig[2:10], pk[2:10]) {
		return errors.New("signed by a different key")
	}
	pub := ed25519.PublicKey(pk[10:])
	var signed []byte
	switch string(sig[:2]) {
	case "ED":
		d := blake2b512(msg)
		signed = d[:]
	case "Ed":
		signed = msg
	default:
		return errors.New("unsupported signature algorithm")
	}
	if !ed25519.Verify(pub, signed, sig[10:]) {
		return errors.New("signature verification failed")
	}
	if !ed25519.Verify(pub, append(append([]byte{}, sig[10:]...), trusted...), global) {
		return errors.New("trusted comment signature verification failed")
	}
	return nil
}
