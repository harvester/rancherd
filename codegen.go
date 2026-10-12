//go:generate /bin/rm -rf pkg/generated
//go:generate go run ./pkg/codegen/main.go

package main

import (
	"fmt"
)

func main() {
	fmt.Println("Run go generate")
}
