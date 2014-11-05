package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"regexp"
	"net/http"
	"os/user"
	"log"
	"path"
	"os/exec"
)
func main() {
	re := regexp.MustCompile(`Version: ([\d\.]+)`)
	cwd, _ := os.Getwd()
	fmt.Println(cwd)
	buf, _ := ioutil.ReadFile("readme.txt")
	local_version := re.FindStringSubmatch(string(buf))[1]
	fmt.Println(local_version)
	resp, _ := http.Get("https://raw.githubusercontent.com/OpenMW/openmw/master/readme.txt")
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	remote_version := re.FindStringSubmatch(string(body))[1]
	fmt.Println(remote_version)

	usr, err := user.Current()
	if err != nil {
		log.Fatal( err )
	}
	settings_folder := path.Join(usr.HomeDir, "Documents/My Games/OpenMW" )
	settings_file := path.Join(settings_folder, "openmw.cfg")
	buf, _ = ioutil.ReadFile(settings_file)
	data_re := regexp.MustCompile(`data="(.+)"`)
	data_path := data_re.FindStringSubmatch(string(buf))[1]
	fmt.Println(data_path)

	cmd := exec.Command("openmw.exe")
	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Waiting for command to finish...")
	err = cmd.Wait()
	log.Printf("Command finished with error: %v", err)

}
