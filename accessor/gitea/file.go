package gitea

import (
	"fmt"
	"io"
	"log"
	"os"
)

func (accessor *Accessor) copyFile(externalFilePath string, giteaPath string) {
	_, err := os.Stat(externalFilePath)
	if os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Warning: cannot copy non-existant attachment file: \"%s\"\n", externalFilePath)
		return
	}

	in, err := os.Open(externalFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	out, err := os.Create(giteaPath)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		log.Fatal(err)
	}

	err = out.Close()
	if err != nil {
		log.Fatal(err)
	}
}
