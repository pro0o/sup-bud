<a href="https://sup-bud.ddns.net/">
  <img src="./web/assets/banner.png" alt="sup-bud" />
</a>

## Notes
Well, [sup-bud](https://www.instagram.com/reel/Cz1jWt_uEJu/) is a tiny, intentionally simple programming language I created as part of my compiler design coursework. It’s built around writing an interpreter from scratch using go.</br>
Here's playground to try out sup-bud: https://sup-bud.ddns.net/ <br/>

---

## Running the Project

```bash
# with docker.
docker build -t sup-bud .
docker run -p 8080:8080 sup-bud

# without docker
./build.sh
```

## Credits

Here are few resources I referenced and learned from while building this project — in no particular order:

- [Writing an Interpreter in Go By Thorstan Bell](https://monkeylang.org/) [took much reference & easy to follow book written by a goated author.]

- [Grammars, parsing, and recursive descent (YouTube)](https://youtu.be/ENKT0Z3gldE?si=seZc5bsaGTnevbFD/)  [understand recursive descent parser better.]

- [Pratt Prasing (YouTube)](https://youtu.be/2l1Si4gSb9A?si=oJhwBtrq08jTxpJl/)  [viusalization of Pratt Parsing in OCaml.]
