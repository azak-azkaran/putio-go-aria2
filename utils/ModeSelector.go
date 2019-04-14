package utils

import (
	"bufio"
	"errors"
	"os"
)

func GetArguments() (string, string, string, error) {
	var configfile string
	var mode string
	argsWithProg := os.Args
	i := len(argsWithProg)

	switch i {
	case 2:
		configfile = argsWithProg[1]
		Info.Println("Mode not select Download(d) or Organize(o)")
		reader := bufio.NewReader(os.Stdin)
		mode, err := reader.ReadString('\n')
		if err != nil {
			Error.Fatalln("Error while reading Input")
			return "", "", "", errors.New("Error while reading")
		}

		Info.Println("reading config file: ", configfile)
		Info.Println("Starting with Mode: ", mode)
		return configfile, mode, "", nil
	case 3:
		configfile = argsWithProg[1]
		mode = argsWithProg[2]

		Info.Println("reading config file: ", configfile)
		Info.Println("Starting with Mode: ", mode)
		return configfile, mode, "", nil
	case 4:
		configfile = argsWithProg[1]
		mode = argsWithProg[2]
		filter := argsWithProg[3]
		Info.Println("reading config file: ", configfile)
		Info.Println("Starting with Mode: ", mode)

		return configfile, mode, filter, nil

	default:
		Error.Fatalln("script was used wrong:\n putio-go-aria2 oauth.secret\nputio-go-aria2 oauth.secret filter")
		return "", "", "", errors.New("Wrong usage")
	}
}
