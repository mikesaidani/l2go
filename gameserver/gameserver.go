package gameserver

import (
	"bytes"
	"fmt"
	"github.com/frostwind/l2go/config"
	"github.com/frostwind/l2go/gameserver/clientpackets"
	"github.com/frostwind/l2go/gameserver/models"
	"github.com/frostwind/l2go/gameserver/serverpackets"
	"gopkg.in/mgo.v2"
	"net"
	"strconv"
)

type GameServer struct {
	clients         []*models.Client
	database        *mgo.Database
	config          config.GameServerConfigObject
	status          gameServerStatus
	databaseSession *mgo.Session
	socket          net.Listener
}

type gameServerStatus struct {
	onlinePlayers uint32
	hackAttempts  uint32
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

	// Listen for client connections
	g.socket, err = net.Listen("tcp", ":"+strconv.Itoa(g.config.GameServer.Port))
	if err != nil {
		fmt.Println("Couldn't initialize the Game Server")
	} else {
		fmt.Printf("Game Server listening on port %s\n", strconv.Itoa(g.config.GameServer.Port))
	}
}

func (g *GameServer) Start() {
	defer g.databaseSession.Close()
	defer g.socket.Close()

	for {
		var err error
		client := models.NewClient()
		client.Socket, err = g.socket.Accept()
		g.clients = append(g.clients, client)
		if err != nil {
			fmt.Println("Couldn't accept the incoming connection.")
			continue
		} else {
			go g.handleClientPackets(client)
		}
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
