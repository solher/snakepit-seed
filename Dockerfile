FROM gliderlabs/alpine
COPY snakepit-seed /usr/bin
COPY swagger.json ~/
COPY config.yaml ~/
EXPOSE 3000
ENTRYPOINT ["snakepit-seed", "run"]
