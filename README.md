# Example

```aml
func get_fib() {
    func fib(n) {
        return n > 2 ? fib(n - 1) + fib(n - 2) : n;
    }
    return fib;
}
var fib = get_fib();
var start = time();
for (var i = 0; i < 10; i = i + 1) {
    print fib(i);
}
print time() - start;
```

# Usage

```bash
git clone https://github.com/zaixrx/AML-lang.git
cd ./AML-lang && go build -o aml
```

you can either use repl
```bash
./aml -repl

>> print "Hello, World!";
"Hello, World!"
```

or interpret an existing file
```bash
./aml ./examples/helloworld.aml

"Hello, World!"
```

# Resources
- crafting interpreters: \
    https://craftinginterpreters.com
- I'm not a go developer: \
    https://go.dev/ \
    https://gobyexample.com/
	.
* useful theory:
    - Regular Expressions: \
        https://web.stanford.edu/class/archive/cs/cs103/cs103.1208/lectures/14-RegExes/Regular%20Expressions.pdf
    - Context Free Grammars: \
        https://web.stanford.edu/class/archive/cs/cs103/cs103.1164/lectures/18/Small18.pdf
