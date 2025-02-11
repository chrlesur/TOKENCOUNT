package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
	tiktoken "github.com/tiktoken-go/tokenizer"
)

func formatNumber(number int) string {
	// Convert the number to a string
	numberString := fmt.Sprintf("%d", number)

	// Get the length of the number string
	stringLength := len(numberString)

	// Create a new string to store the formatted number
	formattedNumberString := ""

	// Iterate over the number string from right to left
	for i := 0; i < stringLength; i++ {
		// Get the current digit
		digit := numberString[stringLength-1-i]

		// Add the digit to the formatted number string
		formattedNumberString = string(digit) + formattedNumberString

		// Add a space after every three digits
		if (i+1)%3 == 0 && (i+1) != stringLength {
			formattedNumberString = " " + formattedNumberString
		}
	}

	// Return the formatted number string
	return formattedNumberString
}

// Global variables
var (
	version   = "1.0.0" // Application version
	debug     bool      // Enable debug mode
	logFile   string    // Specify the log file
	threads   int       // Specify the number of threads to use
	recursive bool      // Explore directories recursively
)

func main() {
	// rootCmd represents the base command when called without any subcommands
	var rootCmd = &cobra.Command{
		Use:     "token-count [files...]",
		Short:   "Counts tokens in files",
		Long:    `tokencount is a simple tool to count tokens in files using a Go native library.`,
		Version: version,
		Args:    cobra.MinimumNArgs(1), // Require at least one file argument
		Run: func(cmd *cobra.Command, args []string) {
			// Initialize logger
			var logOutput io.Writer = os.Stdout // Default log output to stdout
			if logFile != "" {
				// If log file is specified, open it for writing
				file, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
				if err != nil {
					fmt.Printf("Error opening log file: %v\n", err)
					os.Exit(1)
				}
				defer file.Close() // Close the file when the function returns
				logOutput = file   // Set log output to the file
			}
			logger := log.New(logOutput, "", log.LstdFlags) // Create a new logger

			if threads <= 0 {
				threads = runtime.NumCPU() // If threads is not specified, use the number of CPUs
			}

			var files []string         // Create a slice to store the list of files to process
			for _, arg := range args { // Iterate over the arguments
				fileInfo, err := os.Stat(arg) // Get file information
				if err != nil {
					logger.Printf("Error stating file %s: %v\n", arg, err) // Log the error
					continue                                               // Skip to the next argument
				}

				if fileInfo.IsDir() { // If the argument is a directory
					if recursive { // If recursive mode is enabled
						filepath.Walk(arg, func(path string, info os.FileInfo, err error) error { // Walk the directory tree
							if err != nil {
								logger.Printf("Error accessing path %s: %v\n", path, err) // Log the error
								return err                                                // Return the error to stop walking
							}
							if !info.IsDir() { // If the path is not a directory
								files = append(files, path) // Add the path to the list of files
							}
							return nil // Return nil to continue walking
						})
					} else {
						logger.Printf("Skipping directory %s (use --recursive to process)\n", arg) // Log that the directory is being skipped
					}
				} else {
					files = append(files, arg) // If the argument is a file, add it to the list of files
				}
			}

			totalTokens := 0 // Initialize the total token count
			for _, filename := range files {
				content, err := os.ReadFile(filename) // Read the file content
				if err != nil {
					logger.Printf("Error reading file %s: %v\n", filename, err) // Log the error
					continue                                                    // Skip to the next file
				}

				if debug { // If debug mode is enabled
					logger.Printf("File: %s, Content: %s\n", filename, string(content)) // Log the file content
				}

				// Initialize tokenizer
				tokenizer, err := tiktoken.Get("cl100k_base") // Get the tokenizer for the cl100k_base encoding
				if err != nil {
					logger.Printf("Error getting encoding: %v\n", err) // Log the error
					continue                                           // Skip to the next file
				}

				// Tokenize the content
				tokens, _, err := tokenizer.Encode(string(content)) // Encode the file content into tokens
				if err != nil {
					logger.Printf("Error encoding: %v\n", err) // Log the error
					continue                                   // Skip to the next file
				}

				if logFile != "" { // If a log file is specified
					logger.Printf("File: %s (%d tks)\n", filename, len(tokens)) // Log the file name and token count
				} else {
					fmt.Printf("File: %s (%d tks)\n", filename, len(tokens)) // Print the file name and token count to stdout
				}
				totalTokens += len(tokens)
			}
			fmt.Printf("__\n")
			fmt.Printf("Total token count: %s tks\n", formatNumber(totalTokens))
		},
	}

	// Define flags
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug mode")                                             // Enable debug mode
	rootCmd.PersistentFlags().StringVar(&logFile, "log-file", "", "Specify the log file")                                      // Specify the log file
	rootCmd.PersistentFlags().IntVar(&threads, "threads", 0, "Specify the number of threads to use (default: number of CPUs)") // Specify the number of threads to use
	rootCmd.PersistentFlags().BoolVarP(&recursive, "recursive", "r", false, "Explore directories recursively")                 // Explore directories recursively

	if err := rootCmd.Execute(); err != nil { // Execute the root command
		fmt.Println(err) // Print the error to stdout
		os.Exit(1)       // Exit with an error code
	}
}
