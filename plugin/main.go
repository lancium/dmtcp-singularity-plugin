package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/sylabs/singularity/internal/pkg/sylog"
	"github.com/sylabs/singularity/pkg/cmdline"
	"github.com/sylabs/singularity/pkg/runtime/engine/config"
	singularity "github.com/sylabs/singularity/pkg/runtime/engine/singularity/config"
	pluginapi "github.com/sylabs/singularity/pkg/plugin"
)

// Plugin is the only variable which a plugin MUST export.
// This symbol is accessed by the plugin framework to initialize the plugin.
var Plugin = pluginapi.Plugin{
	Manifest: pluginapi.Manifest{
		Name:        "lancium.com/dmtcp-singularity",
		Author:      "Lancium",
		Version:     "0.0.1",
		Description: "This is a plugin to add checkpointing to Singularity with DMTCP",
	},

	Initializer: pluginImplementation{},
}

var(
	BindPaths []string
)

type pluginImplementation struct{}

func (p pluginImplementation) Initialize(r pluginapi.Registry) error {
	r.AddCLIMutator(pluginapi.CLIMutator{
		Mutate: func(manager *cmdline.CommandManager) {
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
			//this will translate to singulairty exec [instance] /.dmtcp/bin/dmtcp_launch [options] [command to run]
			
			
			// get reference to exec Run method
			execCmd := manager.GetCmd("exec")
			if execCmd == nil {
				sylog.Warningf("Could not find exec command")
				return
			}
			execCmdRun := execCmd.Run
			
			// create command: singularity checkpoint exec
			var checkpointExecCmd = &cobra.Command{
				DisableFlagsInUseLine: true,
				Args:                  cobra.MinimumNArgs(2),
				Use:                   "exec [args ...]",
				Short:                 "Execute a program in the instance",
				Long:                  "Execute a program in the given instance with checkpoint",
				Example:               "singularity checkpoint exec <name> <command>",
				Run: func(cmd *cobra.Command, args []string) {
					newArgs := []string{args[0],"sh", "/.dmtcp/scripts/launch.sh"}
					newArgs = append(newArgs, args[1:]...)
					fmt.Println(newArgs)
					execCmdRun(execCmd, newArgs)
				},
				TraverseChildren: true,
			}
			// register checkpoint exec command
			manager.RegisterSubCmd(checkpointCmd, checkpointExecCmd)
			
			
			// get reference to run Run method -- not gonna be necessary, just using exec
			/*
			runCmd := manager.GetCmd("run")
			if runCmd == nil {
				sylog.Warningf("Could not find exec command")
				return
			}
			runCmdRun := runCmd.Run
			*/
			
			// create command: singularity checkpoint run
			var checkpointRunCmd = &cobra.Command{
				DisableFlagsInUseLine: true,
				Args:                  cobra.MinimumNArgs(1),
				Use:                   "run [args ...]",
				Short:                 "Execute the runscript in the instance",
				Long:                  "Execute the runscript in the given instance with checkpoint",
				Example:               "singularity checkpoint run <command>",
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
				Example:               "singularity checkpoint make <command>",
				Run: func(cmd *cobra.Command, args []string) {
					newArgs := []string{args[0],"sh", "/.dmtcp/scripts/checkpoint.sh"}
					fmt.Println(newArgs)
					execCmdRun(execCmd, newArgs)
				},
				TraverseChildren: true,
			}
			// register checkpoint make command
			manager.RegisterSubCmd(checkpointCmd, checkpointMakeCmd)
			
			// Ask for the checkpoint, wait for finish.
			// TODO: Add script for checkpointing first, then run here.
			
		},
	})
	r.AddEngineConfigMutator(pluginapi.EngineConfigMutator{
		Mutate: func(common *config.Common) {
			c, ok := common.EngineConfig.(*singularity.EngineConfig)
			if !ok {
				sylog.Warningf("Unexpected engine config")
				return
			}
			//Add bind for DMTCP if in environment.
			dmtcpLocation := os.Getenv("SINGULARITY_DMTCP")
			if(dmtcpLocation == ""){
				sylog.Errorf("No DMTCP location found. Run install script?")
				return 
			}
			origBind := c.GetBindPath()
			c.SetBindPath(append(origBind, dmtcpLocation+":/.dmtcp/"))
			//fmt.Println(c.GetBindPath())
		},
	})

	return nil
}


