cmd
コマンドプロンプロのとき
https://docs.aws.amazon.com/ja_jp/lambda/latest/dg/golang-package.html

go.exe get -u github.com/aws/aws-lambda-go/cmd/build-lambda-zip


set GOOS=linux

ファイルが複数のとき
go build -o main main.go cloudwatch.go
いつもは
go build -o main main.go

%USERPROFILE%\Go\bin\build-lambda-zip.exe -output main.zip main

