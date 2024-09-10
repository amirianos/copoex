package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"net/http"
	"time"
	"io/ioutil"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/yaml.v2"
)

const (
	metricsPort    = ":8099"
	updateInterval = 10 * time.Second
)

// Command defines the structure of each command from the YAML file
type Command struct {
	Name    string `yaml:"name"`
	Command string `yaml:"command"`
}

var (
	// Define a Prometheus gauge metric
	commandGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "command_output",
		Help: "Output of executed commands",
	}, []string{"command_name"})
)

func init() {
	// Register the gauge metric with Prometheus's default registry
	prometheus.MustRegister(commandGauge)
}

func main() {
	commands, err := loadCommands("commands.yaml")
	if err != nil {
		log.Fatalf("Failed to load commands: %v", err)
	}

	// Set up HTTP server for metrics
	http.Handle("/metrics", promhttp.Handler())

	// Start a goroutine to periodically update the metrics
	go startMetricsUpdater(commands, commandGauge)

	// Start the HTTP server
	log.Printf("Starting HTTP server on %s", metricsPort)
	if err := http.ListenAndServe(metricsPort, nil); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}

// loadCommands reads and parses the YAML file into a slice of Command structs
func loadCommands(filePath string) ([]Command, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading YAML file: %w", err)
	}

	var commands []Command
	if err := yaml.Unmarshal(data, &commands); err != nil {
		return nil, fmt.Errorf("error unmarshaling YAML data: %w", err)
	}

	return commands, nil
}

// startMetricsUpdater periodically runs the commands and updates the Prometheus metrics
func startMetricsUpdater(commands []Command, gauge *prometheus.GaugeVec) {
	ticker := time.NewTicker(updateInterval)
	defer ticker.Stop()

	for {
		for _, cmd := range commands {
			output := runCommand(cmd.Command)
			gauge.WithLabelValues(cmd.Name).Set(output)
		}
		<-ticker.C // Wait for the next tick
	}
}

// runCommand executes a shell command and returns the output as a float64
func runCommand(command string) float64 {
	cmd := exec.Command("sh", "-c", command)
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		log.Printf("Error running command '%s': %v", command, err)
		return 0
	}

	outputStr := strings.TrimSpace(out.String())
	output, err := strconv.ParseFloat(outputStr, 64)
	if err != nil {
		log.Printf("Error converting output '%s' to float: %v", outputStr, err)
		return 0
	}

	return output
}
