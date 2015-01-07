# boom-curl

Work in progress cURL like interface for https://github.com/rakyll/boom

```
USAGE:
   bcurl [options] <URL>
GLOBAL OPTIONS:
   -H, --header '-H option -H option'	custom header to pass to server
   -d, --data 				HTTP POST data
   --cpus '4'				Number of used cpu cores. (default for current machine is 4 cores)
   --requests '200'			Number of requests to run
   --concurrency '50'			Number of requests to run concurrently.
   --help, -h				show help
   --version, -v			print the version
```

## Installation

To install the latest [pre-built binary release](https://github.com/fgrehm/boom-curl/releases)
run the following one-liners:

##### Linux

```sh
L=$HOME/bin/bcurl && curl -sL https://github.com/fgrehm/boom-curl/releases/download/v0.1.0/linux_amd64 > $L && chmod +x $L
```

##### Mac OS

```sh
L=$HOME/bin/bcurl && curl -sL https://github.com/fgrehm/boom-curl/releases/download/v0.1.0/darwin_amd64 > $L && chmod +x $L
```

_The oneliners above assume that `$HOME/bin` is available on your `PATH`, if that's
not the case please change it accordingly._
