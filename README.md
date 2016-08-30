# go-book-teacher

[![Deploy](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy?template=https://github.com/hiromaily/go-book-teacher)

[![Build Status](https://travis-ci.org/hiromaily/go-book-teacher.svg?branch=master)](https://travis-ci.org/hiromaily/go-book-teacher)

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
#### registration for target teacher's ids
1. Inside ./teacher/teacherinfo.go  
  or
2. Outer json file: To use command line arguments ```-f jsonfile```

#### notification
1. Web browser  
  or
2. mail: To set ```MAIL_TO_ADDRESS``` environment variable

#### save current state
1. txt file: ```SAVE_LOG``` environment variable or ```/tmp/status.log``` as default  
 or
2. redis server: To set ```REDIS_URL``` environment variable

### 2. On heroku
```
## Install 
$ heroku create bookteacher --buildpack heroku/go
$ heroku addons:create heroku-redis:hobby-dev -a projectname 
$ heroku addons:create scheduler:standard

## Environment variable
$ heroku config:add HEROKU_FLG=1
$ heroku config:add MAIL_TO_ADDRESS=xxx@gmail.com
$ heroku config:add MAIL_FROM_ADDRESS=xxx@gmail.com
$ heroku config:add SMTP_ADDRESS=xxx@gmail.com
$ heroku config:add SMTP_PASS=xxxxx
$ heroku config:add SMTP_SERVER=smtp.gmail.com
$ heroku config:add SMTP_PORT=587

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
| REDIS_URL         | redis://h:password@$servername:6379 |
| SAVE_LOG          | /var/log/teacher/data.log           |
| HEROKU_FLG        | 1                                   |

### 2. Mail Mode
| NAME              | Value            |
|:------------------|:---------------- |
| MAIL_TO_ADDRESS   | xxx@example.com  |
| MAIL_FROM_ADDRESS | info@example.com |
| SMTP_ADDRESS      | info@example.com |
| SMTP_PASS         | xxxxx            |
| SMTP_SERVER       | smpt.example.com |
| SMTP_PORT         | 587              |

### 3. Port
Heroku server use PORT automatically as environment valuable

| NAME              | Value            |
|:------------------|:---------------- |
| PORT              | 9999             |


## Usage
```
Usage: book [options...]

Options:
  -f     Json file path

e.g.
 $ book -f /var/go/teacher.json
```

