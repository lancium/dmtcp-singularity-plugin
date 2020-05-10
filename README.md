# DMTCP for Singularity
### In a world ravaged by compute tasks that just can't stop...
This plugin is a simple solution that integrates the long-running [DMTCP](github.com/dmtcp/dmtcp) project into Singularity, a containerization platform optimized for security and performance of HPC tasks.

This is just a simple plugin that may be installed in Singularity versions < 3.5 (will need a little reworking for newer versions). A description of how to work with plugins is discussed in the plugin's directory, but it is designed to be installed with the `install.sh` script, provided a Singularity build directory for generating the plugin files as an arg. A simple set of commands are then available:

1. Start a Singularity with `singularity checkpoint start {container} {name}`
   - This is designed to mirror instance starting.
2. Start the desired program with `singularity checkpoint exec {name} {program and options}` or `singularity checkpoint run {name} {run options}`.
   - Again, attempting to mirror the included commands.
3. To run a checkpoint, `singularity checkpoint make {name}`.
4. To run a checkpoint then stop the program and container instance, run `singularity checkpoint stop {name}`.

More information is available on the wiki.