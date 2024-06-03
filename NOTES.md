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

- Provides support for range requests, which can be good when serving large files as a way to provide resumable downloads, e.g. `curl -i -H "Range: bytes=100-199" ...`

# Useful References

- [semver](https://semver.org/)

- [embed](https://pkg.go.dev/embed)
