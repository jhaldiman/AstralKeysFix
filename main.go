package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

func writeLines(lines []string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}

func doesFileOrDirectoryExist(path string) bool {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return true
	}

	return false
}

func readUserInput() (string, error) {
	scanner := bufio.NewScanner(os.Stdin)
	var input string
	if scanner.Scan() {
		input = scanner.Text()
	}

	return input, scanner.Err()
}

func getAstralKeysPath() (string, error) {
	var wowDirectory string

	if doesFileOrDirectoryExist(filepath.Join(os.Getenv("PROGRAMFILES(X86)"), "World of Warcraft")) {
		wowDirectory = filepath.Join(os.Getenv("PROGRAMFILES(X86)"), "World of Warcraft")
	} else if doesFileOrDirectoryExist(filepath.Join(os.Getenv("PROGRAMFILES"), "World of Warcraft")) {
		wowDirectory = filepath.Join(os.Getenv("PROGRAMFILES"), "World of Warcraft")
	} else if doesFileOrDirectoryExist("D:\\Program Files (x86)\\World of Warcraft") {
		wowDirectory = "D:\\Program Files (x86)\\World of Warcraft"
	} else {
		fmt.Println("Unable to Automatically Detect World of Warcraft Installation")
		fmt.Println("Please Provide Path to the World of Warcraft Directory: ")
		var err error
		wowDirectory, err = readUserInput()
		if err != nil {
			return "", fmt.Errorf("WoW Directory Prompt: %s", err)
		}

		if !doesFileOrDirectoryExist(wowDirectory) {
			return "", errors.New("provided directory does not exist")
		}
	}

	fmt.Println("Using World of Warcraft Installation: " + wowDirectory)

	if !doesFileOrDirectoryExist(filepath.Join(wowDirectory, "_retail_")) {
		return "", errors.New("unable to locate retail installation inside World of Warcraft directory")
	}

	if !doesFileOrDirectoryExist(filepath.Join(wowDirectory, "_retail_/Interface/AddOns")) {
		return "", errors.New("unable to locate AddOns directory")
	}

	if !doesFileOrDirectoryExist(filepath.Join(wowDirectory, "_retail_/Interface/AddOns/AstralKeys")) {
		return "", errors.New("unable to locate Astral Keys AddOn")
	}

	fmt.Println("Found Astral Keys AddOn at: " + filepath.Join(wowDirectory, "_retail_/Interface/AddOns/AstralKeys"))

	return filepath.Join(wowDirectory, "_retail_/Interface/AddOns/AstralKeys"), nil
}

func exitPrompt(exitCode int) {
	fmt.Println()
	fmt.Print("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	os.Exit(exitCode)
}

func fatalErrorExitPrompt(message string, err error) {
	log.Printf(message, err)
	exitPrompt(1)
}

func fatalExitPrompt(message string) {
	log.Print(message)
	exitPrompt(1)
}

func main() {
	astralKeysPath, err := getAstralKeysPath()
	if err != nil {
		fatalErrorExitPrompt("Get Astral Keys Path: %s", err)
	}

	communicationsLua := filepath.Join(astralKeysPath, "Communications.lua")
	if !doesFileOrDirectoryExist(communicationsLua) {
		fatalExitPrompt("Unable to find Communications.lua")
	}

	var sourcePath string
	backupPath := filepath.Join(astralKeysPath, "Communications.lua_ORIGINAL")
	if doesFileOrDirectoryExist(backupPath) {
		sourcePath = backupPath
	} else {
		sourcePath = communicationsLua
	}

	fmt.Println("Please Enter the New Phrase to Use: ")
	newKeyTitle, err := readUserInput()
	if err != nil {
		fatalErrorExitPrompt("New Key Title: %s", err)
	}
	newKeyTitle = strings.ReplaceAll(newKeyTitle, "\\", "\\\\")
	newKeyTitle = strings.ReplaceAll(newKeyTitle, "'", "\\'")
	newKeyTitle = strings.ReplaceAll(newKeyTitle, "[", "\\[")
	newKeyTitle = strings.ReplaceAll(newKeyTitle, "]", "\\]")

	fmt.Println("Reading File: " + sourcePath)
	lines, err := readLines(sourcePath)
	if err != nil {
		fatalErrorExitPrompt("Open File: %s", err)
	}

	if communicationsLua == sourcePath {
		fmt.Println("Creating backup at: " + backupPath)
		if err := writeLines(lines, backupPath); err != nil {
			fatalErrorExitPrompt("Backup File: %s", err)
		}
	}

	var outLines []string
	for _, line := range lines {
		if strings.Contains(line, "Astral Keys") {
			outLines = append(outLines, strings.Replace(line, "Astral Keys", newKeyTitle, 1))
		} else {
			outLines = append(outLines, line)
		}
	}

	fmt.Println("Updating Communications.lua")
	if err := writeLines(outLines, communicationsLua); err != nil {
		fatalErrorExitPrompt("Making Changes: %s", err)
	}
}
