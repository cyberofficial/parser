package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/zveinn/parser"
)

type Drive struct {
	Name string
	Size int64
	Type string
}

type Server struct {
	Name     string
	Memory   int64
	CPU      int
	Drives   []Drive
	IsActive bool
}

func main() {
	servers := []Server{
		{
			Name:   "web-server-1",
			Memory: 8589934592, // 8GB
			CPU:    4,
			Drives: []Drive{
				{Name: "OS", Size: 536870912000, Type: "SSD"},    // 500GB
				{Name: "Data", Size: 2199023255552, Type: "HDD"}, // 2TB
			},
			IsActive: true,
		},
		{
			Name:   "db-server-1",
			Memory: 34359738368, // 32GB
			CPU:    8,
			Drives: []Drive{
				{Name: "OS", Size: 1073741824000, Type: "SSD"},   // 1TB
				{Name: "Data", Size: 5497558138880, Type: "SSD"}, // 5TB
			},
			IsActive: true,
		},
		{
			Name:   "backup-server",
			Memory: 4294967296, // 4GB
			CPU:    2,
			Drives: []Drive{
				{Name: "OS", Size: 268435456000, Type: "HDD"},       // 250GB
				{Name: "Backup", Size: 10995116277760, Type: "HDD"}, // 10TB
			},
			IsActive: false,
		},
	}

	examples := []string{
		"Memory > 16GB",                                // Find servers with more than 16GB RAM
		"CPU >= 4 AND IsActive = true",                 // Find active servers with 4+ CPUs
		"Drives.Size > 1TB",                            // Find servers with drives larger than 1TB
		"Name CONTAINS 'web'",                          // Find web servers
		"Memory < 8GiB OR CPU < 4",                     // Find under-powered servers
		"Drives.Type = 'SSD' AND Drives.Size > 500GiB", // Find servers with large SSDs
	}

	fmt.Println("Server Query Examples with Humanized Values")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println()

	for i, query := range examples {
		fmt.Printf("Example %d: %s\n", i+1, query)

		results, err := parser.Parse(query, servers)
		if err != nil {
			log.Printf("Error parsing query '%s': %v", query, err)
			continue
		}

		fmt.Printf("Results (%d found):\n", len(results))
		for _, server := range results {
			fmt.Printf("  - %s (Memory: %s, CPU: %d, Active: %v)\n",
				server.Name,
				formatBytes(server.Memory),
				server.CPU,
				server.IsActive)
			for _, drive := range server.Drives {
				fmt.Printf("    * %s: %s (%s)\n", drive.Name, formatBytes(drive.Size), drive.Type)
			}
		}
		fmt.Println()
	}
}

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
