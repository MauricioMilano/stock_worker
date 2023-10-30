# RabbitMQ Worker Example

This Go application demonstrates a worker that interacts with RabbitMQ and consumes data from the Stooq API.

## Prerequisites

- Go 1.21.3 or later
- Docker (for running RabbitMQ and other dependencies)

## Usage

1. Clone this repository: 

2. Build the Go application:

    ```bash
    go build main.go
    ```

3. Run the application:

    ```bash
    ./main
    ```

4. The worker will:
   - Read messages from one RabbitMQ queue.
   - Process the data (e.g., retrieve stock information from Stooq).
   - Post the processed data to another RabbitMQ queue.

## Configuration

- Update the RabbitMQ connection details in `.env`.

