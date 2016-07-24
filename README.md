# booking-teacher

This program is for booking specific teachers on English lesson usign web scraping.



## Usage
```
Usage: book [options...]

Options:
  -f     Json file path

e.g.
 $ book -f /var/go/teacher.json
```

## Mail Mode
* Set environment variable

| NAME              |
|:------------------|
| MAIL_TO_ADDRESS   |
| MAIL_FROM_ADDRESS |
| SMTP_ADDRESS      |
| SMTP_PASS         |
| SMTP_SERVER       |
| SMTP_PORT         |


## when running on heroku
```
# initial command
$ heroku create bookteacher --buildpack heroku/go

# environment variable
$ heroku config:add HEROKU_FLG=1
$ heroku config:add MAIL_TO_ADDRESS=xxxxs@gmail.com
$ heroku config:add MAIL_FROM_ADDRESS=xxxxx@gmail.com
$ heroku config:add SMTP_ADDRESS=xxxxx@gmail.com
$ heroku config:add SMTP_PASS=xxxxx
$ heroku config:add SMTP_SERVER=smtp.gmail.com
$ heroku config:add SMTP_PORT=587

# Scheduler
$ heroku addons:create scheduler:standard

# Redis
$ heroku addons:create heroku-redis:hobby-dev -a bookteacher 
$ heroku config | grep REDIS

# Server port use below.
os.Getenv("PORT")

# deploy
$ git push -f heroku master


```