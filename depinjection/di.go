package depinjection

import (
	"fmt"
	"io"
	"os"
)

func Greet(writer io.Writer, name string) error {
	fmt.Fprintf(writer, "Hello, %s", name)
	return nil
}

func main() {
	Greet(os.Stdout, "Elodie")
}
