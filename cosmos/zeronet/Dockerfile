FROM node:8.9.0

WORKDIR /usr/src/app

COPY package.json .
COPY server.js .

COPY genesis.json .
COPY cyberd-0.json .
COPY cyberd-1.json .
COPY cyberd-2.json .
COPY core.js .

RUN npm install

EXPOSE 5000

CMD [ "npm", "start" ]