package loginserver

import (
	"code.google.com/p/go.crypto/bcrypt"
	"fmt"
	"github.com/frostwind/l2go/config"
	"github.com/frostwind/l2go/loginserver/clientpackets"
	"github.com/frostwind/l2go/loginserver/packet"
	"github.com/frostwind/l2go/loginserver/serverpackets"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net"
)

type Account struct {
	Username    string `bson:username`
	Password    string `bson:password`
	AccessLevel uint8  `bson:access_level`
}

func handleConnection(conn net.Conn, conf config.ConfigObject, session *mgo.Session) {

	fmt.Println("A client is trying to connect...")

	buffer := serverpackets.NewInitPacket()
	err := packet.Send(conn, buffer, false, false)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Init packet sent.")
	}

	for {
		p, err := packet.Receive(conn)

		if err != nil {
			fmt.Println(err)
			fmt.Println("Closing the connection...")
			conn.Close()
			break
		}

		switch opcode := p.GetOpcode(); opcode {
		case 00:

			// reponse buffer
			var buffer []byte

			requestAuthLogin := clientpackets.NewRequestAuthLogin(p.GetData())

			fmt.Printf("User %s is trying to login\n", requestAuthLogin.Username)

			accounts := session.DB(conf.LoginServer.Database.Name).C("accounts")

			account := Account{}
			err := accounts.Find(bson.M{"username": requestAuthLogin.Username}).One(&account)

			if err != nil {
				fmt.Println("Account not found !")

				if conf.LoginServer.AutoCreate == true {
					hashedPassword, err := bcrypt.GenerateFromPassword([]byte(requestAuthLogin.Password), 10)
					if err != nil {
						fmt.Println("An error occured while trying to generate the password")

						buffer = serverpackets.NewLoginFailPacket(serverpackets.REASON_SYSTEM_ERROR)
					} else {
						err = accounts.Insert(&Account{requestAuthLogin.Username, string(hashedPassword), 50})
						if err != nil {
							fmt.Printf("Couldn't create an account for the user %s\n", requestAuthLogin.Username)

							buffer = serverpackets.NewLoginFailPacket(serverpackets.REASON_SYSTEM_ERROR)
						} else {
							fmt.Printf("Account successfully created for the user %s\n", requestAuthLogin.Username)
							buffer = serverpackets.NewLoginOkPacket()
						}
					}

				}
			} else {
				// Account exists; Is the password ok?
				err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(requestAuthLogin.Password))

				if err != nil {
					fmt.Printf("Wrong password for the account %s\n", requestAuthLogin.Username)

					buffer = serverpackets.NewLoginFailPacket(serverpackets.REASON_USER_OR_PASS_WRONG)
				} else {
					buffer = serverpackets.NewLoginOkPacket()

				}
			}

			err = packet.Send(conn, buffer)

			if err != nil {
				fmt.Println(err)
			}

		case 02:
			serverId := p.GetData()[4+4+1] // Skip the sessionId (2*4bytes) and grab the serverId

			fmt.Printf("The client wants to connect to the server : %X\n", serverId)

			buffer := serverpackets.NewPlayOkPacket()
			err := packet.Send(conn, buffer)

			if err != nil {
				fmt.Println(err)
			}

		case 05:
			buffer := serverpackets.NewServerListPacket()
			err := packet.Send(conn, buffer)

			if err != nil {
				fmt.Println(err)
			}

		default:
			fmt.Println("Couldn't detect the packet type.")
		}
	}

}

func Init(conf config.ConfigObject) {

	// Database setup
	session, err := mgo.Dial(conf.LoginServer.Database.Host + ":" + conf.LoginServer.Database.Port)
	if err != nil {
		panic("Couldn't connect to the database server")
	} else {
		fmt.Println("Successfully connect to the database server")
	}
	defer session.Close()

	// Socket Setup
	ln, err := net.Listen("tcp", ":2106")
	defer ln.Close()

	if err != nil {
		fmt.Println("Couldn't initialize the Login Server")
	} else {
		fmt.Println("Login Server listening on port 2106")
	}

	// Connections handling
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Couldn't accept the incoming connection.")
			continue
		} else {
			go handleConnection(conn, conf, session)
		}
	}
}
