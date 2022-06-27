package metastore

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"time"
)

// ErrCannotStoreRSAPrivateKey indicate key is generated but cannot successfully store and fetch from storage.
var ErrCannotStoreRSAPrivateKey = errors.New("cannot store generated RSA private key into storage")

func packPKCS1PrivateKey(priKey *rsa.PrivateKey) (keyText string) {
	buf := x509.MarshalPKCS1PrivateKey(priKey)
	keyText = base64.StdEncoding.EncodeToString(buf)
	return
}

func unpackPKCS1PrivateKey(keyText string) (priKey *rsa.PrivateKey, err error) {
	var buf []byte
	if buf, err = base64.StdEncoding.DecodeString(keyText); nil != err {
		return
	}
	priKey, err = x509.ParsePKCS1PrivateKey(buf)
	return
}

// FetchRSAPrivateKey read RSA private key from storage.
//
// A new private key will be generate if existed key expires.
//
// Caller can keep resulted `modifyAt` time-stamp and feed it as `currentModifyAt`
// on next invoke. Key parsing will be skip and resulted `ok` will be false
// if the modification time-stamp (ie: `modifyAt`) is not changed.
func (m *MetaStore) FetchRSAPrivateKey(metaKey string, keyBits int, maxAcceptableAge time.Duration, currentModifyAt int64) (ok bool, priKey *rsa.PrivateKey, modifyAt int64, err error) {
	ok, keyText, modifyAt, err := m.fetch(metaKey)
	if nil != err {
		return
	}
	modifyBoundAt := time.Now().Unix() - int64(maxAcceptableAge/time.Second)
	if ok && (modifyBoundAt < modifyAt) {
		if currentModifyAt == modifyAt {
			ok = false
			return
		}
		if priKey, err = unpackPKCS1PrivateKey(keyText); nil == err {
			return
		}
	}
	priKey, err = rsa.GenerateKey(rand.Reader, keyBits)
	if nil != err {
		return
	}
	keyText = packPKCS1PrivateKey(priKey)
	if err = m.store(metaKey, keyText); nil != err {
		priKey = nil
		return
	}
	ok, updatedKeyText, modifyAt, err := m.fetch(metaKey)
	if nil != err {
		priKey = nil
		return
	}
	if !ok {
		priKey = nil
		err = ErrCannotStoreRSAPrivateKey
		return
	}
	if updatedKeyText != keyText {
		if priKey, err = unpackPKCS1PrivateKey(keyText); nil != err {
			ok = false
		}
	}
	return
}
