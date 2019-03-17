openssl genrsa -out private.pem 2048
openssl pkcs8 -topk8 -inform pem -in private.pem -outform pem -nocrypt -out private_pkcs8.pem
openssl rsa -in private.pem -outform PEM -pubout -out public.pem