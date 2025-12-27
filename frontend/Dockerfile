# Stage 1: Build
FROM node:18-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

# Stage 2: Production Server
FROM nginx:stable-alpine
COPY --from=builder /app/build /usr/share/nginx/html
# If you use React Router, you'd need a custom nginx config here
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]