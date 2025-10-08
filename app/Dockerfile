# stage builder
FROM node:22-alpine AS builder

ARG VITE_API_URL
ENV VITE_API_URL=${VITE_API_URL}

WORKDIR /app

COPY package.json yarn.lock ./
RUN yarn install --frozen-lockfile

COPY . .
RUN yarn build


# stage runner
FROM nginx:1.21.6-alpine

RUN apk add --no-cache gettext

COPY --from=builder /app/dist /usr/share/nginx/html
COPY env.js.template /usr/share/nginx/html/env.js.template

COPY env.sh /docker-entrypoint.d/env.sh
RUN chmod +x /docker-entrypoint.d/env.sh

EXPOSE 80
ENTRYPOINT ["/docker-entrypoint.d/env.sh"]
