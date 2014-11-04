package packet

import (
  "fmt"
  "code.google.com/p/go.crypto/blowfish"
)

type packet struct {
  header [2]byte
  data []byte

}

type packetError struct {
	message string
}

func (e packetError) Error() string {
	return fmt.Sprintf("%v", e.message)
}

func Decrypt(buffer []byte) (*packet, error) {

  // Init our packet struct
  p := new(packet)

  // First 2 bytes are the size
  p.header[0] = buffer[0]
  p.header[1] = buffer[1]

  // Calculate the packet size
  size := 0
  size = size + int(p.header[0])
  size = size + (int(p.header[1])*256)

  // Copy the packet body to data
  p.data = buffer[2:size]

  if len(p.data) != size-2 {
    return &packet{}, packetError{"Wrong packet size detected !"}
  }

  fmt.Printf("Raw packet : %X%X\n", p.header, p.data)

  decrypted := blowfishDecrypt(p.data, []byte("[;'.]94-31==-&%@!^+]"), size-2)
  fmt.Printf("Decrypted packet content : %x\n", decrypted[0])

  return p, nil
}

func blowfishDecrypt(encrypted, key []byte, size int) []byte {
  // create the cipher
  dcipher, err := blowfish.NewCipher(key)
  if err != nil {
    // fix this. its okay for this tester program, but...
    panic(err)
  }

  count := len(encrypted) / 8;

  decrypted := make([]byte, size)

  for i := 0; i < count; i++ {
    dcipher.Decrypt(decrypted[i*8:], encrypted[i*8:]);
  }

  return decrypted
}
