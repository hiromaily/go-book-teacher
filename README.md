# go-book-teacher

[![Build Status](https://travis-ci.org/hiromaily/go-book-teacher.svg?branch=master)](https://travis-ci.org/hiromaily/go-book-teacher)
[![Coverage Status](https://coveralls.io/repos/github/hiromaily/go-book-teacher/badge.svg)](https://coveralls.io/github/hiromaily/go-book-teacher)
[![Go Report Card](https://goreportcard.com/badge/github.com/hiromaily/go-book-teacher)](https://goreportcard.com/report/github.com/hiromaily/go-book-teacher)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](https://raw.githubusercontent.com/hiromaily/go-book-teacher/master/LICENSE)
[![Deploy](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy?template=https://github.com/hiromaily/go-book-teacher)

Go-book-teacher is notifier for specific teachers are available on English lesson service.  


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


## Installation
```
$ go get github.com/hiromaily/go-book-teacher ./...
```

## Usage
```
Usage: book [options...]

Options:
  -j      Json file path for teacher information
  -t      Toml file path for config
  -i      Interval for scraping, if 0 it scrapes only once
  -crypto true is that conf file is handled as encrypted value

e.g.
 $ book -j ./testdata/json/teachers.json -t ./config/toml/text-command.toml
```

## Configration
```
./config/toml/*.toml
```
* site
* storage
    * redis
    * text
* notification
    * slack
    * browser
    * mail

※ As needed, secret information can be encrypted.(using AES encryption)

## Environment valuables
- encryption is used for secret value in config files.
- secret value can be encrypted/decrypted by tools
```bash
# encode
go run ./tools/encryption/ -m e important-password
# decode
go run ./tools/encryption/ -m d o5PDC2aLqoYxhY9+mL0W/IdG+rTTH0FWPUT4u1XBzko=
```

### Option
| NAME              | Value                               |
|:------------------|:------------------------------------|
| ENC_KEY           | xxxxx                               |
| ENC_IV            | xxxxx                               |


## registration for target teacher's ids
json file can be used with command line argument `-j`
`testdata/json/teachers.json`


## deploy on heroku
```
## Install 
$ heroku create bookteacher --buildpack heroku/go
$ heroku addons:create heroku-redis:hobby-dev -a projectname 
$ heroku addons:create scheduler:standard

## Environment variable
$ heroku config:add HEROKU_FLG=1
$ heroku config:add ENC_KEY=xxxxx
$ heroku config:add ENC_IV=xxxxx

## Check
$ heroku config | grep REDIS
$ heroku ps -a bookteacher

## Deploy
$ git push -f heroku master

```
