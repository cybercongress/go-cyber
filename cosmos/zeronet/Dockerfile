FROM node:8.9.0

WORKDIR /usr/src/app

COPY package.json .
COPY server.js .

COPY genesis.json .
COPY privkey0.json .
COPY privkey1.json .
COPY privkey2.json .
COPY core.js .

RUN npm install

EXPOSE 5000

CMD [ "npm", "start" ]