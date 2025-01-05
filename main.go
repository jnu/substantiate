// Tool to copy secrets from env variables to a file.
//
// Writes environment variables to files in a directory.
//
// The files will have the same name as the environment variable.
//
// We will attempt to base64-decode the value of the environment variable before writing it to the file.
// If the value is not base64-decodable, we will write the value as is.
//
// Either the `vars` or the environment variable `SUBSTANTIATE` can be used to specify the list of
// environment variables to write to files. If none are given, the script exits without doing anything.
//
// Flags:
//  -directory string (default "/secrets") - Directory where the secrets files will be created.
//  -vars string (default "") - Comma separated list of environment variables to copy to files.

package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Decode base64 content if possible, otherwise return as is.
func maybeDecodeContent(content string) string {
	decoded, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return content
	}
	return string(decoded)
}

// Write the content of an environment variable to a file in the given directory.
//
// Decodes content if possible, creates the directory if necessary.
func writeToFile(directory string, variable string, value string) error {
	// Create the directory if it does not exist.
	_, err := os.Stat(directory)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(directory, 0755)
			if err != nil {
				log.Println("Directory does not exist and could not create it ", err)
				return err
			}
		} else {
			log.Println("Error stat-ing directory ", err)
			return err
		}
	}

	// Create the file.
	fp := filepath.Join(directory, variable)
	file, err := os.Create(fp)
	if err != nil {
		log.Println("Error creating file ", err)
		return err
	}

	content := maybeDecodeContent(value)
	_, err = file.WriteString(content)
	if err != nil {
		log.Println("Error writing to file", err)
		return err
	}

	return nil
}

func main() {
	directory := flag.String("directory", "/secrets", "Directory where the secrets files will be created.")
	vars := flag.String("vars", "", "Comma separated list of environment variables to copy to files.")
	flag.Parse()

	// `variables` is a list of all environment variables.
	raw_vars := os.Getenv("SUBSTANTIATE")
	if *vars != "" {
		raw_vars = *vars
	}

	if raw_vars == "" {
		fmt.Println("No environment variables specified.")
		os.Exit(0)
	}

	variables := strings.Split(raw_vars, ",")

	ctr := 0
	for _, variable := range variables {
		value, exists := os.LookupEnv(variable)
		if exists {
			err := writeToFile(*directory, variable, value)
			if err != nil {
				fmt.Printf("Error writing to file: %s\n", err)
			} else {
				ctr++
			}
		} else {
			fmt.Printf("Environment variable %s does not exist.\n", variable)
		}
	}

	log.Println("Wrote", ctr, "file(s) to", *directory)
	os.Exit(0)
}
