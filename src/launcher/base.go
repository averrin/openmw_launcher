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
	"gopkg.in/qml.v1"
)

func Pos(value interface {}, slice []string) int {
	for p, v := range slice {
		if (v == value) {
			return p
		}
	}
	return -1
}

type Options struct {
	LocalVersion string
	RemoteVersion string
	DataPath string
	CWD string
	LauncherConfigPath string
	LauncherConfig ini.File
	OMWConfigPath string
	OMWConfig ini.File
	Profiles *Profiles
	ContentFiles *ContentFiles
}

type Profiles struct {
	Options *Options
	List []string
	Current string
}


func (profiles *Profiles) Add(p string) {
	profiles.List = append(profiles.List, p)
	qml.Changed(profiles, profiles.Length())
}

func (profiles *Profiles) Length() int {
	return len(profiles.List)
}

func (profiles *Profiles) At(index int) string {
	return profiles.List[index]
}

func (profiles *Profiles) Select(index int) {
	if index != -1 {
		p := profiles.List[index]
		profiles.Options.ChangeProfile(p)
		profiles.Options.ContentFiles.Update()
		fmt.Println(profiles.Options.ContentFiles.List, profiles.Options.ContentFiles.Length)
		qml.Changed(profiles.Options.ContentFiles, &profiles.Options.ContentFiles.Length)
	}
}

type ContentFiles struct {
	Options *Options
	List []string
	Length  int
}

func (content *ContentFiles) Add(c string) {
	content.List = append(content.List, c)
	content.Length = len(content.List)
	qml.Changed(content, &content.Length)
}

func (content *ContentFiles) Update() {
	content.Clear()
	o := content.Options
	profile, _ := o.LauncherConfig.Get("Profiles", "currentprofile")
	content_files, _ := o.LauncherConfig.Get("Profiles", profile.(string))
	switch content_files.(type){
	case string:
		content_files = []string{content_files.(string)}
	}
	content.List = content_files.([]string)
	content.Length = len(content_files.([]string))
}

func (content *ContentFiles) Clear() {
	content.List = make([]string, 0)
	content.Length = 0
}

func (content *ContentFiles) Text(index int) string {
	return content.List[index]
}

func (o *Options) ChangeProfile(profile string) {
	if Pos(profile, o.Profiles.List) != -1 {
		println("Change profile to", profile)
		o.LauncherConfig["Profiles"]["currentprofile"] = profile
		o.Profiles.Current = profile
		o.LauncherConfig.SaveFile(o.LauncherConfigPath)

		content_files, _ := o.LauncherConfig.Get("Profiles", profile)
		switch content_files.(type){
		case string:
			content_files = []string{content_files.(string)}
		}
		o.OMWConfig[""]["content"] = content_files
		o.OMWConfig.SaveFile(o.OMWConfigPath)
	}
}

func (o *Options) IsLatest() bool {
	return o.LocalVersion == o.RemoteVersion
}

func NewOptions() (o *Options) {
	o = new(Options)
	re := regexp.MustCompile(`OpenMW version ([\d\.]+)`)
	o.CWD, _ = os.Getwd()

	version, _ := exec.Command(constants.OpenMWExec, "--version").Output()
	o.LocalVersion = re.FindStringSubmatch((string)(version))[1]

	usr, _ := user.Current()
	settings_folder := path.Join(usr.HomeDir, constants.OpenMWSettingsDir)

	o.LauncherConfigPath = path.Join(settings_folder, "launcher.cfg")
	o.LauncherConfig, _ = ini.LoadFile(o.LauncherConfigPath)
	o.OMWConfigPath = path.Join(settings_folder, "openmw.cfg")
	o.OMWConfig, _ = ini.LoadFile(o.OMWConfigPath)
	d, _ := o.OMWConfig.Get("", "data")
	o.DataPath = strings.Trim(d.(string), `"`)

	o.Profiles = o.GetProfilesList()
	o.ContentFiles = o.GetSelectedContentFiles()

	return o
}

func (o *Options)FetchRemoteVersion() string{
	resp, _ := http.Get(constants.RemoteReadmeUrl)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	re := regexp.MustCompile(`Version: ([\d\.]+)`)
	o.RemoteVersion = re.FindStringSubmatch(string(body))[1]
	return o.RemoteVersion
}

func (o *Options)ImportMWINI() {
	arguments := make([]string, 0)
	arguments = append(arguments, "--game-files");
	arguments = append(arguments, "--encoding");
	arguments = append(arguments, o.OMWConfig.GetWithDefault("", "encoding", "win1251").(string));
	arguments = append(arguments, "--ini");
	arguments = append(arguments, path.Join(o.DataPath, "Morrowind.ini"));
	arguments = append(arguments, "--cfg");
	arguments = append(arguments, o.OMWConfigPath);
	fmt.Println(arguments)
	cmd := exec.Command(constants.OpenMWINIImport, arguments...)
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	err = cmd.Wait()
}

func (o *Options)StartOpenMW() {
	arguments := make([]string, 0)

	arguments = append(arguments, "--encoding");
	arguments = append(arguments, o.OMWConfig.GetWithDefault("", "encoding", "win1251").(string));
	arguments = append(arguments, "--skip-menu=1");
	arguments = append(arguments, "--new-game=1");

	cmd := exec.Command(constants.OpenMWExec, arguments...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func (o *Options)GetProfilesList() *Profiles {
	p := new(Profiles)
	p.Options = o
	p.Current = o.LauncherConfig["Profiles"]["currentprofile"].(string)
	p.List = make([]string, 0, len(o.LauncherConfig["Profiles"]))
	for k := range o.LauncherConfig["Profiles"] {
		if k != "currentprofile" {
			p.List = append(p.List, k)
		}
	}
	return p
}

func (o *Options)GetAvailableContentFiles() []string {
	exts := []string{".esm", ".esp", ".omwgame", ".omwaddon"}
	p := path.Join(o.DataPath, "*.*")
	files, _ := filepath.Glob(p)
	available_content := make([]string, 0)
	for _, f := range files {
		_, c := path.Split(f)
		if Pos(path.Ext(c), exts) != -1 {
			available_content = append(available_content, c)
		}
	}
	return available_content
}

func (o *Options)GetSelectedContentFiles() *ContentFiles {
	c := new(ContentFiles)
	c.Options = o
	c.Update()
	return c
}
