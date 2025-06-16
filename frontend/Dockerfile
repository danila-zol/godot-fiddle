FROM node:22.16-alpine3.22 AS build
WORKDIR /app

COPY package.json package-lock.json ./
RUN npm install

COPY . ./
RUN npm run build

FROM node:22.16-alpine3.22
WORKDIR /app

COPY --from=build /app/.output/ ./

RUN adduser -D frontrunner && \ 
		chown frontrunner .
USER frontrunner

CMD ["node", "/app/server/index.mjs"]
