# deepwork

Inspired by the book [Deep Work](http://calnewport.com/books/deep-work/) from Cal Newport, deepwork is a little command line tool to help you enter a focused work state by closing all communication apps.

Currently only working with Mac OS X, more variety to come soon.

## Set-Up

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