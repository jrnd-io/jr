# JR: streaming Quality Random Data from the Command line

JR is a CLI program that helps you to stream quality random data for your applications.

![jr](https://user-images.githubusercontent.com/89472/235927141-87632730-90d6-469f-97b0-8b638077dd4e.png)


[![img.png](images/goreport.png)](https://goreportcard.com/report/github.com/ugol/jr)
![Build](https://github.com/ugol/jr/actions/workflows/go-linux.yml/badge.svg)
![Build](https://github.com/ugol/jr/actions/workflows/go-mac.yml/badge.svg)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Reference](https://pkg.go.dev/badge/github.com/ugol/jr.svg)](https://pkg.go.dev/github.com/ugol/jr)
[![Docker](https://img.shields.io/badge/docker-latest-blue.svg)](https://hub.docker.com/r/ugol/jr)

![JR-simple](https://user-images.githubusercontent.com/89472/229626362-70ddc95d-1090-4746-a20a-fbffba4193cd.gif)

## Documentation

For full documentation about emitters, referential integrity, how to write templates and more, pls see the full [JR Documentation](https://jrnd.io).

## Building and compiling

JR requires Go 1.22

you can use the `make_install.sh` to install JR. This script does everything needed in one simple command.

```bash
./make_install.sh
```

These are the steps in the `make_install.sh` script if you want to use them separately:

```shell
make all
make copy_templates
sudo make install
```

If you want to run the Unit tests, you have a `make` target for that too:

```bash
make test
```

## Basic usage

JR is very straightforward to use. Here are some examples:

### Listing existing templates
```bash
jr template list
````
Templates are in the directory `$JR_HOME/templates`. JR_HOME defaults to `~/.jr` and can be changed to a different dir, for example:

```bash
JR_HOME=~/jrconfig/ jr template list
````

Templates with parsing issues are showed in <font color='red'>red</font>, Templates with no parsing issues are showed in <font color='green'>green</font>

### Create random data from one of the provided templates

Use for example the predefined `net_device` template to generate a random JSON network device

```bash
jr template run net_device
````

or, with a shortcut:

```bash
jr run net_device
````

### Using Docker

You can also use a [![Docker](https://img.shields.io/badge/docker-latest-blue.svg)](https://hub.docker.com/r/ugol/jr)
image if you prefer.

```bash
docker run -it ugol/jr:latest jr run net_device
```

### Other options for templates

If you want to use your own template, you can:

- put it in the templates directory
- embed it directly in the command using the `--embedded` flag

For a quick and dirty test, the best option is to embed directly a template in the command:

```bash
jr run --template "name:{{name}}"
```

### Create more random data 

Using `-n` option you can create more data in each pass. 
This example creates 3 net_device objects at once:

```bash
jr run net_device -n 3
```
### Continuous streaming data

Using `--frequency` option you can repeat the creation every `f` milliseconds

This example creates 2 net_device every second, for ever:

```bash
jr run net_device -n 2 -f 1s 
```

Using `--duration` option you can time bound the entire object creation.

This example creates 2 net_device every 100ms for 1 minute:

```bash
jr run net_device -n 2 -f 100ms -d 1m 
```

Results are by default written on standard out (`--output "stdout"`) with this output template:

```
"{{.V}}\n"
```

which means that only the "Value" is in the output. You can change this behaviour embedding a different template with `--outputTemplate`

If you want syntax colouring and your output is just json, you can pipe to [jq](https://jqlang.github.io/jq/)

```bash
jr run net_device -n 2 -f 100ms -d 1m | jq
```

Beware that if you, for example, include the key in the output, it won't be possible to use jq:

```bash
jr run net_device -n 2 -f 100ms -d 1m --kcat | jq

parse error: Expected value before ',' at line 1, column 5
```
