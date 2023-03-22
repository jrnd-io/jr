# JR: Quality Random Data from the Command line

JR is a CLI program that helps you to create quality random data for your applications.

## Prerequisites

JR requires golang >= 1.20


## Building and compiling

you can use the `install.sh` to install JR. This script does everything needed in one simple command.

```bash
./install.sh
```

These are the steps in the `install.sh` script:

```shell
make all
make copy_templates
sudo make install
```

If you want to run the Unit tests, you have a `make` target for that:

```bash
make test
```

## Basic usage

JR is very straightforward to use. Here are some examples

### Listing existing templates
```bash
jr list
````
Templates are in the directory `$HOME/.jr/templates`. You can override with the ```--templatePath``` command flag
Templates with issues are showed in <font color='red'>red</font>, Templates with no issues are showed in <font color='green'>green</font>

### Create random data from one of the provided templates

Use predefined net-device template to generate a random JSON network device

```bash
jr run net-device
````

### Other options for templates

If you want to use your own template, you have several options:

- put it in the default directory
- put it in another directory and use the `--templateDir` flag
- put it in another directory and use the `--templateFileName` flag to directly refer to it
- embed it directly in the command using the `--template` flag

For a quick and dirty test, you can refer a template like this:

```bash 
jr run --templateFileName ~/.jr/templates/user.tpl
```

For an even quicker and dirtier test, you can embed directly a template in the command like this:

```bash
jr run --template "name:{{name}}"
```

### Create more random data 

Using `-n` option you can create more data in each pass.

```bash
jr run net-device -n 3
```
### Continuous streaming data

Using `--f` option you can repeat the creation every `f` milliseconds

This example creates 2 net-device every second.
```bash
jr run net-device -n 2 -f 1s 
```

This example creates 2 net-device every 100ms for 1 minute.
```bash
jr run net-device -n 2 -f 100ms -d 1m 
```

Results are by default written on standard out (`--output "stdout"`) with this output template:

```
"{{.V}}\n"
```

which means that only the "Value" is in the output. You can change this behaviour with the `--outputTemplate`

## Use JR to stream data to Apache Kafka

First thing to do is to create a kafka.properties file. 

### 1a. Using Confluent Cloud & Confluent CLI

The easiest way to do that is to use [Confluent Cloud]("https://confluent.cloud/").

You can use the [confluent CLI]("https://docs.confluent.io/confluent-cli/current/overview.html") to create a cluster:

Config your vars as you see fit:
```bash
export CONFLUENT_CLUSTER_NAME=jr-test
export CONFLUENT_CLUSTER_CLOUD_PROVIDER=aws
export CONFLUENT_CLUSTER_REGION=eu-west-1 
```

Then execute the commands

```bash

confluent login --save

OUTPUT=$(confluent kafka cluster create "$CONFLUENT_CLUSTER_NAME" --cloud $CONFLUENT_CLUSTER_CLOUD_PROVIDER --region $CONFLUENT_CLUSTER_REGION --output json 2>&1)
(($? != 0)) && { echo "$OUTPUT"; exit 1; }
CONFLUENT_CLUSTER_ID=$(echo "$OUTPUT" | jq -r .id)
confluent kafka cluster use $CLUSTER 2>/dev/null
echo "Cluster $CONFLUENT_CLUSTER_NAME created, Id: $CONFLUENT_CLUSTER_ID"

confluent api-key create --resource $CONFLUENT_CLUSTER_ID

OUTPUT=$(confluent api-key create --resource $CONFLUENT_CLUSTER_ID -o json)
CONFLUENT_CLUSTER_API_KEY=$(echo "$OUTPUT" | jq -r ".api_key")
CONFLUENT_CLUSTER_API_SECRET=$(echo "$OUTPUT" | jq -r ".api_secret")

echo "API KEY:SECRET  -> $CONFLUENT_CLUSTER_API_KEY:$CONFLUENT_CLUSTER_API_SECRET"

confluent kafka topic create test --cluster $CONFLUENT_CLUSTER_ID

confluent kafka client-config create go --cluster $CONFLUENT_CLUSTER_ID --api-key $CONFLUENT_CLUSTER_API_KEY --api-secret $CONFLUENT_CLUSTER_API_SECRET 1> kafka/config.properties 2>&1
```

### 1b. Using Confluent Cloud & manually creating config file

You can also create a Cluster in [Confluent Cloud]("https://confluent.cloud/") and copy-paste the configuration in the HOME > ENVIRONMENTS > YOUR ENVIRONMENT > YOUR CLUSTER > CLIENTS > New Client section.

You can also fill the gaps in the provided `kafka/config.properties.example`

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

### 2. Writing data to Apache Kakfa

Just use the `--output kafka` flag and `--topic` flag to indicate the topic name:

```bash
jr run net-device -n 5 -f 500ms -o kafka -t test
```

If you don't specify a key, the string "key" will be used for each record. 
Using `--key` you can use a template for the key, embedding it directly in the command:

For example:
```bash
jr run -k '{{key "KEY" 20}}' -f 1s -d 10s net-device -o kafka -t test
```
Another example:
```bash 
jr run -k '{{randoms "ONE|TWO|THREE"}}' -f 1s -d 10s net-device -o kafka -t test
```
It is possible to write to both stdout and kafka at the same time:
```bash
jr run -k '{{randoms "ONE|TWO|THREE"}}' -f 1s -d 10s net-device -o stdout,kafka -t test
```
### Using JR to pipe data to **KCAT**

Another simple way of streaming to Apache Kafka is to use [kcat](https://github.com/edenhill/kcat) in conjunction with JR. 
JR supports **kcat** out of the box. Using the `--kcat` flag the standard output will be formatted with K,V on a single line.

`--kcat` it's a shorthand equivalent for `--output stdout --outputTemplate '{{key}},{{value}}' --oneline`


```bash
jr run -k '{{randoms "ONE|TWO|THREE"}}' -f 1s -d 5s net-device --kcat | kcat -F kafka/config.properties -K , -P -t test
```
