FROM gliderlabs/alpine
RUN useradd -m app
USER app
COPY snakepit-seed /usr/bin
COPY swagger.json /home/app/
COPY config.yaml /home/app/
EXPOSE 3000
ENTRYPOINT ["snakepit-seed", "run"]
