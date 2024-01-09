package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

//Takes a line of the form str1=str2 and returns (str1, str2)
//if the line is not properly formatted, return an error
//Example: vi=nvim will return (vi, nvim)
func parse_line(line string) (string, string, error) {
	alias := strings.Split(line, "=")

	if len(alias) == 2 {
		return strings.Replace(alias[0], "\n", "", -1),
			strings.Replace(alias[1], "\n", "", -1),
			nil
	}

	return "", "", errors.New("Syntax error: " + line)
}

//Parse the aliast list and return a map
func parse_alias_list(path string) map[string] string {
	alias_list := make(map[string] string)

	//Attempt to open the file
	file, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: file not found: %s\n", path)
		return alias_list
	}
	defer file.Close()

	//Read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		s1, s2, err := parse_line(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			continue
		}

		alias_list[s1] = s2;
	}

	return alias_list
}

func main() {
	alias_list_path := ""	
	
	if len(os.Args) == 1 {
		//No arguments provided	
		//Use the path $HOME/.config/dmenu_alias_list
		//We assume $HOME is defined
		home_path := os.Getenv("HOME") + "/.config/dmenu_alias_list"
		alias_list_path = home_path
	} else if len(os.Args) > 1 {
		alias_list_path = os.Args[1]
	}

	alias_list := parse_alias_list(alias_list_path)

	//Read stdin line by line and return the aliased string from the input
	reader := bufio.NewReader(os.Stdin)

	for {
		line, err := reader.ReadString('\n')
	
		//Empty line, end the program
		if len(line) == 0 {
			break
		}

		//Something went wrong with reading stdin, panic
		if err != nil {
			panic(err)
		}
		
		//Strip out the newlines from the line so we get the proper output
		line = strings.Replace(line, "\n", "", -1)
		str, ok := alias_list[line]
		if ok {
			fmt.Println(str)
		} else {
			fmt.Println(line)
		}
	}
}
