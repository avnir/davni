package main

import (
	"fmt"
	"log"
	"os/exec"
)

func main() {
	cmd := exec.Command("sh", "-c", "docker run -v /tmp/pdfs:/pdfs -w /pdfs madnight/docker-alpine-wkhtmltopdf google.com test.pdf")
	_, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("docker run failed")
		log.Fatal(err)
	}

	cmd = exec.Command("sh", "-c", "ls /pdfs")
	out, outErr := cmd.CombinedOutput()
	if outErr != nil {
		log.Println("ls failed")
		log.Fatal(outErr)
	}

	fmt.Println(string(out))
}
