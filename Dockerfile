FROM gliderlabs/alpine
COPY snakepit-seed /usr/bin
COPY swagger.json $HOME/
COPY config.yaml $HOME/
EXPOSE 3000
ENTRYPOINT ["snakepit-seed", "run"]
