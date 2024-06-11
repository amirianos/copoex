package main

import (
	"bytes"
	"log"
	"time"
	"os/exec"
	"strconv"
	"strings"
	"net/http"
	"io/ioutil"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/yaml.v2"
)

var (
	// Define a Prometheus gauge metric
	fileCountGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "command_executer",
		Help: "output_command_executer",
	}, []string{"command"})
)
type Command struct {
	Name    string `yaml:"name"`
	Command string `yaml:"command"`
}
func init() {
	// Register the gauge metric with Prometheus's default registry
	prometheus.MustRegister(fileCountGauge)
}

func main() {
	data, err := ioutil.ReadFile("commands.yaml")
	if err != nil {
		log.Fatalf("Error reading YAML file: %v", err)
	}

	// Parse the YAML file into a slice of Command structs
	var commands []Command
	err = yaml.Unmarshal(data, &commands)
	if err != nil {
		log.Fatalf("Error unmarshaling YAML data: %v", err)
	}
	// Set up a HTTP server to expose the metrics
	http.Handle("/metrics", promhttp.Handler())

	// Start a goroutine to periodically update the metric
	go func() {
		for {
                        for _, cmd := range commands {
			       output := getFileCount(cmd.Command)
			       fileCountGauge.WithLabelValues(cmd.Name).Set(float64(output))
			       // Update the metric every 10 seconds
			       time.Sleep(10 * time.Second)
		       }
		}
	}()

	// Start the HTTP server
	log.Fatal(http.ListenAndServe(":8099", nil))
}

// getFileCount runs the 'ls | wc -l' command and returns the file count
func getFileCount(command string) float64 {
	cmd := exec.Command("sh", "-c", command)
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		log.Printf("Error running command: %v", err)
		return 0
	}

	// Convert the output to an integer
	countStr := strings.TrimSpace(out.String())
	count, err := strconv.ParseFloat(countStr, 64)
	if err != nil {
		log.Printf("Error converting output to integer: %v", err)
		return 0
	}

	return count
}

