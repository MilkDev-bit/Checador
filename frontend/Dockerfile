FROM node:20-alpine AS builder

WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
# nginx:alpine 1.19+ automatically runs envsubst on files in /etc/nginx/templates/
COPY nginx.conf.template /etc/nginx/templates/default.conf.template
# Keep local nginx.conf as fallback (overridden by template at runtime)
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
