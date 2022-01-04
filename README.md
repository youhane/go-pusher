# Go Photo Pusher

Made with Go and VueJs

<br>

## How to start (if you're unfamiliar with Go)

1. Make the folder/clone this repo
2. Run go mod init example.com/try echo in your terminal
3. Run this in your terminal

``` go
go get github.com/labstack/echo
go get github.com/labstack/echo/middleware
go get github.com/mattn/go-sqlite3
go get github.com/pusher/pusher-http-go
```

4. Make sure the code is there
5. Run go run main.go

<br>

## Notes

- If there are errors on the imports in main.go, most problems are that you haven't installed the package
- Don't forget to make and use the pusher app credentials
- Refresh the page if the picture wont appear
- If you want to delete the pics, delete in the public/uploads folder per pic
- If the project won't run, try opening it in another window of the text editor, with Go-Pusher as the root folder
