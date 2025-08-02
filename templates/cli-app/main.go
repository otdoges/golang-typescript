package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"PROJECT_NAME/async"
	"PROJECT_NAME/types"
	"PROJECT_NAME/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	verbose bool
	config  = types.NewMap[string, interface{}]()
)

func main() {
	fmt.Println("🚀 TypeScript-like Go CLI Application")
	
	execute()
}

func execute() {
	rootCmd := &cobra.Command{
		Use:   "cli-app",
		Short: "A CLI application built with TypeScript-like Go patterns",
		Long: `A command-line application that demonstrates TypeScript-like patterns in Go.
		
This CLI shows how to build modern command-line tools using:
- Optional types for safe null handling
- Promise-based async operations
- Event-driven architecture
- Rich error handling
- Collection utilities`,
	}

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cli-app.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	// Initialize configuration
	cobra.OnInitialize(initConfig)

	// Add commands
	rootCmd.AddCommand(processCmd())
	rootCmd.AddCommand(interactiveCmd())
	rootCmd.AddCommand(configCmd())
	rootCmd.AddCommand(exampleCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
		os.Exit(1)
	}
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".cli-app")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil && verbose {
		fmt.Println("📄 Using config file:", viper.ConfigFileUsed())
	}
}

func processCmd() *cobra.Command {
	var inputFile, outputFile string
	var workers int

	cmd := &cobra.Command{
		Use:   "process",
		Short: "Process files with TypeScript-like patterns",
		Long:  "Demonstrates file processing using async operations and TypeScript-like utilities",
		RunE: func(cmd *cobra.Command, args []string) error {
			return processFiles(inputFile, outputFile, workers)
		},
	}

	cmd.Flags().StringVarP(&inputFile, "input", "i", "", "Input file to process")
	cmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file")
	cmd.Flags().IntVarP(&workers, "workers", "w", 4, "Number of worker goroutines")
	
	cmd.MarkFlagRequired("input")

	return cmd
}

func interactiveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "interactive",
		Short: "Interactive mode with TypeScript-like patterns",
		Long:  "Run in interactive mode to explore TypeScript-like Go features",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInteractiveMode()
		},
	}

	return cmd
}

func configCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Configuration management",
		Long:  "Manage application configuration using TypeScript-like patterns",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "show",
		Short: "Show current configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			return showConfig()
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "set <key> <value>",
		Short: "Set configuration value",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return setConfig(args[0], args[1])
		},
	})

	return cmd
}

func exampleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "examples",
		Short: "Run TypeScript-like Go examples",
		Long:  "Demonstrate various TypeScript-like patterns implemented in Go",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runExamples()
		},
	}

	return cmd
}

// Command implementations

func processFiles(inputFile, outputFile string, workers int) error {
	fmt.Printf("📁 Processing file: %s\n", inputFile)
	
	// Check if input file exists using Optional pattern
	fileInfo := checkFileExists(inputFile)
	if fileInfo.IsNone() {
		return types.NewError("Input file not found", types.ValidationError)
	}

	// Use Promise for async file processing
	promise := async.NewPromise(func() (string, error) {
		// Simulate file processing
		fmt.Printf("⚙️  Processing with %d workers...\n", workers)
		
		// Read file content
		content, err := os.ReadFile(inputFile)
		if err != nil {
			return "", err
		}

		// Process content using TypeScript-like utilities
		lines := strings.Split(string(content), "\n")
		
		// Filter non-empty lines
		nonEmptyLines := utils.Filter(lines, func(line string) bool {
			return !utils.Strings.IsBlank(line)
		})

		// Transform lines (example: add line numbers)
		numberedLines := utils.MapWithIndex(nonEmptyLines, func(line string, index int) string {
			return fmt.Sprintf("%d: %s", index+1, line)
		})

		// Join back
		result := strings.Join(numberedLines, "\n")

		// Simulate processing time
		time.Sleep(500 * time.Millisecond)

		return result, nil
	})

	// Await result
	result, err := promise.Await()
	if err != nil {
		return types.NewError("Processing failed", types.InternalError).WithCause(err)
	}

	// Write output
	if outputFile != "" {
		err = os.WriteFile(outputFile, []byte(result), 0644)
		if err != nil {
			return types.NewError("Failed to write output", types.InternalError).WithCause(err)
		}
		fmt.Printf("✅ Output written to: %s\n", outputFile)
	} else {
		fmt.Println("📄 Processed content:")
		fmt.Println(result)
	}

	return nil
}

func runInteractiveMode() error {
	fmt.Println("🎮 Interactive Mode - TypeScript-like Go CLI")
	fmt.Println("Type 'help' for available commands, 'exit' to quit")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)
	
	// Event emitter for commands
	commandEvents := types.NewEventEmitter[CommandEvent]()
	
	// Setup command listeners
	commandEvents.On("command:executed", func(event CommandEvent) {
		if verbose {
			fmt.Printf("📝 Command '%s' executed in %v\n", event.Command, event.Duration)
		}
	})

	for {
		fmt.Print("ts-go> ")
		
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		if input == "exit" || input == "quit" {
			fmt.Println("👋 Goodbye!")
			break
		}

		start := time.Now()
		err := handleInteractiveCommand(input)
		duration := time.Since(start)

		commandEvents.Emit("command:executed", CommandEvent{
			Command:  input,
			Duration: duration,
			Error:    err,
		})

		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
		}
		fmt.Println()
	}

	return nil
}

func showConfig() error {
	fmt.Println("📋 Current Configuration:")
	
	if config.Size() == 0 {
		fmt.Println("  No configuration set")
		return nil
	}

	config.ForEach(func(key string, value interface{}) {
		fmt.Printf("  %s: %v\n", key, value)
	})

	return nil
}

func setConfig(key, value string) error {
	config.Set(key, value)
	fmt.Printf("✅ Set %s = %s\n", key, value)
	return nil
}

func runExamples() error {
	fmt.Println("🧪 TypeScript-like Go Examples")
	fmt.Println("=" * 40)

	// Optional Types Example
	fmt.Println("\n📦 Optional Types:")
	name := types.Some("CLI User")
	age := types.None[int]()
	fmt.Printf("  Name: %s\n", name.GetOrDefault("Anonymous"))
	fmt.Printf("  Age: %d\n", age.GetOrDefault(0))

	// Array Utilities Example
	fmt.Println("\n🔧 Array Utilities:")
	numbers := []int{1, 2, 3, 4, 5}
	doubled := utils.Map(numbers, func(x int) int { return x * 2 })
	evens := utils.Filter(numbers, func(x int) bool { return x%2 == 0 })
	fmt.Printf("  Original: %v\n", numbers)
	fmt.Printf("  Doubled: %v\n", doubled)
	fmt.Printf("  Evens: %v\n", evens)

	// Collections Example
	fmt.Println("\n🗂️  Collections:")
	tasks := types.NewSet[string]()
	tasks.Add("Write code").Add("Write tests").Add("Deploy")
	fmt.Printf("  Tasks: %v\n", tasks.Values())
	fmt.Printf("  Total tasks: %d\n", tasks.Size())

	// Error Handling Example
	fmt.Println("\n🚨 Error Handling:")
	result, err := types.NewTry(func() (string, error) {
		return "Success!", nil
	}).Execute()
	
	if err == nil {
		fmt.Printf("  Result: %s\n", result)
	}

	fmt.Println("\n✅ All examples completed!")
	return nil
}

// Helper functions

func checkFileExists(filename string) types.Optional[os.FileInfo] {
	info, err := os.Stat(filename)
	if err != nil {
		return types.None[os.FileInfo]()
	}
	return types.Some(info)
}

func handleInteractiveCommand(input string) error {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return nil
	}

	command := parts[0]
	args := parts[1:]

	switch command {
	case "help":
		fmt.Println("Available commands:")
		fmt.Println("  help           - Show this help")
		fmt.Println("  status         - Show application status")
		fmt.Println("  list <type>    - List items of type")
		fmt.Println("  count <items>  - Count items")
		fmt.Println("  config <key>   - Get config value")
		fmt.Println("  exit/quit      - Exit interactive mode")

	case "status":
		fmt.Printf("✅ Application running\n")
		fmt.Printf("  Config items: %d\n", config.Size())
		fmt.Printf("  Current directory: %s\n", getCurrentDir().GetOrDefault("unknown"))

	case "list":
		if len(args) == 0 {
			return types.NewValidationError("Please specify what to list")
		}
		handleListCommand(args[0])

	case "count":
		if len(args) == 0 {
			return types.NewValidationError("Please specify items to count")
		}
		count := len(args)
		fmt.Printf("Item count: %d\n", count)

	case "config":
		if len(args) == 0 {
			return showConfig()
		}
		key := args[0]
		if value := config.Get(key); value.IsSome() {
			fmt.Printf("%s: %v\n", key, value.Get())
		} else {
			fmt.Printf("Config key '%s' not found\n", key)
		}

	default:
		return types.NewValidationError(fmt.Sprintf("Unknown command: %s", command))
	}

	return nil
}

func handleListCommand(listType string) {
	switch listType {
	case "files":
		files, err := filepath.Glob("*")
		if err == nil {
			fmt.Printf("Files in current directory: %v\n", files)
		}
	case "config":
		showConfig()
	default:
		fmt.Printf("Don't know how to list: %s\n", listType)
	}
}

func getCurrentDir() types.Optional[string] {
	dir, err := os.Getwd()
	if err != nil {
		return types.None[string]()
	}
	return types.Some(dir)
}

// Event types
type CommandEvent struct {
	Command  string
	Duration time.Duration
	Error    error
}