<a href="https://sup-bud.onrender.com/">
  <img src="./web/assets/banner.png" alt="sup' bud" />
</a>

## Intro...
Well, [sup' bud](https://www.instagram.com/reel/Cz1jWt_uEJu/) is a tiny, intentionally simple programming language I created as part of my compiler design coursework. It’s built around writing an interpreter from scratch using go.</br>
Try it yourself here: https://sup-bud.onrender.com/ <br/>
There are also a few code snippets to try out, just to get a feel for the language.

---

## Features so far 

- [x] Variable Declaration using *sup*.
- [x] Arithmetic & Logic.
- [x] Boolean Logic.
- [x] If/Else Expressions.
- [x] Expression handling using Pratt Parsing.
- [x] Block stmts.
- [x] Return stmt.
- [x] Funcs & Closures using *bud*.
- [x] Scoping & env.
- [x] Runtime Object System [Int, Bool & Null].
- [x] Minimal Web Interface, i love it.
- [x] Inline error throwing.
- [x] Syntax Highlighting.
- [x] Recursion depth limit.
- [x] Unit-test driven approach.
- [ ] REPL-CLI.
- [ ] Built-in funcs.
- [ ] I/O operations.
- [ ] More Concise Error handling.
- [ ] Sequel of the Monkey lang??.

---

## Project Structure
- `ast/`       - AST definitions  
- `cmd/`       - Main entry point  
- `eval/`      - Evaluator logic  
- `lexer/`     - Lexical analysis  
- `object/`    - Object/typ/env system  
- `parser/`    - Parser implementation  
- `repl/`      - CLI (in progress)  
- `token/`     - Token definitions  
- `web/`       - Minimal Web UI   

## Running the Project

```bash
# with docker.
docker build -t sup-bud .
docker run -p 8080:8080 sup-bud

# without docker
./build.sh
```

## Resources 

Here are few resources I referenced and learned from while building this project — in no particular order:

- [Writing an Interpreter in Go By Thorstan Bell](https://monkeylang.org/) [took much reference & easy to follow book written by a goated author.]

- [Grammars, parsing, and recursive descent (YouTube)](https://youtu.be/ENKT0Z3gldE?si=seZc5bsaGTnevbFD/)  [understand recursive descent parser better.]

- [Pratt Prasing (YouTube)](https://youtu.be/2l1Si4gSb9A?si=oJhwBtrq08jTxpJl/)  [viusalization of Pratt Parsing in OCaml.]

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contributing

If you have any suggestions or improvements, please create an issue or a pull request. I'll try to respond to all issues and pull requests.
