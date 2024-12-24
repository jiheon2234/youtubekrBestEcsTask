# 빌드 스테이지
FROM golang:1.23-alpine AS build

WORKDIR /app
# go.mod와 go.sum 복사 후 의존성 다운로드 및 빌드
COPY go.mod go.sum ./
RUN go mod download

# 소스 코드 전체 복사 및 빌드
COPY . .
RUN go build -o crawler .

#실제 도커파일
FROM alpine:latest
WORKDIR /app

ENV GOROUTINE_CNT=2

RUN apk add --no-cache tzdata
ENV TZ=Asia/Seoul

# 빌드된 바이너리 복사
COPY --from=build /app/crawler /app/crawler

#ENTRYPOINT  ["/app/crawler"]
CMD ["/app/crawler"]
