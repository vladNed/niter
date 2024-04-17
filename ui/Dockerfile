# Base image
FROM node:20.12-buster as build-base

WORKDIR /app

COPY package*.json yarn.lock ./
RUN yarn install

COPY . .

RUN yarn build:testnet

# Stage 2
FROM nginx:alpine

COPY --from=build-base /app/build /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
