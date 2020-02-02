
# TODOs
- tests
- logging
- create extensive documentation

- single-run target
- bug: cannot call included public target from CLI, when the included target also has includes
- log: log variables in dumped format: eg. "referencing field on a list:"

- overwrites in cli and in include
- make all field separators configurable
- while loop

- optimizations
    - vartable creates env_vars even if nothing has changes (should reuse previous one)
    - get all value of vartable creates all values even if unchanged
- tooling to syntaxcheck all tamefiles

- BUG: this passes parsing but prints nothing
``` yaml
      - if: true
        then:
          echo "FAILED"
        else:
          print: "OK
```
- BUG: race condition when returning error: next step is still being run
- BUG: this should create a new map, but it does nothing in for loop
``` yaml
            - runBlackBoxTest: { $tc: $tc, $path: $path }
              =: [$passed, $msg]
            - print: (format("PASSED %v", $passed))
            - $result[($path+"::"+$tc.name)]:
                passed: $passed
```
- BUG: async does not wait in parallel exec