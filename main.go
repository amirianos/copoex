package main

import (
        "bytes"
        "fmt"
        "io/ioutil"
        "log"
        "net/http"
        "os/exec"
        "strconv"
        "strings"
        "time"

        "github.com/prometheus/client_golang/prometheus"
        "github.com/prometheus/client_golang/prometheus/promhttp"
        "gopkg.in/ini.v1"
        "gopkg.in/yaml.v2"
)

var (
        // Define a Prometheus gauge metric
        commandOutputGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
                Name: "command_executer_output",
                Help: "Output of the executed command",
        }, []string{"command", "status"})
)

type Command struct {
        Name    string `yaml:"name"`
        Command string `yaml:"command"`
}

func init() {
        // Register the gauge metric with Prometheus's default registry
        prometheus.MustRegister(commandOutputGauge)
}

func main() {
        // Load the INI configuration file
        cfg, err := ini.Load("config.ini")
        if err != nil {
                log.Fatalf("Failed to read config file: %v", err)
        }

        // Read the port number from the INI file
        port := cfg.Section("server").Key("port").String()
        if port == "" {
                log.Fatal("Port number not found in config file")
        }

        // Read the YAML file
        data, err := ioutil.ReadFile("commands.yaml")
        if err != nil {
                log.Fatalf("Error reading YAML file: %v", err)
        }

        // Parse the YAML file into a slice of Command structs
        var commands []Command
        if err := yaml.Unmarshal(data, &commands); err != nil {
                log.Fatalf("Error unmarshaling YAML data: %v", err)
        }

        // Set up a HTTP server to expose the metrics
        http.Handle("/metrics", promhttp.Handler())

        // Start a goroutine to periodically update the metric
        go func() {
                for {
                        for _, cmd := range commands {
                                output, err := executeCommand(cmd.Command)
                                if err != nil {
                                        log.Printf("Error executing command %s: %v", cmd.Name, err)
                                        commandOutputGauge.WithLabelValues(cmd.Name, "error").Set(0)
                                } else {
                                        commandOutputGauge.WithLabelValues(cmd.Name, "success").Set(output)
                                }
                        }
                        // Update the metric every 10 seconds
                        time.Sleep(10 * time.Second)
                }
        }()

        // Start the HTTP server
        log.Printf("HTTP server listening on port %s", port)
        log.Fatal(http.ListenAndServe(":"+port, nil))
}

func executeCommand(command string) (float64, error) {
        cmd := exec.Command("sh", "-c", command)
        var out bytes.Buffer
        cmd.Stdout = &out

        if err := cmd.Run(); err != nil {
                return 0, fmt.Errorf("error running command: %v", err)
        }

        // Convert the output to a float64
        countStr := strings.TrimSpace(out.String())
        count, err := strconv.ParseFloat(countStr, 64)
        if err != nil {
                return 0, fmt.Errorf("error converting output to float64: %v", err)
        }

        return count, nil
}

