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

