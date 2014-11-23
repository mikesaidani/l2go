package crypt

import (
  "errors"
	"github.com/frostwind/l2go/loginserver/crypt/blowfish"
)

func Checksum(raw []byte) bool {
	var chksum int = 0
	count := len(raw) - 8
	i := 0

	for i = 0; i < count; i += 4 {
		var ecx int = int(raw[i])
		ecx |= int(raw[i+1]) << 8
		ecx |= int(raw[i+2]) << 0x10
		ecx |= int(raw[i+3]) << 0x18
		chksum ^= ecx
	}

	var ecx int = int(raw[i])
	ecx |= int(raw[i+1]) << 8
	ecx |= int(raw[i+2]) << 0x10
	ecx |= int(raw[i+3]) << 0x18

	raw[i] = byte(chksum)
	raw[i+1] = byte(chksum >> 0x08)
	raw[i+2] = byte(chksum >> 0x10)
	raw[i+3] = byte(chksum >> 0x18)

	return ecx == chksum
}

func BlowfishDecrypt(encrypted, key []byte) ([]byte, error) {
	cipher, err := blowfish.NewCipher(key)

	if err != nil {
		return nil, errors.New("Couldn't initialize the blowfish cipher")
	}

	// Check if the encrypted data is a multiple of our block size
	if len(encrypted)%8 != 0 {
		return nil, errors.New("The encrypted data is not a multiple of the block size")
	}

	count := len(encrypted) / 8

	decrypted := make([]byte, len(encrypted))

	for i := 0; i < count; i++ {
		cipher.Decrypt(decrypted[i*8:], encrypted[i*8:])
	}

	return decrypted, nil
}

func BlowfishEncrypt(decrypted, key []byte) ([]byte, error) {
	cipher, err := blowfish.NewCipher(key)

	if err != nil {
		return nil, errors.New("Couldn't initialize the blowfish cipher")
	}

	// Check if the decrypted data is a multiple of our block size
	if len(decrypted)%8 != 0 {
		return nil, errors.New("The decrypted data is not a multiple of the block size")
	}

	count := len(decrypted) / 8

	encrypted := make([]byte, len(decrypted))

	for i := 0; i < count; i++ {
		cipher.Encrypt(encrypted[i*8:], decrypted[i*8:])
	}

	return encrypted, nil
}
