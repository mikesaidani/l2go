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
      "port": "7777",

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
      }
    }    
  ]
}`

type ConfigObject struct {
	LoginServer LoginserverType
	GameServers []GameserverType
}

type GameServerConfigObject struct {
	LoginServer LoginserverType
	GameServer  GameserverType
}

type DatabaseType struct {
	Name     string
	Host     string
	Port     int
	User     string
	Password string
}

type CacheType struct {
	Index    int
	Host     string
	Port     int
	Password string
}

type LoginserverType struct {
	Host       string
	AutoCreate bool
	Database   DatabaseType
}

type GameserverType struct {
	Name       string
	InternalIP string
	ExternalIP string
	Port       string
	Database   DatabaseType
	Cache      CacheType
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
