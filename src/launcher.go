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
	"github.com/andlabs/ui"
)

var window ui.Window

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
	Profiles []string
}

func (o *Options) ChangeProfile(profile string) {
	 if Pos(profile, o.Profiles) != -1 {
		 o.LauncherConfig["Profiles"]["currentprofile"] = profile
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
	
//	buf, _ := ioutil.ReadFile("readme.txt")
//	o.LocalVersion = re.FindStringSubmatch(string(buf))[1]

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

	o.Profiles = make([]string, 0, len(o.LauncherConfig["Profiles"]))
	for k := range o.LauncherConfig["Profiles"] {
		if k != "currentprofile" {
			o.Profiles = append(o.Profiles, k)
		}
	}

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

func StartOpenMW(o *Options) {
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

func (o *Options)GetSelectedContentFiles() []string {
	profile, _ := o.LauncherConfig.Get("Profiles", "currentprofile")
	content_files, _ := o.LauncherConfig.Get("Profiles", profile.(string))
	switch content_files.(type){
	case string:
		content_files = []string{content_files.(string)}
	}
	return content_files.([]string)
}


func main() {
	options := NewOptions()
	if options.LauncherConfig["General"]["firstrun"] == "true" {
		fmt.Println("Its a first run of OpenMW, please run official omwlauncher for setting Morrowind path and initial settings")
		os.Exit(1)
	}
	profile := "new"
	options.ChangeProfile(profile)
	fmt.Println(options.DataPath)
	content_files := options.GetSelectedContentFiles()

	fmt.Println("Starting with profile:", profile)
	for _, f := range options.GetAvailableContentFiles(){
		if Pos(f, content_files) != -1 {
			fmt.Print(" [x] ")
		} else {
			fmt.Print(" [ ] ")
		}
		fmt.Println(f)
	}
//	StartOpenMW(options)

	go ui.Do(func() {
		rv_label := ui.NewLabel("Available OpenMW: fetching...")
		play_button := ui.NewButton("Play")
//		profiles_combo := ui.New
		layout := []ui.Control {
			ui.NewLabel("Installed OpenMW: " + options.LocalVersion),
			rv_label,
		}
//		profile, _ := options.LauncherConfig.Get("Profiles", "currentprofile")
//		for _, p := range options.Profiles {
//			c := ui.NewCheckbox(p)
//			if p ==  profile.(string) {
//				c.SetChecked(true)
//			}
//			c.OnToggled((func() {
//				if !c.Checked() {
//					c.SetChecked(true)
//				}
//			}))
//			layout = append(layout, c)
//		}
		layout = append(layout, play_button)
		stack := ui.NewVerticalStack(layout...)
		window = ui.NewWindow("OpenMW Launcher", 200, 200, stack)
		window.OnClosing(func() bool {
			ui.Stop()
			return true
		})

		play_button.OnClicked(func() {
			StartOpenMW(options)
		})
		window.Show()

		go func() {
			ver := options.FetchRemoteVersion()
			if !options.IsLatest() {
				println("Maybe new OpenMW version is available! Go to", constants.SiteUrl)
				rv_label.SetText(fmt.Sprintf("Available OpenMW: %v", ver ))
			} else {
				rv_label.SetText("You have latest version.")
			}
		}()
	})
	err := ui.Go()
	if err != nil {
		panic(err)
	}
}
