# weblogger
Web server for logging data.

## Usage

    ./weblogger PORT TOKEN [PATH]

where 

- `PORT` is the port where the server is started.
- `TOKEN` is the token required to write data to a datafile.
- `PATH` is the path where data files are stored. By default, the path containing the `weblogger` binary is used.

## Writing data

The data is written using a POST request. If the token is `1234` and the server is started on `localhost`, the following HTTP request writes `SOME_DATA` to a timestamped csv file:

    POST /?token=1234&data=SOME_DATA HTTP/1.1
    Host: localhost

## Reading data

If the server is started on `localhost:8080`, the data is simply read by navigating a webbrowser to

    http://localhost:8080

To get data from another date, (e.g. 2019-02-01) use the following url:

    http://localhost:8080?date=2019-02-01

## Deploying the server

On a linux server with systemd, generate the file `/etc/systemd/system/weblogger-sensor.service` with the content:

    [Unit]
    Description=Weblogger
    After=network.target
    StartLimitIntervalSec=0

    [Service]
    Type=simple
    Restart=always
    RestartSec=1
    User=<<USERNAME>>
    ExecStart=<</PATH/TO/>>weblogger <<PORT>> <<TOKEN>>

    [Install]
    WantedBy=multi-user.target

where the tags marked with `<<>>` must be filled out, of course. To start the server, run

    systemctl start weblogger-sensor

and to do so on boot, run

    systemctl enable weblogger-sensor

Note that the port must be opened to the outside world to access it from outside the server.

