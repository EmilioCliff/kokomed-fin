# Build Stage
FROM node:23-alpine AS build
WORKDIR /app
COPY package*.json ./
RUN npm install --frozen-lockfile
COPY . .
RUN npm run build

# Production Stage
FROM node:18
WORKDIR /app
RUN npm install -g serve
COPY --from=builder /app/dist /app/dist
EXPOSE 3000
CMD ["serve", "-s", "dist", "-l", "3000"]

