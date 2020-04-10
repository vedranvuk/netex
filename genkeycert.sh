#
# !/bin/sh
#

openssl req -x509 -newkey rsa:4096 -days 9999 -keyout key.pem -out cert.pem
