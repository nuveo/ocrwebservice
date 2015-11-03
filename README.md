# ocrwebservice.go

```go
package main

import (
	"fmt"

  "github.com/poorny/ocrwebservice.go"
)

func main() {
	path := "/path/to/file.pdf"
	lang := "English"

	result, err := ocrws.OcrWs(path, lang)
	fmt.Println(result.Text())
}
```
