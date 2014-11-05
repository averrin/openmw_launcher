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
	"constants"
)

type Options struct {
	LocalVersion string
	RemoteVersion string
	DataPath string
	CWD string
}

func (o *Options) IsLatest() bool {
	return o.LocalVersion == o.RemoteVersion
}

func NewOptions() (o *Options) {
	o = new(Options)
	re := regexp.MustCompile(`Version: ([\d\.]+)`)
	o.CWD, _ = os.Getwd()
	buf, _ := ioutil.ReadFile("readme.txt")
	o.LocalVersion = re.FindStringSubmatch(string(buf))[1]

	resp, _ := http.Get(constants.RemoteReadmeUrl)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	o.RemoteVersion = re.FindStringSubmatch(string(body))[1]

	usr, _ := user.Current()
	settings_folder := path.Join(usr.HomeDir, constants.OpenMWSettingsDir )
	settings_file := path.Join(settings_folder, "openmw.cfg")
	buf, _ = ioutil.ReadFile(settings_file)
	data_re := regexp.MustCompile(`data="(.+)"`)
	o.DataPath = data_re.FindStringSubmatch(string(buf))[1]

	return o
}

func StartOpenMW() {
	cmd := exec.Command(constants.OpenMWExec)
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Waiting for command to finish...")
	err = cmd.Wait()
	log.Printf("Command finished with error: %v", err)
}

func main() {

	options := NewOptions()
	fmt.Println(options)
	fmt.Println(options.IsLatest())

//	StartOpenMW()

}
