FROM node:10-alpine

COPY package.json yarn.lock /usr/src/app/
WORKDIR /usr/src/app
RUN yarn
COPY . /usr/src/app/
COPY --from=nlepage/grpc-hellow:service /hellow.proto /usr/src/app/

ENTRYPOINT ["yarn", "start"]
