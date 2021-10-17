FROM alpine:latest

ADD citadel_intranet /main
ADD web /var/www
ADD migrations /var/migrations

CMD "/main"
