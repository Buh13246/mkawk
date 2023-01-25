# mkawk
The mkawk tool

This was a Tool written 1989 by M.Wellner to generate awk Scripts from simple templates with a documented simple templating language.
See the man file (mkawk.LOCAL) for details on usage.

The Tool was originally written in C for an old Unix-based system.

The reimplementation was done by me (Philipp Wellner) on January 2023

# Compiling
You need go version 1.19 for this to compile, but if you don't have that, you can test if lowering the version in `go.mod` compiles on older versions
```sh
go build .
```
