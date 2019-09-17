#!/bin/bash

TOKEN=1155234
DATA=1,2,3

cat << EOF | nc localhost 8080
POST /?token=$TOKEN&data=$DATA HTTP/1.1
Host: localhost


EOF

echo ""

