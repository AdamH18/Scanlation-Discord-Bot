package config

import (
	"encoding/json" //used to print errors majorly.
	"log"
	"os" //it will be used to help us read our config.json file.
	"strconv"
)

var (
	Token          string //To store value of Token.
	RemoveCommands bool   //Choose whether or not to remove registered slash commands from Discord upon shutdown.

	config *configStruct //To store value extracted from config.json.
)

type configStruct struct {
	Token          string `json: "Token"`
	RemoveCommands string `json: "RemoveCommands"`
}

func ReadConfig() error {
	log.Println("Reading config file...")
	file, err := os.ReadFile("./config.json") // ioutil package's ReadFile method which we read config.json and return it's value we will then store it in file variable and if an error ocurrs it will be stored in err .

	//Handling error and printing it using fmt package's Println function and returning it .
	if err != nil {
		log.Println(err.Error())
		return err
	}

	// We are here printing value of file variable by explicitly converting it to string .

	log.Println(string(file))
	// Here we performing a simple task by copying value of file into config variable which we have declared above , and if there any error we are storing it in err . Unmarshal takes second arguments reference remember it .
	err = json.Unmarshal(file, &config)

	//Handling error
	if err != nil {
		log.Println(err.Error())
		return err
	}

	// After storing value in config variable we will access it and storing it in our declared variables .
	Token = config.Token
	RemoveCommands, _ = strconv.ParseBool(config.RemoveCommands)

	//If there isn't any error we will return nil.
	return nil

}
