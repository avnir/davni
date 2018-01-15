package main

import (
	"log"
	"os/exec"
)

func main() {
	cmd := exec.Command("sh", "-c", "docker run -ti -v /tmp/pdfs:/pdfs -w /pdfs madnight/docker-alpine-wkhtmltopdf google.com test.pdf")
	_, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("docker run failed")
		log.Fatal(err)
	}

	cmd = exec.Command("sh", "-c", "ls /pdfs")
	_, err = cmd.CombinedOutput()
	if err != nil {
		log.Println("ls failed")
		log.Fatal(err)
	}
}
