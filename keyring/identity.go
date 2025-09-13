package keyring

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"path/filepath"
)

// Identity holds an ed25519 keypair used for GRISHINIUM peer identity.
type Identity struct {
	Private ed25519.PrivateKey
	Public  ed25519.PublicKey
}

// Fingerprint returns a short hex fingerprint of the public key.
func (id Identity) Fingerprint() string {
	if len(id.Public) == 0 {
		return ""
	}
	sum := id.Public
	if len(sum) > 8 {
		sum = sum[:8]
	}
	return hex.EncodeToString(sum)
}

// LoadIdentity loads an ed25519 private key from the path, or generates and saves a new one when the path is empty or file not found.
// The key is stored in raw ed25519 format.
func LoadIdentity(path string) (Identity, error) {
	if path == "" {
		// generate ephemeral
		pub, priv, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return Identity{}, err
		}
		return Identity{Private: priv, Public: pub}, nil
	}
	clean := filepath.Clean(path)
	b, err := os.ReadFile(clean)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// generate and save
			pub, priv, gerr := ed25519.GenerateKey(rand.Reader)
			if gerr != nil {
				return Identity{}, gerr
			}
			if werr := writeFileAtomic(clean, priv); werr != nil {
				return Identity{}, werr
			}
			return Identity{Private: priv, Public: pub}, nil
		}
		return Identity{}, err
	}
	// load
	priv := ed25519.PrivateKey(b)
	if len(priv) != ed25519.PrivateKeySize {
		return Identity{}, errors.New("invalid ed25519 private key size")
	}
	pub := priv.Public().(ed25519.PublicKey)
	return Identity{Private: priv, Public: pub}, nil
}

func writeFileAtomic(path string, data []byte) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(dir, ".keytmp-*")
	if err != nil {
		return err
	}
	defer func() {
		_ = os.Remove(tmp.Name())
	}()
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Chmod(0o600); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Sync(); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), path)
}
