# HOW TO START

### First

- make the folder/clone this
- go mod init example.com/try echo
- go run main.go

### Then Run this

go get github.com/labstack/echo
go get github.com/labstack/echo/middleware
go get github.com/mattn/go-sqlite3
go get github.com/pusher/pusher-http-go


## Put this in the main.go file

```
import (
	"database/sql"
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/mattn/go-sqlite3"
	pusher "github.com/pusher/pusher-http-go"
) 
```

If there are squiggly lines, hover on the red ones then click the get package

Make the pusher app credentials

Then copy all of the code
Understand it
Update the pusher credentials
Then youre all done
Just go run main.go

Refresh the page if the picture wont appear

If you want to delete the files, delete in the public/uploads folder per files

# Folder information

## database

Self explanatory

## public

For public files, such as html and images

### Notes

If the squigly lines wont dissapear in vscode, try opening the project/folder in a new window, that makes the folder a root

Buka foldernya di window baru, yang isinya dia doang