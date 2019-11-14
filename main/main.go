package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/alexodle/wrappy"
	"os"
	"strings"
)

func main() {
	whitelistFile := flag.String("whitelist", "", "file containing newline-delimited list of full struct names (i.e. github.com/alexodle/wrappy.Struct)")

	flag.Parse()
	args := flag.Args()
	if len(args) != 2 {
		fmt.Printf("go run main.go [--whitelist <whitelist_file>] <input_dir> <output_dir>\n")
		os.Exit(1)
	}

	inputDir, outputDir := args[0], args[1]
	whitelist := parseWhitelistFile(whitelistFile)

	fmt.Printf("whitelist:\n")
	fmt.Println(whitelist)
	wrappy.GenerateWrappers(inputDir, outputDir, wrappy.GenerateWrappersOptions{
		StructWhitelist: whitelist,
	})
}

func parseWhitelistFile(whitelistFile *string) map[string]struct{} {
	if whitelistFile == nil || *whitelistFile == "" {
		return nil
	}

	whitelist := map[string]struct{}{}

	file, err := os.Open(*whitelistFile)
	if err != nil {
		panic(err)
	}
	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		structName := strings.TrimSpace(scanner.Text())
		if structName != "" {
			whitelist[structName] = struct{}{}
		}
	}

	if err := scanner.Err(); err != nil {
		panic(fmt.Errorf("failed to open whitelist file: %s", *whitelistFile))
	}

	return whitelist
}
