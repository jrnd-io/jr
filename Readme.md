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
make copy_templates; sudo make install
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
### Continuous streaming data

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

Just use the -t flag to indicate the topic name:

```bash
jr run net-device -n 5 -f 500ms -t test
```
With ```silent``` mode you can stop the standard output:

```bash
jr run net-device -n 5 -f 500ms -s -t test
```

If you don't specify a key, the string "key" will be used for each record. 
Using ```--key``` you can use a template for the key, embedding it directly in the command:

For example:
```bash
jr run -k '{{key "KEY" 20}}' -f 1s -d 10s net-device -s -t test
```
Another example:
```bash 
jr run -k '{{randoms "ONE|TWO|THREE"}}' -f 1s -d 10s net-device -s -t test
```

### Using JR to pipe data to **KCAT**

Another simple way of streaming to Apache Kafka is to use [kcat](https://github.com/edenhill/kcat) in conjunction with JR. 
JR supports **kcat** out of the box. Using the ```--kcat``` flag the standard output will be formatted with K,V on a single line

```bash
jr run -k '{{randoms "ONE|TWO|THREE"}}' -f 1s -d 5s net-device --kcat | kcat -F kafka/config.properties -K , -P -t test
```
