# ocrwebservice.go


### Install

`go get github.com/poorny/ocrwebservice`


### Usage

Export to env `LICENSE_CODE` and `USERNAME`.

_Get LICENSE CODE and USERNAME in your [Dashboard](http://www.ocrwebservice.com/dashboard)_

```go
package main

import (
	"fmt"

  "github.com/poorny/ocrwebservice"
)

func main() {
	path := "/path/to/file.pdf"
	lang := "english" // See http://www.ocrwebservice.com/api/keyfeatures to others

	result, err := ocrws.OcrWs(path, lang)
	fmt.Println(result.Text())
}
```
