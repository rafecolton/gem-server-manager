package gsm

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
	"syscall"
)

type BundleExecer struct {
	Configuration
	sync.Mutex
}

func NewBundleExecer(config Configuration) *BundleExecer {
	return &BundleExecer{
		Configuration: config,
	}
}

func (me *BundleExecer) ProcessInstructions(instructions *Instructions) {
	var err_w error

	me.Lock()
	stat, err := os.Stat(me.ScriptLoc)
	if err == nil {
		if stat.Size() <= 0 || (stat.Mode()&0111) != 0 {
			err_w = me.writeScript()
		}
	} else if os.IsNotExist(err) {
		err_w = me.writeScript()
	} else {
		me.Logger.Printf("os - Error detecting script file %s\n", me.ScriptLoc)
		os.Exit(13)
	}

	if err_w != nil {
		me.Logger.Print("os - Error writing script %s\n", me.ScriptLoc)
		os.Exit(17)
	}
	me.Unlock()

	cmd := exec.Command(me.ScriptLoc)
	cmd.Env = append(cmd.Env,
		fmt.Sprintf("GEMINABOX_HOST=%s", me.GibHost),
		fmt.Sprintf("GEMDIR=%s", me.GemDir),
		fmt.Sprintf("OWNER=%s", instructions.RepoOwner),
		fmt.Sprintf("REPO=%s", instructions.RepoName),
		fmt.Sprintf("REV=%s", instructions.Rev),
		fmt.Sprintf("AUTH_TOKEN=%s", instructions.AuthToken))
	err = cmd.Run()

	if err != nil {
		if msg, ok := err.(*exec.ExitError); ok {
			me.Logger.Printf("os - Error running retrieve-gems, exited %d\n",
				msg.Sys().(syscall.WaitStatus).ExitStatus())
		} else {
			me.Logger.Printf("os - Error running retrieve-gems script: %+v\n", err)
		}
	}
}

func (me *BundleExecer) writeScript() error {
	file, err := os.OpenFile(me.ScriptLoc, os.O_CREATE|os.O_WRONLY, 0755)
	defer file.Close()
	if err != nil {
		me.Logger.Println("os - Error opening script file for writing")
		return err
	}

	_, err = file.WriteString(SCRIPT_STRING)
	if err != nil {
		me.Logger.Println("os - Error writing to script file")
		return err
	}

	return nil
}
