package gitea

import (
	"io"
	"log"
	"os"
)

func (accessor *Accessor) copyFile(externalFilePath string, giteaPath string) {
	_, err := os.Stat(externalFilePath)
	if os.IsNotExist(err) {
		log.Printf("Warning: cannot copy non-existant file: \"%s\"\n", externalFilePath)
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
