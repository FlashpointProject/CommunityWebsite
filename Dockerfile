FROM node:18 as build-stage

# Build clientside first
WORKDIR /app/web
COPY ./web/package*.json ./
RUN npm install
WORKDIR /app
COPY . .
WORKDIR /app/web
RUN npm run build

# Build serverside
FROM golang:1.20

WORKDIR /app
COPY --from=build-stage /app /app
RUN go build -o /app/fpcomm-app ./main/main.go

EXPOSE 8080
CMD ["./fpcomm-app"]