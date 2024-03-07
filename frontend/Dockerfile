FROM node:21-alpine3.18

WORKDIR /app

COPY package.json package-lock.json .

# Behaves like npm install, except that it installs packages according to
# package-lock.json instead of package.json (making more stable and reproducible
# builds, especially for production-releases). If it fails, run npm install on
# your local machine before building within docker. Equivalent of yarn --frozen-lockfile
RUN npm ci

COPY . .

EXPOSE 3000

CMD ["npm", "run", "start"]
