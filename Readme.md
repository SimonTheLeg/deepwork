# deepwork

Inspired by the book [Deep Work](http://calnewport.com/books/deep-work/) from Cal Newport, deepwork is a little command line tool to help you enter a focused work state by closing all communication apps.

Currently only working with Mac OS X, more variety to come soon.

## Installation

Simply grab the desired version from the Github Release page and place it in your $PATH, e.g.

```shell
curl -L https://github.com/SimonTheLeg/deepwork/releases/download/v0.1.0/deepwork-darwin-64 -o /usr/local/bin/deepwork
```


## Configuration

create a config file under ~/.deepwork/config.json and add the names of all communication applications. An example config can be found in [example-config.json](example-config.json). By default Mail and Calendar will be added.

## Usage

Close all communication apps to enter a focused working state.

```shell
deepwork on
```

Open up all communication apps again

```shell
deepwork off
```