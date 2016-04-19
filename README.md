`snakepit-seed` is an example API built using the [snakepit](https://github.com/solher/snakepit) toolbox.

## Key features

- [Cobra](https://github.com/spf13/cobra) based CLI interface.
- Auto configuration from files/env/etcd thanks to [viper](https://github.com/spf13/viper).
- High expressiveness thanks to the [snakepit](https://github.com/solher/snakepit) toolbox.
- Easy debugging thanks to the [logrus](https://github.com/Sirupsen/logrus) logging levels and a "log everything" policy. For each request:
    - Basic logging (method, path, total latency, etc).
    - Context logging (current user, session, etc).
    - Time consuming operations auto-logging (requests/responses unmarshalling/marshalling, database requests and HTTP calls).
    - Stacktraces logging when `500` occurs.
- High integration testability thanks to loose coupling between the app interfaces (the cobra CLI, the logger and the viper config) and the app requests handler.
- Powerful and flexible request handling thanks to dynamic handlers. Controllers and business logic is built at runtime and `ctx` aware. That way, business logic and constants can for example be switched dynamically according to the current user/session/role. This also allows dependency injection without the use of any reflection.
- Swagger documentation.
- [ArangoDB](https://www.arangodb.com) multi-model (key-value, document and graph support) database.

## TODOs

* [ ] Write tests
* [ ] Add a lot of documentation

## License

MIT
