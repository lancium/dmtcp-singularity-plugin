package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/sylabs/singularity/pkg/sylog"
	"github.com/sylabs/singularity/pkg/cmdline"
	"github.com/sylabs/singularity/pkg/runtime/engine/config"
	singularity "github.com/sylabs/singularity/pkg/runtime/engine/singularity/config"
	pluginapi "github.com/sylabs/singularity/pkg/plugin"
	clicallback "github.com/sylabs/singularity/pkg/plugin/callback/cli"
)

// Plugin is the only variable which a plugin MUST export.
// This symbol is accessed by the plugin framework to initialize the plugin.
var Plugin = pluginapi.Plugin{
	Manifest: pluginapi.Manifest{
		Name:        "lancium.com/dmtcp-singularity",
		Author:      "Lancium",
		Version:     "0.1.1",
		Description: "This is a plugin to add checkpointing to Singularity with DMTCP",
	},
	
	Callbacks: []pluginapi.Callback{
		(clicallback.Command)(callbackPluginCmd),
		(clicallback.SingularityEngineConfig)(callbackDMTCP),
	},
}

var(
	BindPaths []string
)
var isCheckpoint = false 

func callbackPluginCmd(manager *cmdline.CommandManager) {
	// create command: singularity checkpoint
	var checkpointCmd = &cobra.Command{
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("invalid command")
		},
		DisableFlagsInUseLine: false,
		Use:           "checkpoint",
		Short:         "Manage instances with checkpoint support",
		Long:          "Manage instances with checkpoint support",
		Example:       "...",
		SilenceErrors: true,
	}

	// register singularity checkpoint command
	manager.RegisterCmd(checkpointCmd)

	//New scheme:
	//singularity checkpoint [options for checkpointing] [instance] [command to run]
	//this will translate roughly to singulairty exec [instance] /.dmtcp/bin/dmtcp_launch 
	//[options] [command to run]

	// get reference to start Run method
	instanceStartCmd := manager.GetCmd("instance_start")
	if instanceStartCmd == nil {
		sylog.Warningf("Could not find instance start command")
		return
	}
	instanceStartCmdRun := instanceStartCmd.Run

	// create command: singularity checkpoint start
	var checkpointStartCmd = &cobra.Command{
		DisableFlagsInUseLine: true,
		Args:                  cobra.MinimumNArgs(3),
		Use:                   "start [args ...]",
		Short:                 "Start an instance",
		Long:                  "Start an instance with checkpoint capabilities",
		Example:               "singularity checkpoint start <image> <name> <checkpoint directory>",
		Run: func(cmd *cobra.Command, args []string) {
			isCheckpoint = true
			//init checkpoint
			//TODO: start a coordinator
			instanceStartCmdRun(instanceStartCmd, args)
		},
		TraverseChildren: true,
	}
	// register checkpoint start command
	manager.RegisterSubCmd(checkpointCmd, checkpointStartCmd)
	checkpointStartCmd.Flags().AddFlagSet(instanceStartCmd.Flags())

	// get reference to exec Run method
	execCmd := manager.GetCmd("exec")
	if execCmd == nil {
		sylog.Warningf("Could not find exec command")
		return
	}
	execCmdRun := execCmd.Run


	// get reference to stop Run method
	instanceStopCmd := manager.GetCmd("instance_stop")
	if instanceStopCmd == nil {
		sylog.Warningf("Could not find instance stop command")
		return
	}
	instanceStopCmdRun := instanceStopCmd.RunE


	var checkpointStopCmd = &cobra.Command{
		DisableFlagsInUseLine: true,
		Args:                  cobra.MinimumNArgs(1),
		Use:                   "stop [args ...]",
		Short:                 "Checkpoint, stop an instance",
		Long:                  "Stop an instance with checkpoint capabilities after checkpointing",
		Example:               "singularity checkpoint stop <name>",
		Run: func(cmd *cobra.Command, args []string) {
			isCheckpoint = true
			//TODO: send a checkpoint command
			//Checkpoint
			newArgs := []string{"instance://"+args[0],"sh", "/.dmtcp/scripts/checkpoint_then_kill.sh"}
			fmt.Println(newArgs)
			ctkCmd := exec.Command("singularity", "exec", "instance://"+args[0], "sh", "/.dmtcp/scripts/checkpoint_then_kill.sh")
			ctkCmd.Start()
			ctkCmd.Wait()
			fmt.Println("Shutting down instance...")
			instanceStopCmdRun(instanceStopCmd, args)
		},
		TraverseChildren: true,
	}
	// register checkpoint stop command
	manager.RegisterSubCmd(checkpointCmd, checkpointStopCmd)
	checkpointStopCmd.Flags().AddFlagSet(instanceStopCmd.Flags())


	// create command: singularity checkpoint exec
	var checkpointExecCmd = &cobra.Command{
		DisableFlagsInUseLine: true,
		Args:                  cobra.MinimumNArgs(2),
		Use:                   "exec [args ...]",
		Short:                 "Execute a program in the instance",
		Long:                  "Execute a program in the given instance with checkpoint",
		Example:               "singularity checkpoint exec <name> <command>",
		Run: func(cmd *cobra.Command, args []string) {
			isCheckpoint = true
			newArgs := []string{args[0],"sh", "/.dmtcp/scripts/launch.sh"}
			newArgs = append(newArgs, args[1:]...)
			fmt.Println(newArgs)
			execCmdRun(execCmd, newArgs)
		},
		TraverseChildren: true,
	}
	// register checkpoint exec command
	manager.RegisterSubCmd(checkpointCmd, checkpointExecCmd)
	checkpointExecCmd.Flags().AddFlagSet(execCmd.Flags())

	// create command: singularity checkpoint job_run
	// both starts an instance and runs a command
	var checkpointJobRunCmd = &cobra.Command{
		DisableFlagsInUseLine: true,
		Args:                  cobra.MinimumNArgs(3),
		Use:                   "job_run [args ...]",
		Short:                 "Start an instance and execute a program",
		Long:                  "Create an instance with DMTCP ready, then start a program with the DMTCP wrappers",
		Example:               "singularity checkpoint job_run <image> <name> <command>",
		Run: func(cmd *cobra.Command, args []string) {
			isCheckpoint = true
			// run the start command with the container img and name
			checkpointStartCmd.Run(checkpointStartCmd, args[0:2])
			// format a slice by copying and modifying single element
			execSlice := make([]string, len(args))
			copy(execSlice, args)
			execSlice[1] = "instance://"+execSlice[1]
			// append singularity command on subset of modified slice to fit
			cmdSlice := append([]string{"checkpoint", "exec"}, execSlice[1:]...)
			// actually exec
			ctkCmd := exec.Command("singularity", cmdSlice[:]...)
			ctkCmd.Stdout = os.Stdout
			ctkCmd.Stderr = os.Stderr
			ctkCmd.Start()
			ctkCmd.Wait()
			// stop instance
			checkpointStopCmd.Run(checkpointStopCmd, args[1:2])
		},
		TraverseChildren: true,
	}
	// register checkpoint job_run command
	manager.RegisterSubCmd(checkpointCmd, checkpointJobRunCmd)
	// must register instance start and exec's, hopefully not overlap/destroying
	checkpointJobRunCmd.Flags().AddFlagSet(instanceStartCmd.Flags())
	checkpointJobRunCmd.Flags().AddFlagSet(execCmd.Flags())

	// create command: singularity checkpoint run
	var checkpointRunCmd = &cobra.Command{
		DisableFlagsInUseLine: true,
		Args:                  cobra.MinimumNArgs(2),
		Use:                   "run [args ...]",
		Short:                 "Execute the runscript in the instance",
		Long:                  "Execute the runscript in the given instance with checkpoint",
		Example:               "singularity checkpoint run <instance> <command>",
		Run: func(cmd *cobra.Command, args []string) {
			newArgs := []string{args[0],"sh", "/.dmtcp/scripts/launch.sh", "/singularity"}
			fmt.Println(newArgs)
			execCmdRun(execCmd, newArgs)
		},
		TraverseChildren: true,
	}
	// register checkpoint run command
	manager.RegisterSubCmd(checkpointCmd, checkpointRunCmd)


	// create command: singularity checkpoint running tasks
	var checkpointMakeCmd = &cobra.Command{
		DisableFlagsInUseLine: true,
		Args:                  cobra.MinimumNArgs(1),
		Use:                   "make [args ...]",
		Short:                 "Checkpoint in the instance",
		Long:                  "Checkpoint in the given instance",
		Example:               "singularity checkpoint make <instance>",
		Run: func(cmd *cobra.Command, args []string) {
			newArgs := []string{args[0],"sh", "/.dmtcp/scripts/checkpoint.sh"}
			fmt.Println(newArgs)
			execCmdRun(execCmd, newArgs)
		},
		TraverseChildren: true,
	}
	// register checkpoint make command
	manager.RegisterSubCmd(checkpointCmd, checkpointMakeCmd)

	// create command: singularity checkpoint restart task
	var checkpointRestartCmd = &cobra.Command{
		DisableFlagsInUseLine: true,
		Args:                  cobra.MinimumNArgs(1),
		Use:                   "restart [args ...]",
		Short:                 "Restart from script",
		Long:                  "Restart from checkpoint script in mounted checkpoint file",
		Example:               "singularity checkpoint restart <instance>",
		Run: func(cmd *cobra.Command, args []string) {
			newArgs := []string{args[0],"sh", "/.dmtcp/scripts/restart.sh"}
			fmt.Println(newArgs)
			execCmdRun(execCmd, newArgs)
		},
		TraverseChildren: true,
	}
	// register checkpoint make command
	manager.RegisterSubCmd(checkpointCmd, checkpointRestartCmd)

	// create command: singularity checkpoint job_restart
	// both starts an instance and restarts the checkpoint in the current directory
	var checkpointJobRestartCmd = &cobra.Command{
		DisableFlagsInUseLine: true,
		Args:                  cobra.MinimumNArgs(2),
		Use:                   "job_restart [args ...]",
		Short:                 "Start an instance and restart from checkpoint",
		Long:                  "Create an instance with DMTCP ready, then restart a program with the DMTCP wrappers",
		Example:               "singularity checkpoint job_restart <image> <name>",
		Run: func(cmd *cobra.Command, args []string) {
			isCheckpoint = true
			// run the start command with the container img and name
			checkpointStartCmd.Run(checkpointStartCmd, args[0:2])
			// format a slice by copying and modifying single element
			execSlice := make([]string, len(args))
			copy(execSlice, args)
			execSlice[1] = "instance://"+execSlice[1]
			// append singularity command on subset of modified slice to fit
			cmdSlice := append([]string{"checkpoint", "restart"}, execSlice[1])
			// actually exec
			ctkCmd := exec.Command("singularity", cmdSlice[:]...)
			ctkCmd.Stdout = os.Stdout
			ctkCmd.Stderr = os.Stderr
			ctkCmd.Start()
			ctkCmd.Wait()
			// stop instance
			checkpointStopCmd.Run(checkpointStopCmd, args[1:2])
		},
		TraverseChildren: true,
	}
	// register checkpoint job_restart command
	manager.RegisterSubCmd(checkpointCmd, checkpointJobRestartCmd)
	// must register instance start and exec's, hopefully not overlap/destroying
	checkpointJobRestartCmd.Flags().AddFlagSet(instanceStartCmd.Flags())
	checkpointJobRestartCmd.Flags().AddFlagSet(execCmd.Flags())

	// Ask for the checkpoint, wait for finish.
	// TODO: Add script for checkpointing first, then run here.
}

func callbackDMTCP(common *config.Common) {
	c, ok := common.EngineConfig.(*singularity.EngineConfig)
	if !ok {
		sylog.Warningf("Unexpected engine config")
		return
	}
	//Add bind for DMTCP if in environment.
	if isCheckpoint{
		dmtcpLocation := "/opt/dmtcp-singularity/"
		dmtcpLibLocation := "/usr/lib/dmtcp"
		if(dmtcpLocation == ""){
			sylog.Errorf("No DMTCP location found. Run install script?")
			return 
		}
		origBind := c.GetBindPath()
		
		//Build new mount path
		var dmtcpBind singularity.BindPath
		dmtcpBind.Source = dmtcpLocation
		dmtcpBind.Destination = "/.dmtcp/"

		var dmtcpLibBind singularity.BindPath
		dmtcpLibBind.Source = dmtcpLibLocation
		dmtcpLibBind.Destination = "/.singularity.d/libs/"

		//Create a directory for checkpoints
		ckptDir := "./.checkpoint"
		_, err := os.Stat(ckptDir)
		if os.IsNotExist(err){
			mkCkptDir := os.MkdirAll(ckptDir, os.ModePerm)
			if mkCkptDir != nil {
				sylog.Errorf("Could not create checkpoint directory. Will use current.")
				ckptDir = "./"
			}
		}

		var ckptBind singularity.BindPath
		ckptBind.Source = ckptDir
		ckptBind.Destination = "/.checkpoint/"
		
		//Option for read only
		var dmtcpOptions = map[string]*singularity.BindOption{
			"ro":        &singularity.BindOption{},
		}


		//Option for read/write
		var ckptOptions = map[string]*singularity.BindOption{
			"rw":        &singularity.BindOption{},
		}
		dmtcpBind.Options = dmtcpOptions
		dmtcpLibBind.Options = dmtcpOptions
		ckptBind.Options = ckptOptions
		
		
		//Set to include this new bind path
		c.SetBindPath(append(origBind, dmtcpBind, dmtcpLibBind, ckptBind))
	}
	return
}



