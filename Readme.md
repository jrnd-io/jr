# JR: Quality Random Data from the Command line

JR is a CLI program that helps you to create quality random data for your applications.

[![img.png](images/goreport.png)](https://goreportcard.com/report/github.com/ugol/jr)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)



![JR-simple](https://user-images.githubusercontent.com/89472/229626362-70ddc95d-1090-4746-a20a-fbffba4193cd.gif)

## Building and compiling

JR requires Go 1.20

you can use the `make_install.sh` to install JR. This script does everything needed in one simple command.

```bash
./make_install.sh
```

These are the steps in the `make_install.sh` script:

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

## Using embedded function documentation to write your own templates

![JR-man](https://user-images.githubusercontent.com/89472/229628592-68619ec7-2b1d-4704-8c76-ba59bb82579d.gif)

JR has plenty of embedded functions that can be used to write yor own templates.
We have included the documentation for all the functions directly into JR.

You can list all the available functions with a simple command:
```bash
jr man -l
```

You can filter by category:
```bash
jr man -c net
```
Or you can filter by name and description:
```bash
jr man -f random
```
You can also execute directly the Example using `-r` flag:

```bash
jr man ip -r
```
which will basically execute this command for you:

```bash
 jr run --template '{{ip "10.2.0.0/16"}}'
 ```
To study more advanced usages, look at the templates in your `templates` directory


## Use JR to stream data to Apache Kafka

First thing to do is to create a Kafka cluster and relative kafka.properties file. The easiest way to do that is to use [Confluent Cloud]("https://confluent.cloud/").

Here we document three different ways of doing that. Choose the one that fits you better!

### 1. Confluent Cloud and downloading the config file

Just create a basic (free!) Cluster with the web console in [Confluent Cloud]("https://confluent.cloud/") and 
copy-paste the configuration in the HOME > ENVIRONMENTS > YOUR ENVIRONMENT > YOUR CLUSTER > CLIENTS > New Client section.

### 2. Confluent Cloud and config file via Confluent CLI

You can use the [confluent CLI]("https://docs.confluent.io/confluent-cli/current/overview.html") to create a cluster and 
the configuration in a programmatic way:

Config your vars as you see fit, for example:
```bash
export CONFLUENT_CLUSTER_NAME=jr-test
export CONFLUENT_CLUSTER_CLOUD_PROVIDER=aws
export CONFLUENT_CLUSTER_REGION=eu-west-1 
```
Then execute the following commands

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

### 3 An existing Kafka cluster & manually creating config file

If you have an existing cluster, just fill the fields in the provided `kafka/config.properties.example`

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
### Writing data to Apache Kakfa

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
### Autocreate topics

Topics autocreation is disabled in Confluent Cloud. 
If you are really lazy you can use the `-a` option, so JR will create the topic for you. 

```bash
jr run -a -k '{{randoms "ONE|TWO|THREE"}}' -f 1s -d 10s net-device -o stdout,kafka -t mynewtopic
```

Alternatively, you can also create it explicitly from JR:

```bash
jr createTopic topic1
```
If you want to specify number of partitions and replication Factor you can use the `-p` and `-r` flags:

```bash
jr createTopic topic1 -p 10 -r 2
```

### Confluent Schema Registry support

There is also support for Confluent Schema Registry. 
At the moment only `json-schema` and `avro-generic` is directly supported.

To use Confluent Schema registry you need first to fill the `registry.properties` provided example with the needed link and user/pwd:

```properties
schemaRegistryURL=https://blabla.europe-west3.gcp.confluent.cloud
schemaRegistryUser=blablabla-saslkey
schemaRegistryPassword=blablabla-saslpwd
```
then use the `--schema` and the `--serializer` flags

Example usage:
```bash
jr run user -o kafka -t topic1 -s --serializer avro-generic
```
or 
```bash
jr run net-device -o kafka -t topic2 -s --serializer json-schema
```
Remember that once you run these commands, `topic1` will be associated with an avro generic schema representing an user 
object, and `topic2` with a json-schema representing a net-device object. 



### Using JR to pipe data to **KCAT**

Another simple way of streaming to Apache Kafka is to use [kcat](https://github.com/edenhill/kcat) in conjunction with JR. 
JR supports **kcat** out of the box. Using the `--kcat` flag the standard output will be formatted with K,V on a single line.

`--kcat` it's a shorthand equivalent for `--output stdout --outputTemplate '{{key}},{{value}}' --oneline`


```bash
jr run -k '{{randoms "ONE|TWO|THREE"}}' -f 1s -d 5s net-device --kcat | kcat -F kafka/config.properties -K , -P -t test
```


## Docker 

### Multi-Arch Build 

```bash
# Create the local builder 
docker buildx create --name local --bootstrap --use
# Local build 
docker buildx build --platform linux/arm64/v8,linux/amd64 --output=local -t jr:latest .
# Push on DockerHub
docker buildx build --platform linux/arm64/v8,linux/amd64  --build-arg=USER="$(whoami)" --build-arg="0.1.0"  --push -t ugol/jr:latest .
```

### How to use jr with local configurations

It is possible to mount config files from your local environment and use them with jr docker image.

```
docker run -it -v $(pwd)/configs:/home/jr-user/configs --rm ugol/jr:latest jr run net-device -n 5 -f 500ms -o kafka -t net-device -F /home/jr-user/configs/kafka.client.properties -s --serializer json-schema --registryConfig /home/jr-user/configs/registry.client.properties
```
![docker](https://user-images.githubusercontent.com/89472/230502463-cb6faaf8-fcf1-48c4-a571-031d46725cc1.gif)
