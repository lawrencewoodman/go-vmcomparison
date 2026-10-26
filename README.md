# go-vmcomparison

A collection of Virtual Machine implementations used to compare benchmarks

In order to compare different types of virtual machines it can be useful
to compare their use for comfort and performance.  This collection is
one such endeavour. 

## Running the Benchmarks

Each VM has a benchmark routine that runs a set of programs on the VM and
reports the speed of each program along with the size of the code in the VM.

The programs to be run on the VMs can be found in the `fixtures/` directory of
each VM and is assembled by the symbolic assembler provided in the VM
directory.

All the benchmarks can be be run without running the tests using:

    $ go test ./... -bench . -run xxx -timeout 0

If the output of the above is redirected to a file, then they can be collected
and turned into a table using `cmd/collect-benchmarks` as follows:

    $ go test ./... -bench . -run xxx -timeout 0 > raw_results.txt
    $ go run cmd/collect-benchmarks/main.go tables raw_results.txt



## Licence
Copyright (C) 2023, Lawrence Woodman <lwoodman@vlifesystems.com>

This software is licensed under an MIT Licence.  Please see the file, LICENCE.md, for details.
