# **Java Class File Parser in GO**

This project was created for the purpose of analyzing and studying the [Java Class File format](https://docs.oracle.com/javase/specs/jvms/se23/html/jvms-4.html). It provides functionality to parse `.class` files, breaking down their structure and extracting meaningful information.

Additionally, this parsing logic could be repurposed as a foundation for regenerating `.class` files. Such regenerated files could be loaded and interpreted by the JVM, opening possibilities for dynamic class manipulation and custom bytecode generation.

## **Quick Start**

```console
$ javac App.java [optional]
$ go run ./Main.go
```