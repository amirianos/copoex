# copoex (Command Prometheus Exporter)

`copoex` is a lightweight Prometheus exporter written in Go that runs commands specified in a YAML file and exposes the results as Prometheus metrics. This tool is ideal for monitoring the status and output of custom shell commands or system scripts in a structured and scalable way.

## Features

-   **Custom Commands**: Define and run commands directly from a YAML configuration file.
-   **Prometheus Metrics**: Expose the output of your commands in a format that Prometheus can scrape.
-   **Extensible**: Add any command that can be executed from the command line.
-   **Written in Go**: Fast and efficient execution with minimal system overhead.

## Installation

1.  Clone the repository:

