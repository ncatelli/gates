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

- LISTEN_ADDR:		`string`	The server address gates binds to.
- SERVICE_TY:		`string`	A service or gate to emulate. [choices: and, not, or, xor]
- OUTPUT_TY:		`url.URL`	A an outputer to send a computed output to. [choices: stdout, http]
- OUTPUT_ADDRS:		`[]url.URL`	An optional list of address for the http outputter to send a computed output to.

## Usage
### Example
#### Compose
A minimal compose environment that showcases a few gates is available in the repo. This can be started by running:

```bash
$ docker-compose up -d
```

This example exposes a few gates that are wired up to eachother, notably a `not`, `or` gate which feeds their output into an `and` gate that outputs its result. Addtionally a `xor` gate is setup as a standalone gate with it's output wired to stdout.

These be sent signals using the curl template, replacing the values inbetween the `<>` with the the value i will describe below.:

```bash
curl -X POST -sD - -d '{"state": <boolean>, "tick": <unsigned integer>}' <gate host>:8080/input/<input id>
```

##### Post Body
The post body is a json object with two fields

- state: stores a boolean value representing if the input is a `0` or `1`
- tick: an unsigned integer representing a cycle count.

It's worth noting that ticks do not have to be in order. Gates will store inputs for a given tick regardless of the order it is received.

##### URL
Each gate has a unique url, and by default listens on port `8080`.

Each gate exposes it's input via a unique path, where the inputs are ordered by incrementing single lowercase characters starting from `a`.

For example, the url for the solitary input of a `not` gate at the address `not_gate:8080` would be `http://not_gate:8080/input/a`. The urls for the two inputs of an `and` gate at the address `and_gate:8080` would be `http://and_gate:8080/input/a` and `http://and_gate:8080/input/b` etc...

##### Putting it all together
Below shows a curl request to the `not` gate, that sets the input for the first tick to `false`. Causing a `true` output to the `and_gate` service's `a` input.

```bash
$ curl -X POST -sD - -d '{"state": false, "tick": 0}' localhost:8080/input/a
HTTP/1.1 202 Accepted
Date: Sun, 19 Jun 2022 16:02:46 GMT
Content-Length: 0

```

Below shows a subsequent two curl requests to the `or` gate, that set the input for the first tick both to `true`. Causing a `true` output to the `and_gate` service's `b` input.

```bash
$ curl -X POST -sD - -d '{"state": true, "tick": 0}' localhost:8081/input/a
HTTP/1.1 202 Accepted
Date: Sun, 19 Jun 2022 16:01:54 GMT
Content-Length: 0

$ curl -X POST -sD - -d '{"state": true, "tick": 0}' localhost:8081/input/b
HTTP/1.1 202 Accepted
Date: Sun, 19 Jun 2022 16:02:05 GMT
Content-Length: 0
```

When both inputs for tick `0` are received, the `and_gate` service outputs. This service is configured to use the `stdout` outputter via the `OUTPUT_TYPE` env var in the `docker-compose.yaml` file. Thus the results can be seen in the service's log to be set to the expected value of `true`.

```bash
$ docker-compose logs and_gate                                              
Attaching to gates_and_gate_1
and_gate_1  | 2022/06/19 16:04:58 Starting server on 0.0.0.0:8080
and_gate_1  | 2022/06/19 16:04:58 Configured as and gate
and_gate_1  | 2022/06/19 16:05:14 tick: 0, state: true
```

Finally, this can all be cleaned up with the following:

```bash
docker-compose down -v --rmi local
```