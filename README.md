
# copoex (Command Prometheus Exporter)

**copoex** is a lightweight Prometheus exporter written in Go that runs commands specified in a YAML file and exposes the results as Prometheus metrics. This tool is ideal for monitoring the status and output of custom shell commands or system scripts in a structured and scalable way.
## Features

 - Custom Commands: Define and run commands directly from a YAML configuration file.
 - Prometheus Metrics: Expose the output of your commands in a format that Prometheus can scrape.
 -  Extensible: Add any command that can be executed from the command line.
 - Written in Go: Fast and efficient execution with minimal system overhead.

## Installation

 1. Clone the repository:

    git clone https://github.com/amirianos/copoex.git
    cd copoex

2. Build the project:
   

     CGO_ENABLED=0 GOOS=linux go build -o copoex
    
 Alternatively, you can use the pre-built binary (You can download compiled binary files from Release page).

## Configuration

The behavior of copoex is controlled through a YAML file. Here's an example configuration (config.yaml):

commands.yaml file :

    - name: rootUsedSpaceCommand
      command: df -h / | awk 'NR==2 {print $5}' | sed 's/%//'
    - name: loadAverageCommand
      command: uptime | awk '{print $12}' | cut -d "," -f 1

You can add as many commands as needed.

## Usage

1. After building or downloading the binary, run the exporter by specifying your configuration file:

    ./copoex -commands CONFIG_PATH -port ":PORT_NUMBER"

   Note that the default port number is `8099`, and the default `commands.yaml` file path is located next to the copoex binary file. If you want to change these default values, you can run it with the `-commands` and `-port` switches.

2. By default, the exporter will be available at http://INSTANCE_IP:8099/metrics for Prometheus to scrape.

## Command-Line Options

    -commands: Path to the YAML commands file (default: commands.yaml).
    -port: Port for exposing the Prometheus metrics (default: 8099).

## Prometheus Integration

Add the following configuration to your Prometheus scrape_configs to scrape metrics from copoex:


    scrape_configs:
      - job_name: 'copoex'
        static_configs:
          - targets: ['localhost:8099']

Replace localhost:8099 with the actual address and port where copoex is running.

## Example Metrics

The commands defined in the configuration file will be executed periodically, and their output will be exposed as metrics. For example:

     HELP copoex_command_duration_seconds Duration of command execution in seconds.
     TYPE copoex_command_duration_seconds gauge
    copoex_command_duration_seconds{name="loadAverageCommand"} 2

## Contributing

Feel free to contribute to this project by submitting issues, feature requests, or pull requests.


