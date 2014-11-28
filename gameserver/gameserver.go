package gameserver

import (
  "errors"
	"bytes"
	"fmt"
	"github.com/frostwind/l2go/config"
	"github.com/frostwind/l2go/packets"
	"github.com/frostwind/l2go/gameserver/clientpackets"
	"github.com/frostwind/l2go/gameserver/models"
	"github.com/frostwind/l2go/gameserver/serverpackets"
	"gopkg.in/mgo.v2"
	"net"
	"strconv"
)

type GameServer struct {
	clients           []*models.Client
	database          *mgo.Database
	config            config.GameServerConfigObject
	status            gameServerStatus
	databaseSession   *mgo.Session
	clientListener    net.Listener
	loginServerSocket net.Conn
}

type gameServerStatus struct {
	onlinePlayers uint32
	hackAttempts  uint32
}

func (g *GameServer) Receive() (opcode byte, data []byte, e error) {
	// Read the first two bytes to define the packet size
	header := make([]byte, 2)
	n, err := g.loginServerSocket.Read(header)

	if n != 2 || err != nil {
		return 0x00, nil, errors.New("An error occured while reading the packet header.")
	}

	// Calculate the packet size
	size := 0
	size = size + int(header[0])
	size = size + int(header[1])*256

	// Allocate the appropriate size for our data (size - 2 bytes used for the length
	data = make([]byte, size-2)

	// Read the encrypted part of the packet
	n, err = g.loginServerSocket.Read(data)

	if n != size-2 || err != nil {
		return 0x00, nil, errors.New("An error occured while reading the packet data.")
	}

	// Print the raw packet
	fmt.Printf("Raw packet : %X%X\n", header, data)

	// Extract the op code
	opcode = data[0]
	data = data[1:]
	e = nil
	return
}

func (g *GameServer) Send(data []byte) error {
	// Calculate the packet length
	length := uint16(len(data) + 2)

	// Put everything together
	buffer := packets.NewBuffer()
	buffer.WriteUInt16(length)
	buffer.Write(data)

	_, err := g.loginServerSocket.Write(buffer.Bytes())

	if err != nil {
		return errors.New("The packet couldn't be sent.")
	}

	return nil
}

func New(cfg config.GameServerConfigObject) *GameServer {
	return &GameServer{config: cfg}
}

func (g *GameServer) Init() {
	var err error

	// Connect to our database
	g.databaseSession, err = mgo.Dial(g.config.GameServer.Database.Host + ":" + strconv.Itoa(g.config.GameServer.Database.Port))
	if err != nil {
		panic("Couldn't connect to the database server")
	} else {
		fmt.Println("Successfully connected to the database server")
	}

	// Select the appropriate database
	g.database = g.databaseSession.DB(g.config.GameServer.Database.Name)

	// Connect to the login server
	g.loginServerSocket, err = net.Dial("tcp", g.config.LoginServer.Host+":9413")
	if err != nil {
		fmt.Println("Couldn't connect to the Login Server")
	} else {
		fmt.Printf("Successfully connected to the Login Server at %s:9413\n", g.config.LoginServer.Host)
	}

	// Listen for client connections
	g.clientListener, err = net.Listen("tcp", ":"+strconv.Itoa(g.config.GameServer.Port))
	if err != nil {
		fmt.Println("Couldn't initialize the Game Server")
	} else {
		fmt.Printf("Game Server listening on port %s\n", strconv.Itoa(g.config.GameServer.Port))
	}
}

func (g *GameServer) Start() {
	defer g.databaseSession.Close()
	defer g.clientListener.Close()

	done := make(chan bool)

	go func() {
    g.Send([]byte{00, 01, 02})

    for {
      opcode, _, err := g.Receive()

      if err != nil {
        fmt.Println(err)
        fmt.Println("Closing the connection...")
        break
      }

      switch opcode {
      case 00:
        fmt.Println("A game server sent a request to register")
      default:
        fmt.Println("Can't recognize the packet sent by the gameserver")
      }
    }
		done <- true
	}()

	go func() {
		for {
			var err error
			client := models.NewClient()
			client.Socket, err = g.clientListener.Accept()
			g.clients = append(g.clients, client)
			if err != nil {
				fmt.Println("Couldn't accept the incoming connection.")
				continue
			} else {
				go g.handleClientPackets(client)
			}
		}

		done <- true
	}()

	for i := 0; i < 2; i++ {
		<-done
	}
}

func (g *GameServer) kickClient(client *models.Client) {
	client.Socket.Close()

	for i, item := range g.clients {
		if bytes.Equal(item.SessionID, client.SessionID) {
			copy(g.clients[i:], g.clients[i+1:])
			g.clients[len(g.clients)-1] = nil
			g.clients = g.clients[:len(g.clients)-1]
			break
		}
	}

	fmt.Println("The client has been successfully kicked from the server.")
}

func (g *GameServer) handleClientPackets(client *models.Client) {
	fmt.Println("A client is trying to connect...")
	defer g.kickClient(client)

	// Client protocol version
	_, data, err := client.Receive(false)
	protocolVersion := clientpackets.NewProtocolVersion(data)

	if err != nil {
		fmt.Println(err)
		fmt.Println("Closing the connection...")
		return
	}

	if protocolVersion.Version < 419 {
		fmt.Println("Wrong protocol version ! <Expected 419> <Got: %d>", protocolVersion.Version)
		return
	}

	fmt.Println("Sending the Xor Key to the client...")

	buffer := serverpackets.NewCryptInitPacket()
	err = client.Send(buffer, false)

	if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println("CryptInit packet sent.")
	}

	for {
		opcode, data, err := client.Receive()

		if err != nil {
			fmt.Println(err)
			fmt.Println("Closing the connection...")
			break
		}

		switch opcode {
		case 0x08:
			fmt.Println("Client is requesting login to the Game Server")

			buffer := serverpackets.NewCharListPacket()
			err := client.Send(buffer)

			if err != nil {
				fmt.Println(err)
			}

		case 0x0e:
			fmt.Println("Client is requesting character creation template")

			buffer := serverpackets.NewCharTemplatePacket()
			err := client.Send(buffer)

			if err != nil {
				fmt.Println(err)
			}

		case 0x0b:
			character := clientpackets.NewCharacterCreate(data)

			fmt.Printf("Created a new character : %s\n", character.Name)

			// ACK
			buffer := serverpackets.NewCharCreateOkPacket()
			err := client.Send(buffer)

			if err != nil {
				fmt.Println(err)
			}

			// Return to the character select screen
			buffer = serverpackets.NewCharListPacket()
			err = client.Send(buffer)

			if err != nil {
				fmt.Println(err)
			}

		default:
			fmt.Println("Couldn't detect the packet type.")
		}
	}

}
