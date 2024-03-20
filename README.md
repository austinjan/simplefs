# Simple File Server with TFTP and HTTP Support

This Go application serves as a simple file server that provides both TFTP (Trivial File Transfer Protocol) and HTTP interfaces for uploading, downloading, and listing files. It allows users to interact with a specified directory through a web interface or TFTP client, making it versatile for various file management and transfer scenarios.

## Prerequisites

- Go 1.15 or higher

## Installation

1. **Clone the Repository**

```
$ git clone https://github.com/austinjan/simplefs
$ cd simplefs
```


2. **Build the Application**

Navigate to the application directory and build the application using the Go command:

```
$ go build .
```

This will generate an executable file in the current directory.

## Usage

To start the server, simply run the executable with optional flags for port and directory configuration:

```
./simplefs -p <port> -d <directory>
```

- `-p` specifies the port to listen on for HTTP requests (default is 8080).
- `-d` specifies the directory from which files will be served (default is the current directory).

### Endpoints

The server provides the following HTTP endpoints:

- `/list` - Lists all files in the specified directory.
- `/upload` - Allows file upload through a POST request with a multipart form.
- `/files/{file_name}` - Accesses a specific file by its name.

In addition to the HTTP endpoints, the server listens for TFTP requests on UDP port 69. It supports both read and write operations through TFTP.

## Contributing

Contributions to improve the file server are welcome. Before contributing, please ensure you're familiar with Go programming and the basics of TFTP and HTTP protocols. If you have an idea for an enhancement or have found a bug, feel free to open an issue or submit a pull request.

## License

Specify your license here or state that the project is licensed under the MIT License, GPL, etc.
