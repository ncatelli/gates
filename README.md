# gates
A gate service submission for the LLJam 0001

## Dependencies
- make
- go 1.17
- docker

## Building
### Docker
The tool can be built and run entirely via docker using the following command.

```sh
$> docker build -t ncatelli/gates .
```

### Locally
The tool can also be built locally via make 

```
$> make build
```

## Testing
Tests can be run using the built in go testing library and a convenient wrapper to test all subpackages has been provided below.

### Locally
Local tests default to running tests on all subpackages along with coverage tests.
Tests can be run with the following make command.

```
$> make test
```

## Configuration
### Services
The gate service can be configured via the following environment variables:

- LISTEN_ADDR:        `string`  The server address gates binds to.
- GATE_TY: `string`  A gate to emulate. [choices: and, not, or, xor]
- OUTPUT_TY:  `url.URL` A an outputer to send a computed output to. [choices: stdout, http]
- OUTPUT_ADDRS:  `[]url.URL` An optional list of address for the http outputter to send a computed output to.
