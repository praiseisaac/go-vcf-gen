package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Replace special characters with numbers
func replaceSpecialCharacters(phone string) (string, error) {
	var result strings.Builder
	for _, char := range phone {
		if char >= '0' && char <= '9' {
			result.WriteRune(char)
		}
	}

	if result.String() == "" {
		return "", errors.New("no phone number found")
	}
	return result.String(), nil
}

func main() {
	totalContacts := 0

	if len(os.Args) < 2 {
		fmt.Println("Please provide a file name as an argument")
		return
	}
	fileName := os.Args[1]
	linesToSkip := 0
	if len(os.Args) > 2 {
		linesToSkip, _ = strconv.Atoi(os.Args[2])
	}

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Create cards directory if it doesn't exist
	if _, err := os.Stat("cards"); os.IsNotExist(err) {
		os.Mkdir("cards", 0755)
	}

	currentLine := 0
	
	for scanner.Scan() {
		line := scanner.Text()
		if currentLine < linesToSkip {
			currentLine++
			continue
		}

		// Get name
		name := strings.Split(line, ",")[0]
		nameParts := strings.Split(name, " ")
		lastName := nameParts[len(nameParts)-1]
		firstName := strings.Join(nameParts[:len(nameParts)-1], " ")

		// Get phone number	
		phone := strings.Split(line, ",")[1]
		formattedPhone, err := replaceSpecialCharacters(phone)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		// Create vcf file
		vcfFile := fmt.Sprintf("cards/%s.%s.vcf", strings.Split(line, ",")[0], formattedPhone)
		formattedPhone = "+1" + formattedPhone

		exportFile, err := os.Create(vcfFile)
		if err != nil {
			fmt.Println("Error creating file:", err)
			return
		}
		defer exportFile.Close()

		// Write vcf file
		exportFile.WriteString("BEGIN:VCARD\n\n")
		exportFile.WriteString("VERSION:3.0\n\n")
		exportFile.WriteString(fmt.Sprintf("N:%s;%s;;;\n\n", lastName, firstName))
		exportFile.WriteString(fmt.Sprintf("FN:%s\n\n", name))
		exportFile.WriteString(fmt.Sprintf("TEL;type=CELL;type=VOICE;type=pref:%s\n\n", formattedPhone))
		exportFile.WriteString("END:VCARD\n\n")
		exportFile.WriteString("\n")
		totalContacts++
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}

	fmt.Printf("Total contacts: %d\n", totalContacts)
}
