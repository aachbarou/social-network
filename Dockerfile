# Stage 1: Build the React client
FROM node:18-alpine AS client-build

WORKDIR /app/client

# Copy client package.json and install dependencies
COPY client/package.json ./
COPY client/yarn.lock ./
RUN yarn install

# Copy the rest of the client source code
COPY client/ ./

# Build the React app for production
RUN yarn build


# Stage 2: Build and run the Node.js server
FROM node:18-alpine AS server-build

WORKDIR /app

# Copy server package.json and install dependencies
COPY server/package.json ./
COPY server/yarn.lock ./
RUN yarn install

# Copy the rest of the server source code
COPY server/ ./

# --- IMPORTANT ---
# Copy the built client files from the previous stage
COPY --from=client-build /app/client/build ./client/build

# Expose the port your server is running on
EXPOSE 5000

# Command to start the server
CMD ["node", "server.js"]
