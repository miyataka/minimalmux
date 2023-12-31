# minimalmux

`minimalmux` is a simple and efficient HTTP router for Golang. It extends the functionality of the standard library with support for path parameters, middleware integration, and routing restrictions based on HTTP methods. This library is ideal for Golang developers seeking simplicity and extensibility.

## Features

- **Path Parameter Support**: Supports path parameters for flexible URL pattern matching, enhancing dynamic routing capabilities.
- **Middleware Support**: Easily integrates middleware for authentication, logging, request handling, etc.
- **Routing Restriction by HTTP Methods**: Simplifies setting up routing for specific HTTP methods such as GET, POST, PUT, etc.

## Quick Start

TODO: write example.

## License

`minimalmux` is released under the MIT license. For more details, see the [LICENSE](https://github.com/miyataka/minimalmux/blob/main/LICENSE) file.

## Future Work
- [ ] host routing
- [ ] documentation
- [ ] benchmarking
- [ ] default middlewares
    - [ ] logger
    - [ ] debug middlewares
    - [ ] json response headers
    - [ ] compress middlewares
    - [ ] nocache middlewares
