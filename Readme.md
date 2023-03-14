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
- put it in another directory and use the ```--templatePath``` flag to directly refer to it
- embed it directly in the command using the ```--template``` flag

For a quick and dirty test, you can refer a template like this:

```jr --templatePath ./templates/user.json run ```

For an even quicker and dirtier test, you can embed directly a template in the command like this:

```jr --template "name:{{name}}" run```


### Create more random data 

Using ``` --n ``` option you can create more data

```bash
jr run net-device --n 3
```
### Continuos streaming data

Using ``` --f ``` option you can repeat the creation every ```f``` milliseconds

```bash
jr run net-device --n 2 --f 500 
```
### Use JR to stream data to Apache Kafka

A simple way of streaming to Apache Kafka is to use kcat in conjunction with JR.
kcat needs K,V to be on a single line, so if your template generates multiline data you have to use the ```oneline``` 
option to strip all newlines. The alternative is obviously to create a template without newlines, but that's not very readable!

The following line generates 5 net-device random data every half-second and writes them to topic test:

```bash
jr run net-device --n 5 --f 500 --oneline | kcat -T -F kcat/librdkafka.config -K , -P -t test
```
