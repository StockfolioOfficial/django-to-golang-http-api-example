# django-to-golang-rest-api-example
Golang RESTful API example like django project

# Reference
https://github.com/bxcodec/go-clean-arch

# Libs
- [Echo Framework](https://echo.labstack.com/) - HTTP Router
- [GORM](https://gorm.io/) - ORM
- [logrus](https://github.com/sirupsen/logrus) - LOG Util
- [validator](https://github.com/go-playground/validator) - data validate util

# Database
[SQLite](https://www.sqlite.org/)

# 시작방법
1. Golang 1.7이상 설치
2. `go mod download` 터미널에서 실행 (프로젝트 루트에서)  
`pip install -r requirements.txt`나 `npm install`이랑 같은 맥락
```bash
# go mod download all
```
3. `go run .` 터미널에서 실행 (프로젝트 루트에서)
```bash
# go run .
```

# 빌드 후 실행 방법
1. [#시작방법](#시작방법)-1 동일
2. [#시작방법](#시작방법)-2 동일
3. `go build -o app .` 터미널에서 실행 (프로젝트 루트에서)  
**윈도우의 경우 `app` 이 아니라 `app.exe` 로 해주세요**
```bash
# go build -o app .
```
4. `./app` 터미널에서 실행 (프로젝트 루트에서)
```bash
# ./app
```

# API
## Article Fetch (List)
목록 가져오기
```bash
# curl -X GET http://localhost:8000/articles
```

## Article Create
생성
```bash
# curl -X POST -H "Content-Type: application/json" -d "{\"title\":\"제목에 대한 데이터\",\"content\":\"컨텐츠에 대한 데이터\",\"author\":{\"id\":1}}" http://localhost:8000/articles
```

## Article Get (One)
읽기
```bash
# curl -X GET http://localhost:8000/articles/{id}
```

## Article Delete
삭제 
```bash
# curl -X DELETE http://localhost:8000/articles/{id}
```
