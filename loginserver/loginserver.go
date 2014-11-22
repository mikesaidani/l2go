package loginserver

import (
	"fmt"
  "code.google.com/p/go.crypto/bcrypt"
	"github.com/frostwind/l2go/config"
	"github.com/frostwind/l2go/loginserver/clientpackets"
	"github.com/frostwind/l2go/loginserver/models"
	"github.com/frostwind/l2go/loginserver/packet"
	"github.com/frostwind/l2go/loginserver/serverpackets"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net"
	"strconv"
)

type LoginServer struct {
	clients            []*models.Client
	database           *mgo.Database
	config             config.ConfigObject
	internalServerList []byte
	externalServerList []byte
	status             loginServerStatus
	databaseSession    *mgo.Session
	socket             net.Listener
}

type loginServerStatus struct {
	SuccessfulLoginAttempts uint32
	FailedLoginAttempts     uint32
	HackAttempts            uint32
}

func New(cfg config.ConfigObject) *LoginServer {
	return &LoginServer{config: cfg}
}

func (l *LoginServer) Init() {
	var err error

	// Connect to our database
	l.databaseSession, err = mgo.Dial(l.config.LoginServer.Database.Host + ":" + strconv.Itoa(l.config.LoginServer.Database.Port))
	if err != nil {
		panic("Couldn't connect to the database server")
	} else {
		fmt.Println("Successfully connect to the database server")
	}

	// Select the appropriate database
	l.database = l.databaseSession.DB(l.config.LoginServer.Database.Name)

	l.socket, err = net.Listen("tcp", ":2106")
	if err != nil {
		fmt.Println("Couldn't initialize the Login Server")
	} else {
		fmt.Println("Login Server listening on port 2106")
	}
}

func (l *LoginServer) Start() {
	defer l.databaseSession.Close()
	defer l.socket.Close()

	for {
		var err error
		client := models.NewClient()
		l.clients = append(l.clients, client)
		client.Socket, err = l.socket.Accept()
		if err != nil {
			fmt.Println("Couldn't accept the incoming connection.")
			continue
		} else {
			go l.handleClientPackets(client)
		}
	}
}

func (l *LoginServer) handleClientPackets(client *models.Client) {
	fmt.Println("A client is trying to connect...")
	defer client.Socket.Close()

	buffer := serverpackets.NewInitPacket()
	err := packet.Send(client.Socket, buffer, false, false)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Init packet sent.")
	}

	for {
		p, err := packet.Receive(client.Socket)

		if err != nil {
			fmt.Println(err)
			fmt.Println("Closing the connection...")
			break
		}

		switch opcode := p.GetOpcode(); opcode {
		case 00:
			// response buffer
			var buffer []byte

			requestAuthLogin := clientpackets.NewRequestAuthLogin(p.GetData())

			fmt.Printf("User %s is trying to login\n", requestAuthLogin.Username)

			accounts := l.database.C("accounts")
			err := accounts.Find(bson.M{"username": requestAuthLogin.Username}).One(&client.Account)

			if err != nil {
				if l.config.LoginServer.AutoCreate == true {
					hashedPassword, err := bcrypt.GenerateFromPassword([]byte(requestAuthLogin.Password), 10)
					if err != nil {
						fmt.Println("An error occured while trying to generate the password")

						buffer = serverpackets.NewLoginFailPacket(serverpackets.REASON_SYSTEM_ERROR)
					} else {
						err = accounts.Insert(&models.Account{requestAuthLogin.Username, string(hashedPassword), ACCESS_LEVEL_PLAYER})
						if err != nil {
							fmt.Printf("Couldn't create an account for the user %s\n", requestAuthLogin.Username)

							buffer = serverpackets.NewLoginFailPacket(serverpackets.REASON_SYSTEM_ERROR)
						} else {
							fmt.Printf("Account successfully created for the user %s\n", requestAuthLogin.Username)
							buffer = serverpackets.NewLoginOkPacket(client.SessionID)
						}
					}
				} else {
					fmt.Println("Account not found !")
					buffer = serverpackets.NewLoginFailPacket(serverpackets.REASON_USER_OR_PASS_WRONG)
				}
			} else {
				// Account exists; Is the password ok?
				err = bcrypt.CompareHashAndPassword([]byte(client.Account.Password), []byte(requestAuthLogin.Password))

				if err != nil {
					fmt.Printf("Wrong password for the account %s\n", requestAuthLogin.Username)

					buffer = serverpackets.NewLoginFailPacket(serverpackets.REASON_USER_OR_PASS_WRONG)
				} else {

					if client.Account.AccessLevel >= ACCESS_LEVEL_PLAYER {
						buffer = serverpackets.NewLoginOkPacket(client.SessionID)
					} else {
						buffer = serverpackets.NewLoginFailPacket(serverpackets.REASON_ACCESS_FAILED)
					}

				}
			}

			err = packet.Send(client.Socket, buffer)

			if err != nil {
				fmt.Println(err)
			}

		case 02:
			requestPlay := clientpackets.NewRequestPlay(p.GetData())

			fmt.Printf("The client wants to connect to the server : %d\n", requestPlay.ServerID)

			var buffer []byte
			if len(l.config.GameServers) >= int(requestPlay.ServerID) && (l.config.GameServers[requestPlay.ServerID-1].Options.Testing == false || client.Account.AccessLevel > ACCESS_LEVEL_PLAYER) {
				buffer = serverpackets.NewPlayOkPacket()
			} else {
				buffer = serverpackets.NewPlayFailPacket(serverpackets.REASON_ACCESS_FAILED)
			}
			err := packet.Send(client.Socket, buffer)

			if err != nil {
				fmt.Println(err)
			}

		case 05:
			buffer := serverpackets.NewServerListPacket(l.config.GameServers, client.Socket.RemoteAddr().String())
			err := packet.Send(client.Socket, buffer)

			if err != nil {
				fmt.Println(err)
			}

		default:
			fmt.Println("Couldn't detect the packet type.")
		}
	}

}
