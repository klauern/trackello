# trackello
[![Go Report Card](https://goreportcard.com/badge/github.com/klauern/trackello)](https://goreportcard.com/report/github.com/klauern/trackello)
[![codecov](https://codecov.io/gh/klauern/trackello/branch/master/graph/badge.svg)](https://codecov.io/gh/klauern/trackello)
[![wercker status](https://app.wercker.com/status/0d8fcf888b7cf0be0f4e0678b179c341/m "wercker status")](https://app.wercker.com/project/bykey/0d8fcf888b7cf0be0f4e0678b179c341)

Tracking your work in Trello from the command-line

# Setup

## Pre-configured Go environment

If you are proficient in developing and building Go applications, you can simply install with the following command:

```sh
go get github.com/klauern/trackello/cmd/trackello
```

## Make file

An alternative setup is to use the Makefile:

```sh
make setup
make
make trackello
./bin/trackello # For running
```

# Motivation

There's a decent set of blog posts that I have been working on to explain the rationale:

* [Trackello](http://blog.nickklauer.info/2016/trackello/)
* [Formatting and Parallelism](http://blog.nickklauer.info/2016/trackello_parallelism/)

# Future work

- [ ] Display Key identifying summary fields
- [ ] Add color-coding for labels
- [ ] lots more tests
- [ ] create helper to configure the Trello API/Token/Board settings on first run