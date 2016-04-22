FROM gliderlabs/alpine
COPY snakepit-seed /usr/bin
COPY swagger.json /
EXPOSE 3000
ENTRYPOINT ["snakepit-seed", "run"]
