# SoDependencyAnalysiser

SoDependencyAnalysiser, as its name indicates, is a tool to analysis the depenedcy between to so and show dependency

## show dependency between so

``` bash
go run . --depender depender.so --dependee dependee.so

# output
depender.so -> dependee_1.so -> dependee.so
```

## show dependency 

``` bash
go run . --search-path /search/path --show-dependence-of target.so

# output
1 target.so
  1.1 level_1_dependency.so
    1.1.1 level_2_dependency.so
```