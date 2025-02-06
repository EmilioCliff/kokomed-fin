FROM node:23-alpine3.20 AS build
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm install
RUN npm run build
COPY . .

EXPOSE 3000
CMD ["npm", "run", "preview"]
 