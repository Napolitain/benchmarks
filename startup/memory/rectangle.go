package main

import (
	"fmt"
	"os"
	"time"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
)

type RectangleData struct {
	A float64 `yaml:"a"`
	B float64 `yaml:"b"`
	C float64 `yaml:"c"`
	D float64 `yaml:"d"`
}

func computeRectangleArea(data RectangleData) float64 {
	width := abs(data.C - data.A)
	height := abs(data.D - data.B)
	return width * height
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	var rootCmd = &cobra.Command{
		Use:   "rectangle [yaml-file]",
		Short: "Calculate rectangle area from YAML file",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			yamlFile := args[0]

			start := time.Now()

			fileContents, err := os.ReadFile(yamlFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
				os.Exit(1)
			}

			var data RectangleData
			err = yaml.Unmarshal(fileContents, &data)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error parsing YAML: %v\n", err)
				os.Exit(1)
			}

			area := computeRectangleArea(data)

			elapsed := time.Since(start)

			fmt.Printf("Rectangle area: %.2f\n", area)
			fmt.Printf("Time: %.6f ms\n", float64(elapsed.Nanoseconds())/1e6)
		},
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
