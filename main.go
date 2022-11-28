package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	var words []string

	args := os.Args[1:]

	// 2 file names or quit program
	if len(args) != 2 {
		fmt.Println("invalid arguments")
		return
	}

	file, err := os.Open(args[0])

	// wrong file name = error
	if err != nil {
		fmt.Println(err)
		return

	} else {

		// make it string array
		makeByte, _ := file.Stat()
		array := make([]byte, makeByte.Size())
		file.Read(array)
		words = strings.Fields(string(array))
		file.Close()
	}

	// ' where it supposed to be
	count := 0
	for i, word := range words {
		if word == "'" && count == 0 {
			count += 1
			words[i+1] = word + words[i+1]
			words = append(words[:i], words[i+1:]...)
		}
	}

	// space before ,
	for i := range words {
		if words[i][0] == ',' && len(words[i]) > 1 {
			words[i-1] = words[i-1] + string(words[i][0])
			words[i] = words[i][1:]
			if len(words[i]) == 0 {
				remove(words, i)
				i--
			}
		}
	}

	// find modifiers in the text and change stuff
	for i := 0; i < len(words); i++ {
		switch words[i] {

		case "(up)":
			words[i-1] = strings.ToUpper(words[i-1])
			for j := 0; j <= i; j++ {
				words[i] = strings.ToUpper(words[i-j])
			}
			remove(words, i)
			i--

		case "(low)":
			words[i-1] = strings.ToLower(words[i-1])
			for j := 0; j <= i; j++ {
				words[i] = strings.ToLower(words[i-j])
			}
			remove(words, i)
			i--

		case "(cap)":
			words[i-1] = strings.Title(words[i-1])
			for j := 0; j <= i; j++ {
				words[i] = strings.Title(words[i-j])
			}
			remove(words, i)
			i--

		case "(up,":
			next := words[i+1]
			count, _ := strconv.Atoi(strings.TrimSuffix(next, ")"))
			for j := 0; j <= count; j++ {
				words[i-j] = strings.ToUpper(words[i-j])
			}
			remove(words, i)
			remove(words, i)
			i--

		case "(low,":
			next := words[i+1]
			count, _ := strconv.Atoi(strings.TrimSuffix(next, ")"))
			for j := 0; j <= count; j++ {
				words[i-j] = strings.ToLower(words[i-j])
			}
			remove(words, i)
			remove(words, i)
			i--

		case "(cap,":
			next := words[i+1]
			count, _ := strconv.Atoi(strings.TrimSuffix(next, ")"))
			for j := 0; j <= count; j++ {
				words[i-j] = strings.Title(words[i-j])
			}
			remove(words, i)
			remove(words, i)
			i--

		case "a":
			next := words[i+1]
			nextLetter := next[:1]
			switch nextLetter {
			case "a", "e", "i", "o", "u", "h", "A", "E", "I", "O", "U", "H":
				words[i] = "an"
			}

		case "(hex)":
			num, err4 := strconv.ParseInt(words[i-1], 16, 64)
			if err4 != nil {
				fmt.Println(err4)
			}
			words[i-1] = fmt.Sprint(num)
			remove(words, i)

		case "(bin)":
			num, err4 := strconv.ParseInt(words[i-1], 2, 64)
			if err4 != nil {
				fmt.Println(err4)
			}
			words[i-1] = fmt.Sprint(num)
			remove(words, i)

		case ".", ",", "!", "?", ":", ";", "'":
			words[i-1] = words[i-1] + string(words[i][0])
			words[i] = words[i][1:]
			if len(words[i]) == 0 {
				remove(words, i)
				i--
			}
		}
	}

	// making new file for result
	result, err2 := os.Create(args[1])
	if err2 != nil {
		fmt.Println(err2)
	}

	// join and add spaces, remove spaces
	String := strings.Join(words, " ")
	String = strings.TrimRight(String, " ")

	// Write string to the output file
	toWrite := String
	_, err3 := result.Write([]byte(toWrite))
	if err3 != nil {
		fmt.Println(err3)
		os.Exit(1)
	}
	result.Close()
}

// remove stuff from text (cap, low, etc...)
func remove(array []string, i int) []string {
	copy(array[i:], array[i+1:])
	array[len(array)-1] = ""
	array = array[:len(array)-1]
	return array
}
