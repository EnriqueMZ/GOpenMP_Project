GOpenMP
=======

A set of libraries and tools that's implements OpenMP interface in Go language.

Install
-------

Install using go get

    go get github.com/DeusCoNWeT/GOpenMP_Project/GoMP

and this will build the `GoMP` binary in `$GOPATH/bin`.

If you have build problems with auxiliar packages ("goprep", "var_processor", "import_processor", etc), try to revise your $GOPATH, or modify the import declarations in code.

It will also pull in a set of examples that use GoMP, many of them in serial and parallel version.

What is GOpenMP?
----------------

Simple: Go + OpenMP
 
A set of libraries and tools, supported by the characteristics and models of concurrency implemented in Go, that add functionalities that are typical features of standard OpenMP: An Application Programming Interface (API) flexible, portable and scalable, that supports multi-platform shared memory multiprocessing programming in C, C++, or Fortran, and Multiprocessor Architectures oriented.

OpenMP is based on "fork -join" model, paradigm comes from Unix systems, where a task is divided into K-threads ("fork") with less weight, then "collect" their results at the end and unite them in one result ("join").

Syntactically, OpenMP consists of compiler directives, called "pragmas". These directives are included in the code and determine the behavior of the same. Can be incorporated into an existing code and modify its execution without adding extra code. Not only that, but in the event that no parallel execution is desired, the compiler can simply ignore them.

For more information, please consult the official website of OpenMP:

[OpenMP Official Website](http://openmp.org/wp/)

The main idea is to add features that allow programming in Go language using the syntax that OpenMP provides in other languages, like Fortran and C/C++ (ie, similar to the latter), thus providing the Go users with tools to program parallel processing structures in a simple and transparent way, in the same way you would use C/C++.

The conceptual idea behind GOpenMP library is a code preprocessor module. This module takes an original source code, written in Go, and which have been added various GopenMP directives, and becomes a new source code, also in Go, which when it is compiled and executed, behaves in parallel.

Using GoMP
----------

To use GoMP, type:

    $ GoMP input_file output_file
  
I use two parameters: an input file that your want to parallelize, including pragma_gomp directives, and a new output file witch contain the original code rewrite for parallel execution.

You can use the complete route or relative to indicate the input and output files.

If you want rewrite the original code just put the name of the input_file as second parameter. But be careful!! You lose the original code.

Pragmas GoMP
------------

Pragmas (compiler directives) in GOpenMP follow the same structure as in C / C ++. In this case, pragmas are preceded by double bar, like comments on Go. Thus, if parallel execution is not desired , the compiler simply ignored as comments.

The basic structure of a pragma in GOpenMP is:

    //pragma gomp pragma_type [pragma_clause [[,] pragma_clause...]]
    
Here we can see an example:

    //pragma gomp parallel for default(none) shared(cont) private(i, n) reduction(+:sum)
    
Where "parallel for" is the type of pragma that applies in this case, and "default", "shared", "private" and "reduction" are clauses applicable to this pragma, each with its variables (or operators) associated.

Pragmas and Clauses
-------------------

In its latest version , the preprocessor supports the following pragmas :

  * `Parallel`
  * `For`
  * `Parallel for`
  
Additionally , the preprocessor supports "nesting": parallel regions within others; or pragma "for" into pragma "parallel" (It behaves in a similar way if we use the pragma "parallel for ").

It also recognised the following clauses:

  * `Num_threads` (Number of threads in parallel region).
  * `Default` (variable declaration mode in pragmas).
  * `Shared` (variables marked as shared).
  * `Private` (variables marked as private, default initialized).
  * `Firstprivate` (variables marked as private, initialized with the value they had in the main thread).
  * `Reduction` (reduction operator and variables).
  * `If` (parallel clause condition).
  
Although " pragma_processor " is capable of recognizing more clauses , are not yet fully operational.
