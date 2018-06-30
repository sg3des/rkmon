
NAME=rkmon

build: templates bindata
	go build -o ${NAME} ./

templates:
	gotemplator ./

bindata:
	go-bindata ${DEBUG} -ignore='\.scss' -pkg=main -o=rkn-assets.go -nomemcopy assets