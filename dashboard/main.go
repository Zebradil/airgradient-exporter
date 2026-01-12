package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	dash, err := buildDashboard()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error building dashboard: %v\n", err)
		os.Exit(1)
	}

	jsonData, err := json.MarshalIndent(dash, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling dashboard to JSON: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(string(jsonData))
}
