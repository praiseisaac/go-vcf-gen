package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
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
	var fileName string
	var linesToSkip int
	var org string

	flag.StringVar(&fileName, "file", "", "CSV file to process")
	flag.IntVar(&linesToSkip, "skip", 0, "Number of lines to skip")
	flag.StringVar(&org, "org", "", "Organization name")
	flag.Parse()

	if fileName == "" {
		fmt.Println("Please provide a file name using -file flag")
		return
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

		// Split row into components
		parts := strings.Split(line, ",")
		if len(parts) < 3 {
			fmt.Println("Invalid line format:", line)
			continue
		}

		// Get name
		name := parts[0]
		nameParts := strings.Split(name, " ")
		lastName := nameParts[len(nameParts)-1]
		firstName := strings.Join(nameParts[:len(nameParts)-1], " ")

		// Get phone number	
		phone := parts[1]
		formattedPhone, err := replaceSpecialCharacters(phone)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		// Get birthday
		birthday := parts[2]
		if birthday != "" {
			dateParts := strings.Split(birthday, "/")
			if len(dateParts) == 3 {
				birthday = dateParts[2] + "-" + fmt.Sprintf("%02s", dateParts[0]) + "-" + fmt.Sprintf("%02s", dateParts[1])
			}
		}

		// Create vcf file
		vcfFile := fmt.Sprintf("cards/%s.%s.vcf", parts[0], formattedPhone)
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
		if org != "" {
			exportFile.WriteString(fmt.Sprintf("ORG:%s;\n\n", org))
		}
		if birthday != "" {
			exportFile.WriteString(fmt.Sprintf("BDAY:%s\n\n", birthday))
		}
		exportFile.WriteString("END:VCARD\n\n")
		exportFile.WriteString("\n")
		totalContacts++
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}

	fmt.Printf("Total contacts: %d\n", totalContacts)
}
