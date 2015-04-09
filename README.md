GOpenMP
=======

A set of libraries and tools that's implements OpenMP interface in Go! language.

Install
-------

Install using go get

    go get github.com/DeusCoNWeT/GOpenMP_Project/GoMP

and this will build the `GoMP` binary in `$GOPATH/bin`.

If you have build problems with auxiliar packages ("goprep", "var_processor", "import_processor", etc), try to revise your $GOPATH, or modify the import declarations in code.

It will also pull in a set of examples that use GoMP, many of them in serial and parallel version.

Using GoMP
----------

To use GoMP, type:

    $ GoMP input_file output_file
  
I use two parameters: an input file that your want to parallelize, including pragma_gomp directives, and a new output file witch contain the original code rewrite for parallel execution.

You can use the complete route or relative to indicate the input and output files.

If you want rewrite the original code just put the name of the input:file as second parameter. But be careful!! You lose the original code.
