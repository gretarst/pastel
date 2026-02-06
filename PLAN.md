# Pastel Interpreter Development Plan

This document outlines a plan for cleaning up, refactoring, and extending the Pastel interpreter.

## 1. Project Overview

Pastel is a simple interpreter for a subset of the Pascal language, written in Go. It currently consists of a lexer, a parser that builds an Abstract Syntax Tree (AST), and an interpreter that evaluates the AST.

**Current Features:**
- `program` and `var` declarations (integers only).
- `begin...end.` blocks for program body.
- Integer literals and arithmetic operations (`+`, `-`, `*`, `/`).
- `:=` for assignments.
- `writeln` for printing.
- Basic error reporting during parsing and runtime.

## 2. Cleanup and Refactoring

This section focuses on improving the quality and maintainability of the existing code.

- **Add Automated Testing:**
  - Implement unit tests for the `lexer` to verify correct tokenization of various inputs.
  - Write unit tests for the `parser` to ensure it correctly constructs the AST for valid programs and reports errors for invalid ones.
  - Add tests for the `interpreter` to check that programs execute correctly and produce the expected output.

- **Improve Error Handling:**
  - The `main` function currently panics on file read errors or if no file is provided. It should handle these errors gracefully and print a user-friendly message.
  - The interpreter also has the potential to panic. Refactor it to return errors consistently, similar to how the parser does.
  - Enhance the lexer to track line and column numbers. This information can then be used in parser and interpreter errors to provide more precise locations for a better developer experience.

- **Strengthen the AST with Type Safety:**
  - The `parser.Expr` and `parser.Stmt` types are currently `any` (or `interface{}`). Define these as interfaces with methods (e.g., `exprNode()` and `stmtNode()`) and have the AST node structs implement them. This will make the AST more type-safe and reduce the need for type assertions.

- **Restructure Packages:**
  - The AST definitions in `parser/ast.go` are used by both the `parser` and `interpreter`. Consider moving the AST-related code into its own `ast` package to make dependencies clearer.

- **Code Formatting and Linting:**
  - Run `gofmt` on the entire codebase to ensure consistent formatting.
  - Introduce a linter like `golangci-lint` to catch potential bugs and style issues.

- **Remove Unused Code:**
  - The `PrintExpr` function in `parser/parser.go` appears to be a debugging utility and is not used. It can be removed.
  - The file `.vscode/GTEMP.pas` seems temporary. It should either be removed or added to a `.gitignore` file.

- **Refactor the Interpreter:**
  - The interpreter is currently a collection of functions. It could be refactored into a struct that holds the environment, making the design more object-oriented.

## 3. Potential Next Steps

This section lists ideas for new features to extend the capabilities of the interpreter.

- **Expand Data Types:**
  - Add support for `real` (floating-point numbers).
  - Implement the `boolean` type with `true` and `false` literals.
  - Introduce `char` and `string` types.

- **Implement More Control Structures:**
  - `if-then-else` conditional statements.
  - Loop constructs: `for`, `while`, and `repeat-until`.

- **Add Support for Procedures and Functions:**
  - Allow declaration and calling of procedures and functions.
  - Implement parameter passing (both by value and by reference).
  - Handle return values from functions.

- **Enhance the Language Grammar:**
  - Support comments, both single-line `{...}` and multi-line `(*...*)`.
  - Add more operators:
    - Integer division (`div`) and modulo (`mod`).
    - Relational operators: `=`, `<>`, `<`, `>`, `<=`, `>=`.
    - Logical operators: `and`, `or`, `not`.

- **Introduce Composite Types:**
  - Implement support for `array`s.
  - Add `record` types for more complex data structures.

- **Build a More Robust Type System:**
  - Implement a type-checking phase that runs after parsing to catch type mismatches before execution.

- **Expand the Standard Library:**
  - Implement more built-in Pascal functions and procedures like `readln`, `sqrt`, `abs`, etc.

- **Create an Interactive REPL:**
  - Develop a Read-Eval-Print Loop (REPL) mode for the interpreter. This would allow for interactive experimentation without needing to create a file for every program.
