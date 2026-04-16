rm ./bin/*

FILE_PATH="bin/fileserver.exe"
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -trimpath -o $FILE_PATH"temp" .
upx --lzma -o $FILE_PATH $FILE_PATH"temp" 
rm $FILE_PATH"temp"

FILE_PATH="bin/fileserver_linux"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -o $FILE_PATH"temp" .
upx --lzma -o $FILE_PATH $FILE_PATH"temp" 
rm $FILE_PATH"temp"

FILE_PATH="bin/fileserver_mac"
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -trimpath -o $FILE_PATH .