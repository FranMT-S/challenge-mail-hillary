
FROM node:22-alpine AS builder

ARG VITE_API_URL
ENV VITE_API_URL=${VITE_API_URL}

WORKDIR /app

COPY package.json yarn.lock ./
RUN yarn install --frozen-lockfile

COPY . .
RUN yarn build

FROM node:22-alpine AS runner

WORKDIR /app

COPY --from=builder /app/dist ./dist

RUN yarn global add serve

EXPOSE 5173
CMD ["serve", "-s", "dist", "-l", "5173"]
