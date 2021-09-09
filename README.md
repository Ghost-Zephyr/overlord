# Overlord

Full management suite for libvirt and kubernetes clusters!

## Running it for yourself
Remember to go get dependencies.

### Config
Create a file in the working directory of Overlord named lord.json.
You may also specify a config file with -c command line argument.
This file may look like this.
Everything is optional!
```json
{
  "LibVirtHosts": [
    "qemu:///system"
  ],
  "MongoDbStr": "mongodb://localhost:27017",
  "LogLevel": 0,
  "LogFilePath": "lord.log"
}
```
No log file means stdout only!
The loglevels are;
```go
type LogLevel int
const (
	TRACE LogLevel = iota
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
)
```
Use the corresponding integer starting at 0 for TRACE up to 4 for FATAL.

### Docker
Docker compose based development environment in the workings.
There will also be a Overlord docker image.
