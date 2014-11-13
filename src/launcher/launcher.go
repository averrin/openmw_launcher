package main

import (
	"fmt"
	"os"
	"log"
	"gopkg.in/qml.v1"
//	"time"
//	"github.com/jmcvetta/randutil"
)

func run() error {
	options := NewOptions()
	engine := qml.NewEngine()
	log.Println(options.Profiles)

	controls, err := engine.LoadFile("src/main.qml")
	if err != nil {
		return err
	}

	context := engine.Context()
	context.SetVars(options)
//	context.SetVar("ProfilesModel", options.Profiles)
//	context.SetVar("ContentModel", options.ContentFiles)
	fmt.Println(options.Profiles.Current, options.Profiles.List)
	ci := Pos(options.Profiles.Current, options.Profiles.List)
	context.SetVar("CurrentProfile", ci)
	window := controls.CreateWindow(nil)

	window.Show()
	
	go func(){
		options.FetchRemoteVersion()
		window.ObjectByName("Rlabel").Set("text", "Available version: " + options.RemoteVersion)
	}()
	
//	go func(){
//		for i := 0; i<100; i++ {
//			time.Sleep(1000 * time.Millisecond)
//			p, _ := randutil.ChoiceInt([]int{0,1,2})
//			println(p)
//			options.Profiles.Select(p)
//		}
//	}()
	
	window.Wait()
	return nil
}


func main() {
	options := NewOptions()
	if options.LauncherConfig["General"]["firstrun"] == "true" {
		fmt.Println("Its a first run of OpenMW, please run official omwlauncher for setting Morrowind path and initial settings")
		os.Exit(1)
	}

	for _, f := range options.GetAvailableContentFiles() {
		if Pos(f, options.ContentFiles.List) != -1 {
			fmt.Print(" [x] ")
		} else {
			fmt.Print(" [ ] ")
		}
		fmt.Println(f)
	}

	if err := qml.Run(run); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
