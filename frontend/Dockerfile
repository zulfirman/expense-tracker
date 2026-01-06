########################################
# Builder
########################################
FROM node:20-alpine AS builder

WORKDIR /app

# 1) Install dependencies
COPY package*.json ./
RUN npm ci

# 2) Build static site
COPY . .
RUN npm run build


########################################
# Final: nginx to serve static assets + proxy API
########################################
FROM nginx:1.27-alpine

WORKDIR /usr/share/nginx/html

# Remove default config and add our own
RUN rm /etc/nginx/conf.d/default.conf
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Copy built assets
COPY --from=builder /app/build .

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]

