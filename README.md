# ReASM

[Wiki](https://github.com/AsynchronousAI/reasm/wiki)
> An experimental **RISC-V IM** compatible assembler/disassembler to **Luau**.
## Example:
```c
void printf(const char *, ...); /* manually define printf if we are not using stdlib.h */

int fib(int n) {
    if (n <= 1)
        return n;
    return fib(n-1) + fib(n-2);
}

void printFib(int n, int i) {
    if (i < n) {
        printf("%d ", fib(i));
        printFib(n, i+1);
    }
}

int main() {
    int terms = 10;
    printFib(terms, 0);
    return 0;
}
```
```bash
clang -S -target riscv32 -march=rv32im main.c -o main.s
asmdecomp main.s -o main.luau # where the magic happens
luau main.luau
```

## Usage:
```bash
asmdecomp main.S -o main.luau --mode {module|main|bench} --trace --comments
```

Input file can either be a `.S` assembly file, or a `.elf` file which is linked *(experimental)*.

### Options
- `--comments`: This will place comments all around the generated code with details such as the instruction's purpose, operands, and any relevant debug information.
- `--trace`: Everytime a jump happens it will log to output, this is a more extreme option and should only be used for debug.
- `--mode`:
  * `module` will automatically expose memory, API to inject functions, and registers to whoever imports.
  * `main` will generate a simple Luau file which runs on its own.
  * `bench` will generate a module prepared for benchmarking with [Scriptbench](https://devforum.roblox.com/t/scriptbench-free-opensource-heavy-duty-benchmarker/3815286) or [Benchmarker](https://devforum.roblox.com/t/benchmarker-plugin-compare-function-speeds-with-graphs-percentiles-and-more/829912).

## Resources:
Super helpful resources in development below:
- https://www.cs.sfu.ca/~ashriram/Courses/CS295/assets/notebooks/RISCV/RISCV_CARD.pdf
- https://msyksphinz-self.github.io/riscv-isadoc/
- https://godbolt.org/

## TODO:
- Floating point & Double support
- Handle overflows correctly.
- Vector operations
- Work on support for ELF files. (or decide to remove it)
  * `No bindings for function '<register??>'`
  * ECALL & AUIPC
  * Read labels
