# Warp Installation Guide

## Building for Your System
If you are using any system, you can compile Warp using the following command:
```
go build -o warp main.go
```
Create an alias for your executable with the following command:
```
alias warp='~/path/to/your/executable/warp'
```

This ensures the compilation and building of Warp tailored to your system.


## Installing a package
Use the following command to install a package:
```
warp install <package-name> <sysetm archetype>
```

### Example 
```
warp install pegasus amd64
```
To install a package from a GitHub repository:
```
warp install <package-name> --github/-G (clone-method) --http/--ssh
```
To install via docker
```
warp install <package-name> --docker/-D
```
