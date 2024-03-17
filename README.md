# Multithreaded Proxy Server with LRU Cache and Concurrent Logging

This project is a multithreaded proxy server implemented in Go. It includes features such as caching responses and concurrent logging to handle multiple simultaneous requests efficiently. Safety measures such as channels, goroutines, and mutexes are used to ensure thread safety and prevent race conditions.

## Features

- **Proxy Server**: Handles incoming HTTP requests and forwards them to the target server.
- **Caching**: Stores responses in a cache to serve subsequent requests for the same resource quickly.
- **Concurrent Logging**: Logs access and error messages concurrently to handle multiple log writes and reads simultaneously.
- **Thread Safety**: Utilizes channels, goroutines, mutexes, and other synchronization mechanisms to maintain thread safety and prevent data races.

## Dependencies

- Go (version 1.16 or higher)
- External packages (specified in the `go.mod` and `go.sum` files)

## Contributing

Contributions are welcome! If you have suggestions, bug reports, or feature requests, please open an issue or submit a pull request to the repository.

## License

This project is licensed under the MIT License. 
