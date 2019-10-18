package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type PSI struct {
	Timestamp string    `json:"timestamp"`
	Kind      string    `json:"kind"`
	Avg10     string    `json:"avg10"`
	Avg60     string    `json:"avg60"`
	Avg300    string    `json:"avg300"`
	Total     string    `json:"total"`
}

func handleError(err error) {
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func parseFlags() (path, format string) {
	flag.StringVar(&path, "path", "", "the path to the export file.")
	flag.StringVar(&format, "format", "json", "the output format for the psi information. Available options are 'csv' and 'json. The default format is json.")

	flag.Parse()

	format = strings.ToLower(format)

	if format != "csv" && format != "json" {
		fmt.Println("Error: invalid format. Use 'json' or 'csv' instead.")
		flag.Usage()
		os.Exit(1)
	}

	return
}

func collectPSIs(samples int) (psis []PSI) {
	BYTES := 128
	// reader := csv.NewReader(f)
	// reader.Comma = ' '

	// handleError(err)

	// some avg10=0.00 avg60=0.00 avg300=0.00 total=753129657

	for i := 1; i < samples; i++ {
		c, err := os.Open("/proc/pressure/cpu")
		i, err := os.Open("/proc/pressure/io")
		m, err := os.Open("/proc/pressure/memory")
		handleError(err)
	//	reader := csv.NewReader(f)
	//	line, err := reader.Read()
		timestamp := time.Now().String()
		cpu := make([]byte, BYTES)
		c.Read(cpu)
		fmt.Println("cpu", string(cpu), timestamp)
		io := make([]byte, BYTES)
		i.Read(io)
		fmt.Println("io", string(io), timestamp)
		mem := make([]byte, BYTES)
		m.Read(mem)
		fmt.Println("mem", string(mem), timestamp)
	//	timestamp := time.Now()
	//	psi := PSI{
	//		Timestamp: timestamp.String(),
	//		Kind:      string(line[0]),
	//		Avg10:     string(line[1]),
	//		Avg60:     string(line[2]),
	//		Avg300:    string(line[3]),
	//		Total:     string(line[4]),
	//	}
	//	psis = append(psis, psi)
		time.Sleep(1 * time.Second)
		c.Close()
		i.Close()
		m.Close()
	}

	return
}

func main() {
	var output io.Writer
	//var samples int
	samples := 10000
	path, format := parseFlags()
	psis := collectPSIs(samples)

	if path != "" {
		f, err := os.Create(path)
		handleError(err)
		defer f.Close()
		output = f
	} else {
		output = os.Stdout
	}

	if format == "json" {
		data, err := json.MarshalIndent(psis, "", "  ")
		handleError(err)
		output.Write(data)
	} else if format == "csv" {
		output.Write([]byte("timestamp,kind,avg10,avg60,avg300,total\n"))
		writer := csv.NewWriter(output)
		for _, psi := range psis {
			err := writer.Write([]string{psi.Timestamp, psi.Kind, psi.Avg10, psi.Avg60, psi.Avg300, psi.Total})
			handleError(err)
		}
		writer.Flush()
	}
}

