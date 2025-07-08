# Stage 1: Build the React application
FROM node:18-alpine AS builder

WORKDIR /app

# Copy package.json and package-lock.json
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy the rest of the application source code
COPY . .

# Set the API URL for the backend (can be overridden at build time)
ARG REACT_APP_API_URL_CELL_A=/api/cell-a
ARG REACT_APP_API_URL_CELL_B=/api/cell-b
ENV REACT_APP_API_URL_CELL_A=${REACT_APP_API_URL_CELL_A}
ENV REACT_APP_API_URL_CELL_B=${REACT_APP_API_URL_CELL_B}

# Build the static files
RUN npm run build

# Stage 2: Serve the static files with Nginx
FROM nginx:1.25-alpine

# Copy the built static files from the builder stage
COPY --from=builder /app/build /usr/share/nginx/html

# Copy a custom nginx config to handle client-side routing and proxying
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Expose port 80
EXPOSE 80

# Start Nginx
CMD ["nginx", "-g", "daemon off;"]