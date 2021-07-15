# Simple website

## Setup

- golang 1.16+
- make
- docker
- [gauge](https://docs.gauge.org/getting_started/installing-gauge.html?os=macos&language=javascript&ide=vscode)
- heroku

```bash
# for generating unit tests mocks
## install mockgen
go install github.com/golang/mock/mockgen@v1.6.0
## setup path to the installed mockgen if needed
export PATH=$PATH:$(go env GOPATH)/bin
```

```bash
# for automated tests
## install gauge
https://docs.gauge.org/getting_started/installing-gauge.html
```

## Documentation

1. Make command to test
```bash
# generate mocks
make generate

# test
make test

# run automated test against local web server
make local-automated-test

# clean up
make local-clean
```

2. API documentation [swagger doc format](api.yaml)

3. Production

Setup and deploy
```bash
# setup heruko app in databashboard then add remote
heroku git:remote -a desmond-reap-backend-challenge

# setup postgres database
heroku addons:create heroku-postgresql:hobby-dev

# manually run flyway on heroku postgres database
# can't find any buildpack for this
docker run -v ${PWD}/.db:/flyway/sql -ti flyway/flyway:7.11.1-alpine -url=jdbc:postgresql://${DATABASE_URL} -user=<username copied from heroku data> -password=<password copied from heroku data> -connectRetries=60 migrate

# deploy to heroku 
git push heroku main
```

Example curls
```bash
# register / login
curl -v http://desmond-reap-backend-challenge.herokuapp.com/v1/auth/login --data "{ \"username\": \"somebody\", \"password\": \"password\" }"

# list posts
curl -v -H"Authorization: Bearer ${token returned in login}" http://reap-backend-challenge.herokuapp.com/v1/user/post
```