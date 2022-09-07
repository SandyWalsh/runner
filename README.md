# runner design doc

## Background
This repo is intended to hold the client, server (api), and underlying (library) golang code for a gRPC-based runner of linux commands. 

(phew ... that's a mouthful)

Not only that, but it should be secured by TLS, honor cgroups and, optionally, linux namespaces, and stream data back to the client from the running process. 
Impossible you say?! Maybe. Think of it as making docker runtime without the image management.

This document will outline my approach to make this a reality. 

### Broad strokes

I hope to implement this in a series of PR's that follow this schedule:

0. this doc
1. create the library for the process runner and cgroup/ns support (perhaps based on client id roles)
2. basic TLS-secured gRPC client-server chat. The API will wrapper the process runner library
3. support for hard coded gRPC streaming data. Something simple like a clock tick
4. Spawn some basic jobs, capture stderr/stdout and pump them into the streams
5. Work on job control. Cancelling a process, ensuring the channels/goroutines/resources clean up correctly, etc. 
6. Add cgroups and resource isolation 

(this may change)

Stretch goals

7. some key tests
8. docker packaging
9. a nice client
10. magic

## API

Exposed gRPC "Runner" service methods:
```
Run(ClientID, cmd, args...) (Process, error)
Stdout(ClientID, Process) (Stream, error)
Status(ClientID, Process) (Status, error)
Abort(ClientID, Process) (Status, error)
```

where: 
 * the `Run` method that can launch new processes for authenticated callers.
   * the "role" of the authenticated client will determine the cgroup style the process runs in. This will be extracted from client mTLS cert.
   * `Run` takes `cmd` and `args` args for the intended linux command and related cmdline parameters and returns a `Process` message
 * `Process` is an abstracted process ID for the process control calls: `Stdout`, `Status`, and `Abort`. It will be a UUID that maps to the underlying os pid. 
 * `Status` will return an enum of `COMPLETED|RUNNING|ERROR|ABORTED` and, ideally, the exit code

## AuthN / AuthZ

mTLS client certs will be used for authN. The client cert will map to a set of "roles" that determine the class of cgroup processes will run under. 

| cgroup template | CPU | Memory | Disk IO | 
| ------------| --- | ------ | ------ |
| **good**    | x   | x      | x      |
| **better**  | x   | x      | x      |
| **best**    | x   | x      | x      |

| client | assigned cgroup | 
| -------- | --------------- |
| **cert-A** | good            |
| **cert-B** | better        |
| **cert-C** | best            |


## Client UX

The client initially will just be a golang program that runs the server through some paces (just an integration test harness). But in later releases I will try to turn this in to a simple cmdline tool with basic support for each operation. Exact syntax TBD but it may look something like:

```
> client -cert cert.pem run ls -laR
aaaa-bbbb-ccc-ddd-eeee
> client stdout -cert cert.pem -process aaaa-bbbb-ccc-ddd-eeee
.
..
README.md
> client status -cert cert.pem -process aaaa-bbbb-ccc-ddd-eeee
COMPLETED
Exit code: 0
> client abort -cert cert.pem -process aaaa-bbbb-ccc-ddd-eeee
already terminated
```

Beyond the raw client code generated by the gRPC code there will not be a separate "sugar" client library. Although it would be handy for managing the stdout streams. 

## Implementation tradeoffs & security considerations

### Shortcuts

* Long lived gRPC services can be tricky with streaming messages. I will be ignoring all the retry/re-establish connection code. 
* All the authenticated clients will have hardcoded "roles" which map to cgroups. 
* Knowing when the stdout streaming has completed could be tricky. In order to keep `client.Recv` from blocking I may keep sending empty messages server-side. Hacky I know. 
* I may limit cgroup control to just the cpu controller, memory, and disk I/O. 
* "domain" cgroup types will be used. "threaded" won't be supported. 

## Security
The CA and server cert will be generated with openssl and based on a 4096 bit RSA key pair. 
An X.509 cert will be fabricated from these key pairs.
Client certs will also be X.509 as above. 

There are a number of security concerns with such a service:
* in order to exec commands the process will need to run privileged system calls. This requires root level access. Which means a bad actor could basically do anything. We'll need to use cgroups and resource isolation to protect against this. 
* Setting the file permissions on the cgroup v2 files will be important for limiting child processes from manipulating the cgroup settings. 
* the uid/gid for the user to run the processes as will be hardcoded
* The TLS certs will be stored in source control. Ideally these should be in a good secret manager and merged at runtime. 

