FROM debian:stretch-slim

WORKDIR /

COPY nodereservasion /usr/local/bin

CMD ["nodereservasion"]