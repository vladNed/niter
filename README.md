# Niter

> NOTE: This is still under heavy development and will be updated frequently.

Niter is an open-source decentralized cryptocurrency exchange platform or DEX.

With Niter, users can trade cryptocurrencies directly with each other in a peer-to-peer manner. The platform is built on top of the WebRTC protocol, which allows for secure and private communication between users. Niter uses atomic swaps to ensure that trades are executed fairly and securely.

## Running locally

To run the entire PoC for Niter, you need to run the following commands:

```bash
docker-compose up --build
```

This command is necessary only for the first time you run the project or when changes happened to the code and you
need to rebuild the app. After that, you can simply run:

```bash
docker-compose up
```

To stop all services, you can run:

```bash
docker-compose down
```
> Note: Additionally if you want to run each container separately, you can do so by running the following commands:

For signalling server:
```bash
cd ./signalling
make build-docker
make run-docker
```