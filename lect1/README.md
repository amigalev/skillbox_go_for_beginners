# First lecture
In this folder we have an example from the first lecture + homework in the same project.

# Main tips
Terminal:
1. go mod init skillbox - creates new project
2. go mod tidy - resolves dependencies (deps inside go.mod file) or removes unnecessarry deps (if you add deps in go.mod, but don't use it in prject)
3. go mod download - download and resolve deps from the go.mod, adds the version instead of "latest"
4. go run main.go - compiles and executes main.go 
