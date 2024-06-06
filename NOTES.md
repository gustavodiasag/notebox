# Notes

## Supply Chain Attacks

> Modern software engineering is collaborative, and based on reusing Open Source software. That exposes targets to supply chain attacks, where software projects are attacked by compromising their dependencies.

## Network Addresses

The TCP network address accepted by `http.ListenAndServe()` must be in the format `"host:port"`. When instead defined as a a number the port is written as a named one like `"http"`, Go will check for the relevant port number from `/etc/services` when starting the server.

## Fixed and Subtree Paths

In Go's servemux, fixed path patterns such as `/foo/bar` are only matched and handled by the corresponding handler when the request URL path exactly matches the fixed path. In contrast, a subtree path like `/foobar/` are matched and handled whenever the **start** of a request URL path matches the subtree path.

## DefaultServeMux

When a ServeMux instance is not defined for registering routes, a `net/http` global variable called `DefaultServeMux` is used instead, which is not recommended for production applications, any package is able to access it and register a route. If one of those third-party packages is compromised, they could use `DefaultServerMux` to expose a malicious handler tot the web.

## Go Interfaces

A few reasons for using interfaces in go involve its help for reducing duplication or boilerplace code, its ability to make it easier to use mocks instead of real objects in unit tests and as an architectural tool, to help enforce decoupling between parts of a codebase.

## Go File Server

- Sanitizes all request paths by running them through the `path.Clean()` function before searching for a file, helping to stop **directory traversal attacks**. It is particularly useful in the context of a fileserver in conjunction with a router that doesn't automatically sanitize URL paths.

- Provides support for range requests, which can be good when serving large files as a way to provide resumable downloads, e.g. `curl -i -H "Range: bytes=100-199" ...`.

## Concurrent Logging

> Custom loggers created by `log.New()` are concurrency-safe. You can share a single logger and use it across multiple goroutines and in your handlers without needing to worry about race conditions.

## Dependency Availability

- **Global variables**: defining dependencies as global variables is only a possible fit if the application is small and simple, so keeping track of globals isn't a problem, HTTP handlers are spread across multiple packages but the dependency-related code lives in one package are there isn't a need develop a mock for testing purposes.

- **Dependency injection**: has the advantage of clearly defining what dependencies the handlers have and what values they take at runtime. Another benefit is that any unit tests for handlers can be completely self-contained, meaning that there's no need to rely on any global variables **outside of the test**. In general, dependency injection is a useful approach when:
    - There is a common set of dependencies that handlers need access to.
    - All HTTP handlers live in one package, whereas dependency-related code don't.

# Go's sql.Result

The type returned by `DB.Exec()` provides two methods, `LastInsertId()`, which returns the integer returned by the database in response to a command, and `RowsAffected()`, which returns the number of rows affected by the statement.

> Not all drivers and databases support the methods mentioned. PostgreSQL for example does not support the `LastInsertId()` command (check docs).

# Prepared statement

In database management systems, a prepared statement is a feature where ethe database pre-compiles SQL code and stores the results, separating it from data.

Benefits:

- **Efficiency**: they can be easily used repeatedly without re-compiling.

- **Security**: by reducing SQL injection attacks.

# Useful References

- [semver](https://semver.org/)

- [embed](https://pkg.go.dev/embed)

- [12factor](https://12factor.net/)
