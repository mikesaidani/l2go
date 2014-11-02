package loginserver

import (
  "fmt"
  "net"
  "os"
  "crypto/rand"
  "crypto/rsa"
  _"code.google.com/p/go.crypto/blowfish"
)

var modulus []byte

func handleConnection(conn net.Conn, modulus []byte) {
  _content := []byte{0x00, 0xfd, 0x8a, 0x22, 0x00, 0x5a, 0x78, 0x00, 0x00, 0x0e, 0xea, 0x0b, 0xf3, 0x3a, 0x65, 0xc6, 0xc4, 0x62, 0xc7, 0x77, 0x2e, 0x95, 0xde, 0xbc, 0x8c, 0xe0, 0xf1, 0xc9, 0x87, 0xcb, 0x5f, 0xe5, 0x0e, 0x85, 0xa6, 0xf4, 0xac, 0x49, 0xb6, 0x29, 0xe3, 0xa5, 0x11, 0xbe, 0x85, 0x5d, 0x4c, 0x2a, 0x87, 0x0d, 0xd5, 0x17, 0x48, 0x87, 0x0a, 0xd4, 0xa8, 0x9b, 0x9b, 0x8b, 0x0f, 0xad, 0xa3, 0x4d, 0x60, 0x23, 0x6f, 0x2c, 0x53, 0xcc, 0xfb, 0x90, 0xea, 0xa2, 0x91, 0x24, 0x0e, 0x55, 0x6b, 0xb7, 0xb6, 0x6e, 0x30, 0x26, 0x7f, 0xf9, 0x49, 0xd8, 0xb2, 0x2a, 0x47, 0x17, 0xce, 0xd7, 0x10, 0xfc, 0x7d, 0x6f, 0xbc, 0x83, 0xb4, 0xd4, 0x53, 0x04, 0x6e, 0x08, 0x14, 0x7b, 0x92, 0xca, 0xb1, 0x52, 0x55, 0xf7, 0x45, 0x4c, 0xaa, 0xe9, 0xb0, 0x01, 0x1e, 0xac, 0xe2, 0x9b, 0x68, 0x21, 0x29, 0x68, 0x21, 0xe1, 0x93, 0x70, 0xbd, 0x3f, 0x13, 0x16, 0xab,0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

  for i := 0; i < len(modulus); i++ {
    _content[9+i] = modulus[i];
  }

  length := len(_content) + 2

  fmt.Println("A client is trying to connect...")

  buffer := make([]byte, 155)

  buffer[0] = byte(length & 0xff)
  buffer[1] = byte((length >> 8) & 0xff)
  copy(buffer[2:], _content)

  conn.Write([]byte(buffer))
  data := make([]byte, 1024)

  for {
    n, error := conn.Read(data)
    if error != nil {
      fmt.Printf("Error: Reading data : %s \n", error)
      conn.Close()
      break
    } else {
      fmt.Println(string(data[0:n]))
    }
  }
}

/*func blowfishDecrypt(et, key []byte) []byte {*/
	//// create the cipher
	//dcipher, err := blowfish.NewCipher(key)
	//if err != nil {
		//// fix this. its okay for this tester program, but...
		//panic(err)
	//}
	//// make initialisation vector to be the first 8 bytes of ciphertext.
	//// see related note in blowfishEncrypt()
	//div := et[:blowfish.BlockSize]
	//// check last slice of encrypted text, if it's not a modulus of cipher block size, we're in trouble
	//decrypted := et[blowfish.BlockSize:]
	//if len(decrypted)%blowfish.BlockSize != 0 {
		//panic("decrypted is not a multiple of blowfish.BlockSize")
	//}
	//// ok, we're good... create the decrypter
	//dcbc := dcipher.NewCBCDecrypter(dcipher, div)
	//// decrypt!
	//dcbc.CryptBlocks(decrypted, decrypted)
	//return decrypted
//}

//func blowfishEncrypt(ppt, key []byte) []byte {
	//// create the cipher
	//ecipher, err := blowfish.NewCipher(key)
	//if err != nil {
		//// fix this. its okay for this tester program, but ....
		//panic(err)
	//}
	//// make ciphertext big enough to store len(ppt)+blowfish.BlockSize
	//ciphertext := make([]byte, blowfish.BlockSize+len(ppt))
	//// make initialisation vector to be the first 8 bytes of ciphertext. you
	//// wouldn't do this normally/in real code, but this IS example code! :)
	//eiv := ciphertext[:blowfish.BlockSize]
	//// create the encrypter
	//ecbc := ecipher.NewCBCEncrypter(ecipher, eiv)
	//// encrypt the blocks, because block cipher
	//ecbc.CryptBlocks(ciphertext[blowfish.BlockSize:], ppt)
	//// return ciphertext to calling function
	//return ciphertext
/*}*/

func generateRSA() []byte {
  privatekey, err := rsa.GenerateKey(rand.Reader, 1024)

  if err != nil {
   fmt.Println(err.Error)
   os.Exit(1)
  }

  var publickey *rsa.PublicKey
  publickey = &privatekey.PublicKey
  scrambledModulus := publickey.N.Bytes() // modulus to bytes

  for i :=0 ; i < 4 ; i++ {
    temp := scrambledModulus[0x00+i];
    scrambledModulus[0x00+i] = scrambledModulus[0x4d+i];
    scrambledModulus[0x4d+i] = temp;
  }

  // step 2 xor first 0x40 bytes with last 0x40 bytes
  for i := 0; i < 0x40; i++ {
      scrambledModulus[i] = byte(scrambledModulus[i] ^ scrambledModulus[0x40+i]);
  }

  // step 3 xor bytes 0x0d-0x10 with bytes 0x34-0x38
  for i := 0; i < 4; i++ {
      scrambledModulus[0x0d+i] = byte(scrambledModulus[0x0d+i] ^ scrambledModulus[0x34+i]);
  }

  // step 4 xor last 0x40 bytes with first 0x40 bytes
  for i := 0; i < 0x40; i++ {
      scrambledModulus[0x40+i] = byte(scrambledModulus[0x40+i] ^ scrambledModulus[i]);
  }

  return scrambledModulus
}

func Init() {
  modulus = generateRSA()
  fmt.Println(modulus)

  ln, err := net.Listen("tcp", ":2106")
  defer ln.Close()

  if err != nil {
    fmt.Println("Couldn't initialize the Login Server")
  } else {
    fmt.Println("Login Server initialized.")
    fmt.Println("Listening on 127.0.0.1:2106.")
  }

  for {
    conn, err := ln.Accept()
    if err != nil {
      fmt.Println("Couldn't accept the incoming connection.")
      continue
    }

    go handleConnection(conn, modulus)
  }
}
