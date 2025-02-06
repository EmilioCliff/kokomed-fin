# Build Stage
FROM node:23-alpine AS build
WORKDIR /app
COPY package*.json ./
RUN npm install --frozen-lockfile
COPY . .
RUN npm run build

EXPOSE 3000
CMD ["npx", "vite", "preview", "--host"]
