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
  fmt.Printf("Header : %X\n", p.header)
  fmt.Printf("Data : %X\n", p.data)

  decrypted := blowfishDecrypt(p.data, []byte("[;'.]94-31==-%&@!^+]\000"), size-2)
  fmt.Printf("Decrypted packet content : %X\n", decrypted)

  decrypted2 := blowfishDecrypt(p.data, []byte("_;5.]94-31==-%xT!^[$\000"), size-2)
  fmt.Printf("Decrypted packet content : %X\n", decrypted2)

  decrypted3 := blowfishDecrypt(p.data, []byte("[;'.]94-&@%!^+]-31==\000"), size-2)
  fmt.Printf("Decrypted packet content : %X\n", decrypted3)

  decrypted4 := blowfishDecrypt(p.data, []byte("31==-%&@!^+][;'.]94-\000"), size-2)
  fmt.Printf("Decrypted packet content : %X\n", decrypted4)

  decrypted5 := blowfishDecrypt(p.data, []byte("_;V.]05-31!|+-%XT!^[\000"), size-2)
  fmt.Printf("Decrypted packet content : %X\n", decrypted5)
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
