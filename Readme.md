# JR: Quality Random Data from the Command line

JR is a CLI program that helps you to create quality random data for your applications.

## Prerequisites

JR requires golang >= 1.20


## Building and compiling

you need just to use the provided Makefile to build JR

```bash
make all
```

If you want to run the Unit tests, run:

```bash
make test
```

To install the binary and the templates:
```bash
sudo make install
```

## Basic usage

JR is very straightforward to use. Here are some examples

### Listing existing templates
```bash
jr list
```
Templates are in the directory ```$HOME/.jr/templates```. You can override with the ```--templatePath``` command flag

### Create random data from one of the provided templates

Use predefined net-device template to generate a random JSON network device

```bash
jr run net-device
```

### Other options for templates

If you want to use your own template, you have several options:

- put it in the default directory
- put it in another directory and use the ```--templateDir``` flag
- put it in another directory and use the ```--templateFileName``` flag to directly refer to it
- embed it directly in the command using the ```--template``` flag

For a quick and dirty test, you can refer a template like this:

```bash 
jr run --templateFileName ~/.jr/templates/user.json
```

For an even quicker and dirtier test, you can embed directly a template in the command like this:

```bash
jr run --template "name:{{name}}"
```


### Create more random data 

Using ``` -n ``` option you can create more data in each pass.

```bash
jr run net-device -n 3
```
### Continuos streaming data

Using ``` --f ``` option you can repeat the creation every ```f``` milliseconds

This example creates 2 net-device every second.
```bash
jr run net-device -n 2 -f 1s 
```

This example creates 2 net-device every 100ms for 1 minute.
```bash
jr run net-device -n 2 -f 100ms -d 1m 
```

### Use JR to stream data to Apache Kafka

A simple way of streaming to Apache Kafka is to use kcat in conjunction with JR.
First thing to do is to create a kafka.properties file. The easiest way to do that is to use [Confluent Cloud]("https://confluent.cloud/") and copy-paste 
the configuration in the HOME > ENVIRONMENTS > YOUR ENVIRONMENT > YOUR CLUSTER > CLIENTS > New Client section.

You can also fill the gaps in the provided ```kafka/config.properties.example```

```properties
# Kafka configuration
# https://github.com/confluentinc/librdkafka/blob/master/CONFIGURATION.md

bootstrap.servers=
security.protocol=SASL_SSL
sasl.mechanisms=PLAIN
sasl.username=
sasl.password=
compression.type=gzip
compression.level=9
# statistics.interval.ms=1000
```
kcat needs K,V to be on a single line, so if your template generates multiline data you have to use the ```oneline``` 
option to strip all newlines. The alternative is obviously to create a template without newlines, but that's not very readable!

The following line generates 5 net-device random data every half-second and writes them to topic test:

```bash
jr run net-device -n 5 -f 500ms -o | kcat -T -F kafka/config.properties -K , -P -t test
```

You can do the same thing with jr, just use the -t flag to indicate teh topic name:

```bash
jr run net-device -n 5 -f 500ms -t test
```

With ```silent``` mode you can stop the standard output. Do not do this when piping to kcat or the pipe won't work!:

```bash
jr run net-device -n 5 -f 500ms -s -t test
```