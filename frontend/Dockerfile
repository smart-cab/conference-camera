FROM node:21-alpine3.18 as builder

WORKDIR /app

COPY package.json package-lock.json .

# Behaves like npm install, except that it installs packages according to
# package-lock.json instead of package.json (making more stable and reproducible
# builds, especially for production-releases). If it fails, run npm install on
# your local machine before building within docker. Equivalent of yarn --frozen-lockfile
RUN npm ci

COPY . .

RUN npm run build


FROM joseluisq/static-web-server:2-alpine

WORKDIR /app

COPY --from=builder /app/build /app /app/.env.production

EXPOSE 8787

CMD ["static-web-server", "-p", "8787", "--root", "/app"]
