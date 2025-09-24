# Тестовое задание компании Тагес

> Хотел делать тестовое на виртуалке с убунтой, но там не "работает" GoLand  
> к тому же, стоит старая версия голанга (и дебагера dlv соответственно)  
> Поэтому у проекта нет Makefile'а и запускать сервис придётся руками вводя команды  
> А самое печальное то, что репа написана только под виндузню. Следовательно, сервак не запустится
> ни на макпуке, ни на линухе.

### Требование для запуска
- Go 1.25.0
- ОС Windows

### Запуск сервиса
На всякий случай, возможно, потребуется пересобрать `proto` файл из папки `api`:
```shell
$ protoc -I .\api\protos\ --go_out=. --go-grpc_out=. .\api\protos\file.proto
```
Транслятор protobuf на винде можно поставить так:
```shell
$ go install google.golang.org/protobuf/cmd/protoc-gen-go
$ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
```
Затем скачать [архив](https://github.com/protocolbuffers/protobuf/releases) для винды, распаковать и добавить путь 
до бинаря protoc в переменную `PATH`.

Ставим все либы или вендорим проект:
```shell
$ go mod download
```
или
```shell
$ go mod vendor
```
И запускаем сервер:
```shell
$ go build
```