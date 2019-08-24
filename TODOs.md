
# TODOs

- test multiple return values possible

- vartable creates env_vars even if nothing has changes (should reuse previous one)
- document variables and conversions

- resolve $values in variable names
- make field separators configurable
- non-string map keys
- make append to list possible

- BUG: you is empty: this should rather fail

``` yaml
targets:
  echo:
    run:
    - sh: echo $okay
      $: $myvar

  caller:
    run:
    - call:
        echo: {}
      opts: [ silent ]
      $: $yo
    - sh: echo $yo
```

- proper check of return should be possible: internal/parse/parse.go#L27

- select values: approachability, composability, extensibility, maintainability, simplicity, velocity
