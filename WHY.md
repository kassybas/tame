# Why not use shell scripts

The following examples have simple tasks that are more hassle than they should be.

Bash is great and close to being a standard, but boilerplate make code less readable and prone to mistakes.

Tame aims to make larger scripts more modular to make them more extendable, reusable and easier to debug.

Here are some of the painpoints in Shell/Make that Tame wants to address.

## Passing parameters

Passing parameters to a shell script or function comes with a lot of boilerplate and the syntax is less than obvious.

### Shell

While getopts makes things easier, it's still comes quite a lot of typing:

``` shell
while getopts ":a:p:" opt; do
  case $opt in
    f) format="$OPTARG"
    ;;
    m) msg="$OPTARG"
    ;;
    \?) echo "Invalid option -$OPTARG" >&2
    ;;
  esac
done
```

Positional arguments, like `foo="${1? error: missing foo arg}"` are realitively easy to write, but even easier to make mistakes when rearranging or adding arguments.

Default arguments are also not obvious.

While one in a script is fine, passing arguments to functions, adding the `while getopts` everywhere adds a lot of boilerplate, discouraging creating functions.

Also often leads to exp

### Tame

An example script that prepends the current time to our log message:

``` yaml
Log:
  args: {$msg, $format: "%T"}
  run:
  - sh: printf "[%s] %s" "$(date +${format})" "${msg}"
```

We can pass parameters from command line: `tame Log --msg "Hello world"`. Passing `format` is omitted to use the default value

Calling `Log` as a function from another target also has little overhead:

``` yaml
Foo:
  run:
  - call Log: {$msg: "my message"}
```

See more on targets in `example/`.

## Save output vs See output

Keeping the script output clean and understandable, while still passing errors can become complex.

Saving and printing both stdout and stderr can also be tricky, using using named pipes or tmeporary files

### Bash

``` shell
#  captures stderr, letting stdout through:
{ output=$(command 2>&1 1>&3-) ;} 3>&1
```

Capturing both stdout and stderr is close to impossible.

Controlling which fd (stdout or stderr) to let through and which to redirect to variable required some good insight into filedescriptors.

Keeping outputs clean and consistent, leads to the practice of redirecting outputs to `/dev/null`, but we lose debuggability.

To debug, we havet to reproduce errors in "debug" mode with `set -x`.


### Tame

Capturing but not printing both stdout and stderr.

``` yaml
Example:
  run:
  - sh: |
      echo "hello stderr" 1>&2
      echo "hello stdout"
    $: [$fooStdout, $barStderr]  # we capture both stdout and stderr, which can be used in proceeding steps
    # in case you only want to capture, not let stdout and stderr through, the 'silent' option can be set
    opts: silent
    # we can use stdout and stderr after
  - sh: echo "stdout:${fooStdout} stderr:${barStderr}"
```

## Returning results

Returning results in shell most often happens via stdout and stderr. This often creates the problem that subshells must keep silent and not log stats, otherwise their output would be less usable.

### Shell

```shell
foo(){
    echo "starting foo" # this message would diry the function output the output
    # this also means that we have to redirect useful output to /dev/null:
    command -v date > /dev/null || echo "date binary does not exist" # exiting now with 1 would require embedding in if
    date "+%T"
}
```

### Tame

We can execute as commands and print to stdout and still return only selected results

``` yaml
Foo:
  run:
  - sh: echo "starting foo"
  - sh: |
       command -v date && exit 0
       echo "date binary does not exist" && exit 1
    $: [$dateBin, null, $rc] # saving stdout, and return code
  - if: $rc != 0     # returning in case date binary is not needed
    return: null
  - sh: date "+%T"
    $: $currentDate
  - return: $currentDate

FooCaller:
  run:
  - call Foo: {}
    $
```

## Variable Scope









