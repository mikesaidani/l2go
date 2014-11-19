package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/user"
)

var defaultServerConfig = `{
  "loginserver": {
    "host": "127.0.0.1",
    "autoCreate": true,
    "database": {
      "name": "l2go-login",
      "host": "127.0.0.1",
      "port": 27017,
      "user": "",
      "password": ""
    } 
  },

  "gameservers": [
    {
      "name": "Bartz",
      "secret": "CHANGE_ME_PLEASE",
      "internalIP": "127.0.0.1",
      "externalIP": "192.168.1.2",
      "port": 7777,

      "database": {
        "name": "l2go-server",
        "host": "127.0.0.1",
        "port": 27017,
        "user": "",
        "password": ""
      },

      "cache": {
        "host": "127.0.0.1",
        "port": 6379,
        "password": ""
      },

      "options": {
        "maxPlayers": 10000,
        "testing": false
      }
    }    
  ]
}`

type ConfigObject struct {
	LoginServer LoginServerType
	GameServers []GameServerType
}

type GameServerConfigObject struct {
	LoginServer LoginServerType
	GameServer  GameServerType
}

type DatabaseType struct {
	Name     string
	Host     string
	Port     int
	User     string
	Password string
}

type CacheType struct {
	Host     string
	Port     int
	Password string
}

type LoginServerType struct {
	Host       string
	AutoCreate bool
	Database   DatabaseType
}

type GameServerType struct {
	Name       string
	InternalIP string
	ExternalIP string
	Port       int
	Database   DatabaseType
	Cache      CacheType
	Options    OptionsType
}

type OptionsType struct {
	MaxPlayers uint16
	Testing    bool
}

func Read() ConfigObject {
	usr, _ := user.Current()
	dir := usr.HomeDir

	var jsontype ConfigObject
	file, e := ioutil.ReadFile(dir + "/.l2go/config/server.json")

	if e != nil {
		fmt.Println("Couldn't load the server configuration file. Using the default preset.")
		json.Unmarshal([]byte(defaultServerConfig), &jsontype)
	} else {
		json.Unmarshal(file, &jsontype)
	}

	return jsontype
}
