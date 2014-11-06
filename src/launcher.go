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
	"github.com/averrin/go-ini"
	"path/filepath"
	"strings"
)

type Options struct {
	LocalVersion string
	RemoteVersion string
	DataPath string
	CWD string
	LauncherConfig ini.File
	OMWConfig ini.File	
}

func (o *Options) IsLatest() bool {
	return o.LocalVersion == o.RemoteVersion
}

func NewOptions() (o *Options) {
	o = new(Options)
	re := regexp.MustCompile(`OpenMW version ([\d\.]+)`)
	o.CWD, _ = os.Getwd()
	
//	buf, _ := ioutil.ReadFile("readme.txt")
//	o.LocalVersion = re.FindStringSubmatch(string(buf))[1]

	version, _ := exec.Command(constants.OpenMWExec, "--version").Output()
	o.LocalVersion = re.FindStringSubmatch((string)(version))[1]
	

	resp, _ := http.Get(constants.RemoteReadmeUrl)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	re = regexp.MustCompile(`Version: ([\d\.]+)`)
	o.RemoteVersion = re.FindStringSubmatch(string(body))[1]

	usr, _ := user.Current()
	settings_folder := path.Join(usr.HomeDir, constants.OpenMWSettingsDir)

	o.LauncherConfig, _ = ini.LoadFile(path.Join(settings_folder, "launcher.cfg"))
	o.OMWConfig, _ = ini.LoadFile(path.Join(settings_folder, "openmw.cfg"))
	d, _ := o.OMWConfig.Get("", "data")
	o.DataPath = strings.Trim(d.(string), `"`)

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

func Pos(value interface {}, slice []string) int {
	for p, v := range slice {
		if (v == value) {
			return p
		}
	}
	return -1
}

func main() {

	options := NewOptions()
//	fmt.Println(options)
	profile, _ := options.LauncherConfig.Get("Profiles", "currentprofile")
	content_files, _ := options.LauncherConfig.Get("Profiles", profile.(string))
	switch content_files.(type){
	case string:
		content_files = []string{content_files.(string)}
	}
	fmt.Println(options.DataPath)
	fmt.Printf("Content to load (%v):\n", profile)
//	fmt.Println(options.IsLatest())

	p := path.Join(options.DataPath, "/*.esm")
	files, _ := filepath.Glob(p)
	available_content := make([]string, 0)
	for _, f := range files {
		_, c := path.Split(f)
		available_content = append(available_content, c)
	}

	for _, f := range available_content{
		if Pos(f, content_files.([]string)) != -1 {
			fmt.Print(" [x] ")
		} else {
			fmt.Print(" [ ] ")
		}
		fmt.Println(f)
	}

//	StartOpenMW()

}
