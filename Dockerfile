FROM golang:1.23-alpine

# 작업 디렉토리 설정
WORKDIR /

# 의존성 파일 복사 및 다운로드
COPY go.mod go.sum ./
RUN go mod download

# 소스 코드 복사
COPY . .

# air를 사용하여 애플리케이션 실행
CMD ["go", "run", "./cmd/main.go"]