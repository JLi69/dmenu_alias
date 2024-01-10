package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"sort"
)

// Takes a line of the form str1=str2 and returns (str1, str2)
// if the line is not properly formatted, return an error
// Example: vi=nvim will return (vi, nvim)
func parse_line(line string) (string, string, error) {
	alias := strings.Split(line, "=")

	if len(alias) == 2 {
		return strings.Replace(alias[0], "\n", "", -1),
			strings.Replace(alias[1], "\n", "", -1),
			nil
	}

	return "", "", errors.New("Syntax error: " + line)
}

// Parse the aliast list and return a map
func parse_alias_list(path string) map[string]string {
	alias_list := make(map[string]string)

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

		alias_list[s1] = s2
	}

	return alias_list
}

func read_dmenu_output(alias_list map[string] string) {
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

func output_dmenu_input(alias_list map[string] string) {	
	dmenu_input := make([]string, 0)
	set := make(map[string] bool)

	//Insert alias strings into the set and input
	for str := range alias_list {
		dmenu_input = append(dmenu_input, str)
		set[str] = true
	}
	
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
		_, exists := set[line]
		//Make sure we don't double insert an element
		if !exists {
			dmenu_input = append(dmenu_input, line)
		}
	}

	sort.Strings(dmenu_input)

	//Output input to stdout to be piped into dmenu
	for _, str := range dmenu_input {
		fmt.Println(str)
	}
}

const (
	output_mode int = iota
	input_mode
)

func get_path_from_args() string {
	for i, arg := range os.Args {
		//Ignore argument 0 (the executable name)
		if i == 0 {
			continue
		}

		if arg != "-o" && arg != "-i" {
			return arg
		}
	}

	//No arguments provided
	//Use the path $HOME/.config/dmenu_alias_list
	//We assume $HOME is defined
	return os.Getenv("HOME") + "/.config/dmenu_alias_list"
}

func read_options() int {
	for _, arg := range os.Args {
		if arg == "-o" {
			//-o means read output from dmenu
			return output_mode
		} else if arg == "-i" {
			//-i means generate input for dmenu
			return input_mode
		}
	}
	return output_mode
}

func main() {
	alias_list := parse_alias_list(get_path_from_args())	
	switch read_options() {
		case output_mode:
			read_dmenu_output(alias_list)
			break
		case input_mode:
			output_dmenu_input(alias_list)
			break
		default:
			break
	}
}
