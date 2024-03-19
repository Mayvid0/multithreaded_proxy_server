# Multithreaded Proxy Server with LRU Cache , Concurrent Access and Efficient Logging

This project is a multithreaded proxy server implemented in Go, designed to handle incoming HTTP requests efficiently. It includes features such as caching responses using an LRU (Least Recently Used) cache and concurrent access to manage multiple requests concurrently, along with access logs.

## Features

- **Proxy Server**: The server accepts incoming HTTP requests and forwards them to the target server. It utilizes goroutines and channels to handle concurrent requests effectively.
  
- **LRU Cache**: The server implements an LRU cache to store responses from the target server. This cache helps in serving subsequent requests for the same resource quickly by retrieving cached responses instead of fetching them again from the target server.

- **Concurrent Logging**: The project incorporates concurrent logging to log access and error messages efficiently. Goroutines and channels are used to input log messages concurrently into log files, ensuring that multiple log writes and reads can occur simultaneously without conflicts.

- **Thread Safety**: To maintain thread safety and prevent data races, the project employs mutexes for cache operations. Mutexes ensure that only one goroutine accesses the cache at a time, preventing concurrent writes or reads that could lead to inconsistencies or race conditions.

## Contributing

Contributions to this project are welcome! If you have suggestions, bug reports, or feature requests, please open an issue or submit a pull request to the repository.

## License

This project is licensed under the MIT License, allowing for flexibility in usage and modification while retaining copyright and liability limitations.
