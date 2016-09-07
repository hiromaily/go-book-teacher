# go-book-teacher

[![Build Status](https://travis-ci.org/hiromaily/go-book-teacher.svg?branch=master)](https://travis-ci.org/hiromaily/go-book-teacher)
[![Coverage Status](https://coveralls.io/repos/github/hiromaily/go-book-teacher/badge.svg)](https://coveralls.io/github/hiromaily/go-book-teacher)
[![Deploy](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy?template=https://github.com/hiromaily/go-book-teacher)

Go-book-teacher is for booking specific teachers on English lesson service by web scraping.
When running on local PC, it continues to run until stop.
It notices available teachers every 2 minutes when finding and changing state. 


## Installation
```
$ go get github.com/hiromaily/go-book-teacher ./...
```

#### For docker environment
```
$ ./docker-create.sh
```


## Configration

### 1. Common settings
#### TOML file

```${PWD}/config/settings.toml```

* mail settings
* Redis settings  
※ As needed, secret information can be ciphered.(using AES encryption)

#### registration for target teacher's ids
1. Inside ./teacher/teacherinfo.go  
  or
2. Outer json file: To use command line arguments ```-f jsonfile```

#### notification
1. Web browser  
  or
2. mail: To set mail info on settins.toml

#### save current state
1. txt file: To set status_file on settings.toml  
 or
2. redis server: To set redis_url on settings.toml

### 2. On heroku
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

### 3. On Docker

#### Docker related files
* docker-create.sh
* docker-compose.yml
* docker-entrypoint.sh
* Dockerfile
* ./docker_build/*


## Environment valuable e.g.
### 1. Option
| NAME              | Value                               |
|:------------------|:------------------------------------|
| ENC_KEY           | xxxxx                               |
| ENC_IV            | xxxxx                               |
| HEROKU_FLG        | 1                                   |


## Usage
```
Usage: book [options...]

Options:
  -j     Json file path
  -t     Toml file path
  -i     Interval for scraping

e.g.
 $ book -j /var/go/teacher.json -t settings.toml -i 120
```

