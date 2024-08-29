# poc-cond-bridge

PoC for Conditional Bridge for OPI

This PoC contains a bridge code functionality to check the ability to use build tags to enable features
and capabilities in Go.  The code is based on a generic framework with redis support and allows enabling
the inventory collection for the target environment.

The build can be done using:

make compile - this will get the dependencies and build the binaries with the inventory enabled.

make get - this will get the dependencies

make base - this will build the binaries without the inventory enabled.

As additional capabilities are added to the PoC, they can be enabled with build tags and subsequent additions to the
Makefile can be done for individual or multiple services in the binaries.
