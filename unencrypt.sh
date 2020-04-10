#
# !/bin/sh
#

openssl rsa -in key.pem -out key.unencrypted.pem -passin pass:$1
