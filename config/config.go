package config

import (
	"encoding/json"
	"fmt" //used to print errors majorly.
	"os"  //it will be used to help us read our config.json file.
)

var (
	Token     string //To store value of Token from token.txt.
	BotPrefix string // To store value of BotPrefix from config.json.

	config *configStruct //To store value extracted from config.json.
)

type configStruct struct {
	BotPrefix string `json : "BotPrefix"`
}

func ReadConfig() error {
	fmt.Println("Reading config file...")
	file, err := os.ReadFile("./config.json") // ioutil package's ReadFile method which we read config.json and return it's value we will then store it in file variable and if an error ocurrs it will be stored in err .

	//Handling error and printing it using fmt package's Println function and returning it .
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// We are here printing value of file variable by explicitly converting it to string .

	fmt.Println(string(file))
	// Here we performing a simple task by copying value of file into config variable which we have declared above , and if there any error we are storing it in err . Unmarshal takes second arguments reference remember it .
	err = json.Unmarshal(file, &config)

	//Handling error
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	tokenFile, err := os.ReadFile("./token.txt")

	if err != nil {
		fmt.Printf("Unable to extract bot token, error: %s\n", err.Error())
		return err
	}

	// After storing value in config variable we will access it and storing it in our declared variables .
	Token = string(tokenFile)
	BotPrefix = config.BotPrefix

	//If there isn't any error we will return nil.
	return nil

}
