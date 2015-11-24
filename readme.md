## Build the Project

### Mac & Linux

##### Requeriments

- Go
- Git
- GCC
- [Gigo](https://github.com/LyricalSecurity/gigo)

##### How to Install

###### Install Gigo

```bash

## Clone Gigo's Repo
$ git clone https://github.com/LyricalSecurity/gigo.git && cd gigo

## Set Go's Path
$ export GOPATH=`pwd`

## Get Gigo's Dependencies
$ go get github.com/LyricalSecurity/gigo/actions

## Build Gigo
$ go build -o dist/gigo main.go

## Move Gigo to Go's Bin Directory
$ mv dist/gigo /usr/local/go/bin/gigo

## Add Go's Bin directory to the PATH (If not added already)


```

###### Compile or Run the App

```bash

## To run the app
go run basic.go

## To Compile the App (The generated binary will run without the need of having installed go, gcc or git)
go build -o dist/basic basic.go

## the last command generates the binary file in the `dist` folder

```

> Important: To open the Binary, you need to be in the root of the project (or have the shaders folder in the same folder)
> - If the shaders folder is not in the same path the program will break when running.
>
> ```bash
> 
> ## from the project's root run, to open the app
> $ dist/basic
>
> ```

> If the terminal is closed next time is opened the `GOPATH` variable needs to be set again with `export GOPATH=`pwd`/go_modules`


### Windows

#### Requirements

- GCC is required to compile the OpenGL biddings and the GLFW Wrapper.
- Git is required to fetch the dependencies.
- Go is required to fetch dependencies and compile.

#### How to compile

**Step 1**

- Install `Go` from [https://golang.org/dl/](https://golang.org/dl/)
- Add the folder `C:\Go\bin` to the PATH (If not added by the installer)

**Step 2**

- Install Git from [https://git-scm.com/download/win](https://git-scm.com/download/win)
- Add Git to the PATH

**Step 3**

Install TDM-GCC from [http://tdm-gcc.tdragon.net/](http://tdm-gcc.tdragon.net/) (Choose the x64 and x86 version, for the rest use the default settings)

> TDM-GCC is the recommended GCC compiler for Go in windows (is the only one that is officially supported by the Go Team)

**Step 4**

- Clone the project's repo (if already cloned skip to step 4)

**Step 5**

- Open the MinGW terminal installed by TDM-GCC (or open the windows CMD and run the `C:\TDM-GCC-64\mingwvars.bat` script to add the GCC to the path in that session)
- Check that the `go` command exists by doing `go version`
- Check that the `git` command exists by doing `git --version`
- Check that the `gcc` command exists by doing `gcc -v`
- If both commands work, continue to the step 5

**Step 6**

- Go to to the path where the project is located with `cd c:/<path-to-project>`

**Step 7**

- Set the `GOPATH` variable to be the same as the project path with /go_modules at the end 

```bash

set GOPATH=<path-to-the-project>/go_modules

```

**Step 8**

- Fetch the dependencies:

```bash

## Gets the OpenGL Versions Used by this project
go get github.com/go-gl/gl/all-core/gl

## Gets the GLFW Window Wrapper
go get github.com/go-gl/glfw/v3.1/glfw

## Gets the MGL32 Math Library
go get github.com/go-gl/mathgl/mgl32

## This library is used to find the path where the executable is when running a compiled version
go get github.com/kardianos/osext

```

- If all the above steps worked without errors continue to step 9

**Step 9**

- All the dependencies are now installed, and the project can be compiled.

```bash

## To run the app
go run basic.go

## To Compile the App (The generated .exe will run without the need of having installed go, gcc or git)
go build -o dist/basic.exe basic.go

## the last command generates the `.exe` file in the `dist` folder

```

> Important: To open the Binary, you need to be in the root of the project (or have the shaders folder in the same folder)
> - If the shaders folder is not in the same path the program will break when running.
>
> ```bash
> 
> ## Copy the shaders folder to dist (select `D` when asked)
> xcopy shaders "dist/shaders" 
> 
> ## Move to the `dist` folder
> cd dist
> 
> ## Run the App
> basic.exe
>
> ```

> If the terminal is closed next time is opened the `GOPATH` variable needs to be set again with `set GOPATH=<path-to-the-project>/go_modules`

