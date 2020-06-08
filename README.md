# Trading212 Trading Bot

A trading bot created for [Trading 212](https://www.trading212.com) platform using [Selenium](https://www.selenium.dev/) framework.

It was created for [Upwork](https://www.upwork.com/jobs/Golang-Job-mini-backend-project-using-Selenium-similar_~015a9caf0b60510928) project.

Currently app is in demo version, with minimal features implemented.

## Dependencies
-----------------------------------------------------

**Note:** skip this section if you are going to use Docker.

1. Download [Go](https://www.golang.org) and install it.
2. Download [Selenium standalone server](https://www.selenium.dev/downloads/), run it: `java -jar server.jar`, you will need Java8.0+ version to be able to run it, alternatively you can use already downloaded one in selenium folder.

## Build
-----------------------------------------------------

### Manual run

1.  Cd to repo folder and run: `chmod +x build.sh; ./build.sh`, it will create a program executable in bin/ folder.
2.  You need Trading212 username/email and password to be able to run the program.
3.  There is a config/conf.ini file, you can put your credentials there or you can pass your credentials as program arguments. If you cd to bin folder and run: `./web --help`, you can view additional arguments for program or you can input them in conf.ini file.
4.  If you filled config/conf.ini with your creds, you can run program: `./web --inifile` and it will start listening on specified port(default 4000).

### Docker run

1. Pull Selenium images from Docker hub, by running following command: 
    `docker pull selenium/hub:latest selenium/base:latest selenium/standalone-chrome:late selenium/standalone-firefox:latest`
2. Run Selenium standalone server as a container (you can choose firefox or chrome):
    `docker run -d -p 4444:4444 -v /dev/shm:/dev/shm selenium/standalone-firefox:latest`
3. Build Trading212 Trading Bot image:
    `docker build -t trading212 .`
4. Run Bot program as a container(you have to input same browser in conf.ini as Selenium standalone server container), it will start listening on specified port(default 4000):
    `docker run --network "host" trading212`

## Description
-----------------------------------------------------

For the purposes of demo program, there is just 1 route coded for now. Using this route, you can buy commodity of your choice at market price, limit/stop is not available for now(read Known Issues). I recommend you to test program using Postman. Just create a new POST request there, input http://127.0.0.1:4000/commodity as a url and finally choose Body -> raw and input following json:

JSON request syntax:
```
    {
        "name": "Ethereum" (string), // should be full name of commodity
        "quantity": 3.4 (float), // amount you want to buy
        "order": "Buy" (string) // type of order, buy or sell(for now just buy is implemented)
        "price": 200.10 (float) // not implemented yet
    }
```

JSON response syntax:
```
    {
        "success": true/false (boolean), // indication of request success
        "error": "" (string), // usually empty or server error message
        "data": "" (interface) // output of Alert box after order is done
    }
```

### List of implemented routes

* http://127.0.0.1:4000/commodity - POST

## Known issues
-----------------------------------------------------

Buying commodity using limit/stop order is not implemented yet, because of famous Selenium issue: element couldn't be scrolled into view. This issue is going to be fixed, as soon as project gets approval.