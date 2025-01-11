# cendit

Crypto backed Banking - better still - KUDA for cryptocurrencies

## Architecture

Cendit is built on golang's multi-module workspaces setup and implements a service oriented approach for resource sharing and service distribution - `one codebase, one server, multiple services`.

Who needs a queue when you can just `go build`? Should we need a queue? Then it definitely won't be for cross service communication.

Known services in this design are:

- [X] gate - access translator
- [X] auth - authentication, users, kyc, etc
- [X] chain - cryptocurrency, blockchain
- [X] bank - core banking, bill payments, cards etc
- [X] connect - third party communication
- [X] signal - events, alerts, broadcasts, emails, sms, in-app, push, etc

## Spinning up

Add the requried environment variables into a `.env` file in the `gate`  folder and execute `make run`.

## Resources on gqlgen

- [implement multiple resolver files](https://github.com/99designs/gqlgen/issues/1427)
- [how to configure gqlgen](https://gqlgen.com/config/)
- [authentication with gqlgen](https://gqlgen.com/recipes/authentication/)
