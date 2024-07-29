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

## Go's sql.Result

The type returned by `DB.Exec()` provides two methods, `LastInsertId()`, which returns the integer returned by the database in response to a command, and `RowsAffected()`, which returns the number of rows affected by the statement.

> Not all drivers and databases support the methods mentioned. PostgreSQL for example does not support the `LastInsertId()` command (check docs).

## Prepared Statement

In database management systems, a prepared statement is a feature where the database pre-compiles SQL code and stores the results, separating it from data.

Benefits:

- **Efficiency**: they can be easily used repeatedly without re-compiling.

- **Security**: by reducing SQL injection attacks.

## Closing a Resultset

After an execution of a SQL statement, it is crucial to close a resultset because as long as it is open, it will keep the underlying database connection open, so if something goes wrong in this method and the resultset isn't closed, it can rapidly lead to all the connections in the pool being used.

## Reusable Prepared Statements

Example in Go:

```go
// We need somewhere to store the prepared statement for the lifetime of our
// web application. A neat way is to embed in the model alongside the connection
// pool.
type ExampleModel struct {
    DB         *sql.DB
    InsertStmt *sql.Stmt
}

// Create a constructor for the model, in which we set up the prepared
// statement.
func NewExampleModel(db *sql.DB) (*ExampleModel, error) {
    // Use the Prepare method to create a new prepared statement for the
    // current connection pool. This returns a sql.Stmt object which represents
    // the prepared statement.
    insertStmt, err := db.Prepare("INSERT INTO ...")
    if err != nil {
        return nil, err
    }

    // Store it in our ExampleModel object, alongside the connection pool.
    return &ExampleModel{db, insertStmt}, nil
}

// Any methods implemented against the ExampleModel object will have access to
// the prepared statement.
func (m *ExampleModel) Insert(args...) error {
    // Notice how we call Exec directly against the prepared statement, rather
    // than against the connection pool? Prepared statements also support the
    // Query and QueryRow methods.
    _, err := m.InsertStmt.Exec(args...)

    return err
}

// In the web application's main function we will need to initialize a new
// ExampleModel struct using the constructor function.
func main() {
    db, err := sql.Open(...)
    if err != nil {
        errorLog.Fatal(err)
    }
    defer db.Close()

    // Create a new ExampleModel object, which includes the prepared statement.
    exampleModel, err := NewExampleModel(db)
    if err != nil {
       errorLog.Fatal(err)
    }

    // Defer a call to Close() on the prepared statement to ensure that it is
    // properly closed before our main function terminates.
    defer exampleModel.InsertStmt.Close()
}
```

# Middleware

A common way of organizing shared functionalities that are used for many, or even all, HTTP requests in a web application is to set them up as **middleware**. In essence, it is defined by some self-contained code which independently acts on a request before or after the general application handlers.

# Same-Origin Policy

Mechanism that restricts how a document or script loaded by one origin[^1] can interact with a resource from another origin.

[^1]: Defined by the *scheme*, *hostname* and *port* of the URL to access it.  

# Map information on HTML templates

For maps, it's possible to access the value for a given key by simply chaining the key name. So, for example, to render a validation error for a `title` field, the tag `{{ .Form.FieldErrors.title }}` can be used in a template.

# Form Decoder

When `app.formDecoder.Decode()` is called, it requires a **non-nil** pointer as the target decode destination. If something is passed that isn't a non-nil pointer, `Decode()` will return a `form.InvalidDecoderError`.

This error must be checked specifically and managed as a special case, rather than just returning a `400 Bad Request` response.

# Restricting cipher suites

For some applications, it may be desirable to limit a HTTPS server to only support some of the cipher suites that Go suppoerts. For example, one might want to only suppoer cipher suites which user forward secrecy and not support weal cipher suites that use RC4, 3DES or CBC.

One thing to notice, however, is that restricting the supported cipher suites to only include modern ciphers can mean that users with certain older browsers won't be able to use the website.

# Observable Response Discrepancy

The product provides different responses to incomnig requests in a way that reveals internal state information to an unauthorized actor outise of the intended control scope. This issue frequently occurs during authentication, where a difference in failed-login messages could allow an attacker to determine if the username is valid or not. These exposures can be inadvertent or intentional.

# Recording Responses

To assist in testing HTTP handlers, Go provides the `net/http/httptest` package, which contains tools such as the `httptest.ResponseRecorder` type. This is essentially an implementation of `http.ResponseWriter` which records the response status code, headers and body instead of actually writing them to an HTTP connection. 

So an easy way to unit test handlers is to create a new `httptest.ResponseRecorder` object, pass it to the handler function, and then examine it again after the handler returns.

# Count Flag for Tests

The `count` flag is used to tell `go test` how many times you want to execute each test. It's a non-cacheable flag, which means that any time it is used `go test` will neither read or write the test rsults to the cache. so using `count=1` is a trick to avoid the cache without other wise affecting how tests run.

Alternatively, cached results can be cleaned for all tests with the `go clean` command:

```sh
$ go clean -testcache
```

# Ignored files and folders

The Go tool ignores any directories called testdata, so any content inside them will be ignored when compiling a Go application. It also ignores any directories or files which have names that begin with an `_` or `.` character too.

# Useful References

- [semver](https://semver.org/)

- [embed](https://pkg.go.dev/embed)

- [12factor](https://12factor.net/)

- [CSP](https://developer.mozilla.org/en-US/docs/Web/HTTP/CSP)

- [OWASP](https://cheatsheetseries.owasp.org/index.html)

- [CWE](https://cwe.mitre.org/index.html)
