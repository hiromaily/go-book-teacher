# go-book-teacher

[![Build Status](https://travis-ci.org/hiromaily/go-book-teacher.svg?branch=master)](https://travis-ci.org/hiromaily/go-book-teacher)
[![Coverage Status](https://coveralls.io/repos/github/hiromaily/go-book-teacher/badge.svg)](https://coveralls.io/github/hiromaily/go-book-teacher)
[![Go Report Card](https://goreportcard.com/badge/github.com/hiromaily/go-book-teacher)](https://goreportcard.com/report/github.com/hiromaily/go-book-teacher)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](https://raw.githubusercontent.com/hiromaily/go-book-teacher/master/LICENSE)
[![Deploy](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy?template=https://github.com/hiromaily/go-book-teacher)

Go-book-teacher is notifier for that specific teachers are available on English lesson service.  
For now, only [DMM eikaiwa](https://eikaiwa.dmm.com/) is expected.

This project has started since 2016 to study Golang and code is quite messy. Now it's under refactoring.

#### Console
```
----------- Rita M / Portugal / 7093 -----------
----------- Gagga / Serbia / 5252 -----------
----------- Kaytee / Serbia / 7646 -----------
----------- Milica J / Serbia / 6294 -----------
----------- Marine / France / 8519 -----------
----------- Lavinija / Serbia / 5656 -----------
2017-07-06 20:30:00
2017-07-07 16:00:00
2017-07-07 18:00:00
2017-07-07 18:30:00
----------- Aleksandra S / Serbia / 6214 -----------
----------- Jekica / Serbia / 4806 -----------
----------- Yovana / Serbia / 6550 -----------
2017-07-07 03:30:00
2017-07-07 04:00:00
2017-07-07 04:30:00
----------- Olivera V / Serbia / 5380 -----------
----------- Emilia / Serbia / 2464 -----------
----------- Indre / Lithuania / 3486 -----------
----------- Joxyly / Serbia / 4808 -----------
----------- Milica Ml / Serbia / 4107 -----------
```


#### Slack
![slack](https://raw.githubusercontent.com/hiromaily/go-book-teacher/master/images/slack_image.png)


## Requirements
- Golang 1.15+
- Docker compose
  - Redis
- [direnv](https://github.com/direnv/direnv) for MacOS user


## Installation
### for MacOS user
```
$ brew install hiromaily/tap/go-book-teacher

# run
$ book-teacher
```

### for development
```
# clone
$ git clone https://github.com/hiromaily/go-book-teacher.git

# setup config file
$ cp configs/default.example.toml configs/default.toml
# and modify `configs/default.toml` as you want

# setup your favorite teacher's information
$ cp configs/teacher/default.example.json configs/teacher/default.json
# and modify `configs/teacher/default.json` as you want

# setup .envrc
$ cp example.envrc .envrc
# and modify `.envrc` as you want

# build
$ make build

# run
book
```

## Usage
```
Usage: book [options...]

Options:
  -json      Json file path for teacher information
  -toml      Toml file path for config
  -day       range of schedule to get teacher's availability: 0: all day, 1:today, 2: tomorrow
  -v         show version

e.g.
 $ book -day 1
 $ book -h      # show help
```

## Environment valuables
- envrc is used for MacOS user, please install direnv](https://github.com/direnv/direnv)
- encryption is used for secret value in config files.
- secret value can be encrypted/decrypted by tools (see Makefile)

```Makefile
.PHONY: tool-encode
tool-encode:
	go run ./tools/encryption/ -encode important-password

.PHONY: tool-decode
tool-decode:
	go run ./tools/encryption/ -decode o5PDC2aLqoYxhY9+mL0W/IdG+rTTH0FWPUT4u1XBzko=
```

| NAME              | Value                                         |
|:------------------|:----------------------------------------------|
| GO_BOOK_CONF      | default config file path                      |
| GO_BOOK_JSON      | default teacher info json file path           |
| ENC_KEY           | 16byte_string_xx                              |
| ENC_IV            | 16byte_string_xx                              |


## deploy on heroku
```
## Install 
$ heroku create bookteacher --buildpack heroku/go
$ heroku addons:create heroku-redis:hobby-dev -a projectname 
$ heroku addons:create scheduler:standard

## Environment variable
$ heroku config:add ENC_KEY=xxxxx
$ heroku config:add ENC_IV=xxxxx

## Check
$ heroku config | grep REDIS
$ heroku ps -a bookteacher

## Deploy
$ git push -f heroku master

```
