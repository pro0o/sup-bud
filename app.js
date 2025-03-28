let wasmLoaded = false
const go = new Go()

let evaluateOlaf

async function initWasm() {
  const loadingElement = document.getElementById("loading")
  const errorMessageElement = document.getElementById("error-message")
  const errorTextElement = document.getElementById("error-text")
  const runButton = document.getElementById("run-button")
  const clearButton = document.getElementById("clear-button")
  const codeInput = document.getElementById("code-input")
  const wasmStatusElement = document.getElementById("wasm-status")

  try {
    await new Promise((resolve) => setTimeout(resolve, 1500))

    const loadingTextElement = document.querySelector(".typed-text")
    const originalText = loadingTextElement.textContent
    loadingTextElement.textContent = ""

    for (let i = 0; i < originalText.length; i++) {
      await new Promise((resolve) => setTimeout(resolve, 30))
      loadingTextElement.textContent += originalText[i]
    }

    await new Promise((resolve) => setTimeout(resolve, 500))
    console.log("Initializing WASM...")
    const result = await WebAssembly.instantiateStreaming(fetch("sup-bud.wasm"), go.importObject)
    console.log("WASM Instance:", result.instance)

    go.run(result.instance)
    console.log("WASM Started")

    console.log("evaluateOlaf:", typeof window.evaluateOlaf)
    if (typeof window.evaluateOlaf === "function") {
      evaluateOlaf = window.evaluateOlaf
      console.log("evaluateOlaf assigned successfully.")
    } else {
      console.error("evaluateOlaf is not defined.")
    }

    wasmLoaded = true
    runButton.disabled = false
    clearButton.disabled = false

    initComponents()

    loadingElement.style.opacity = "0"
    setTimeout(() => {
      loadingElement.classList.add("hidden")
      loadingElement.style.opacity = "1"
    }, 500)

    setTimeout(() => {
      codeInput.focus()
    }, 100)
  } catch (err) {
    console.error("Failed to load WASM:", err)

    errorTextElement.textContent = "Failed to load the interpreter. Please check the console for details."
    errorMessageElement.classList.add("visible")

    if (wasmStatusElement) {
      wasmStatusElement.textContent = "Error"
      wasmStatusElement.style.color = "var(--terminal-red)"
    }

    loadingElement.classList.add("hidden")
  }
}

function applySyntaxHighlighting(textarea) {
  const code = textarea.value
  const tokens = tokenize(code)

  let highlightDiv = document.getElementById("syntax-highlight-container")
  if (!highlightDiv) {
    highlightDiv = document.createElement("div")
    highlightDiv.id = "syntax-highlight-container"
    highlightDiv.className = "syntax-highlight-container"
    textarea.parentNode.insertBefore(highlightDiv, textarea)
  }

  highlightDiv.innerHTML = ""

  tokens.forEach((token) => {
    const span = document.createElement("span")
    span.textContent = token.text
    span.className = `token-${token.type}`
    highlightDiv.appendChild(span)
  })

  highlightDiv.scrollTop = textarea.scrollTop
  highlightDiv.scrollLeft = textarea.scrollLeft
}

function setupScrollSync() {
  const textarea = document.getElementById("code-input")
  const highlightDiv = document.getElementById("syntax-highlight-container")

  if (textarea && highlightDiv) {
    textarea.addEventListener("scroll", () => {
      highlightDiv.scrollTop = textarea.scrollTop
      highlightDiv.scrollLeft = textarea.scrollLeft
    })

    window.addEventListener("resize", () => {
      highlightDiv.scrollTop = textarea.scrollTop
      highlightDiv.scrollLeft = textarea.scrollLeft
    })

    highlightDiv.scrollTop = textarea.scrollTop
    highlightDiv.scrollLeft = textarea.scrollLeft
  }
}

function setupTabHandling() {
  const textarea = document.getElementById("code-input")

  if (textarea) {
    textarea.addEventListener("keydown", (e) => {
      if (e.key === "Tab") {
        e.preventDefault()

        const start = textarea.selectionStart
        const end = textarea.selectionEnd

        textarea.value = textarea.value.substring(0, start) + "  " + textarea.value.substring(end)

        textarea.selectionStart = textarea.selectionEnd = start + 2

        applySyntaxHighlighting(textarea)
      }
    })
  }
}

function initComponents() {
  const runButton = document.getElementById("run-button")
  const clearButton = document.getElementById("clear-button")
  const codeInput = document.getElementById("code-input")
  const outputElement = document.getElementById("output")
  const snippetButtons = document.querySelectorAll(".code-example")
  const MAX_CODE_LENGTH = 200

  fixTextareaCursor(codeInput)
  setupTabHandling()

  codeInput.addEventListener("keydown", (event) => {
    if (event.ctrlKey && event.key === "Enter") {
      event.preventDefault()
      executeCode(codeInput.value, outputElement, MAX_CODE_LENGTH)
    }
  })

  codeInput.addEventListener("input", () => {
    applySyntaxHighlighting(codeInput)
  })

  applySyntaxHighlighting(codeInput)

  runButton.addEventListener("click", () => {
    executeCode(codeInput.value, outputElement, MAX_CODE_LENGTH)
  })

  clearButton.addEventListener("click", () => {
    codeInput.value = ""
    outputElement.textContent = "// Output cleared"
    applySyntaxHighlighting(codeInput)
    setTimeout(() => {
      codeInput.focus()
    }, 0)
  })

  snippetButtons.forEach((button) => {
    button.addEventListener("click", () => {
      codeInput.value = button.getAttribute("data-code")
      applySyntaxHighlighting(codeInput)
      codeInput.focus()
    })
  })
}

function fixTextareaCursor(textarea) {
  textarea.style.cursor = "text"

  textarea.setAttribute("autocomplete", "off")
  textarea.setAttribute("autocorrect", "off")
  textarea.setAttribute("autocapitalize", "off")
  textarea.setAttribute("spellcheck", "false")
}

function executeCode(code, outputElement, maxLength = 1000) {
  if (!wasmLoaded) {
    outputElement.textContent = "Error: Interpreter not loaded yet"
    return
  }

  const trimmedCode = code.trim()

  if (trimmedCode.length > maxLength) {
    outputElement.textContent = `Error: Code exceeds maximum length of code. Hold onn—`
    return
  }

  if (!trimmedCode) {
    outputElement.textContent = "Error: No code to execute"
    return
  }

  try {
    outputElement.textContent = ""
    const typingText = "Executing..."

    let i = 0
    const typingInterval = setInterval(() => {
      if (i < typingText.length) {
        outputElement.textContent += typingText[i]
        i++
      } else {
        clearInterval(typingInterval)

        setTimeout(() => {
          try {
            const result = evaluateOlaf(trimmedCode)
            outputElement.textContent = result.error ? "Error: " + result.error : result.result
          } catch (error) {
            outputElement.textContent = "Error: " + error.message
            console.error(error)
          }
        }, 300)
      }
    }, 50)
  } catch (error) {
    outputElement.textContent = "Error: " + error.message
    console.error(error)
  }
}

const TOKEN_TYPES = {
  KEYWORD: ["sup", "bud", "true", "false", "if", "else", "return"],
  OPERATOR: ["+", "-", "*", "/", "=", "==", "!=", "<", ">"],
  PUNCTUATION: [",", ";", "(", ")", "{", "}"],
  NUMBER: /^\d+$/,
  IDENTIFIER: /^[a-zA-Z_][a-zA-Z0-9_]*$/,
}

function tokenize(code) {
  const tokens = []
  let currentToken = ""
  const currentType = null

  function pushCurrentToken() {
    if (currentToken) {
      let type = "identifier"

      if (TOKEN_TYPES.KEYWORD.includes(currentToken)) {
        type = "keyword"
      } else if (currentToken.match(TOKEN_TYPES.NUMBER)) {
        type = "number"
      }

      tokens.push({
        text: currentToken,
        type: type,
      })

      currentToken = ""
    }
  }

  for (let i = 0; i < code.length; i++) {
    const char = code[i]

    if (TOKEN_TYPES.OPERATOR.includes(char) || TOKEN_TYPES.PUNCTUATION.includes(char)) {
      if (i + 1 < code.length) {
        const twoChars = char + code[i + 1]
        if (TOKEN_TYPES.OPERATOR.includes(twoChars)) {
          pushCurrentToken()
          tokens.push({
            text: twoChars,
            type: "operator",
          })
          i++
          continue
        }
      }

      pushCurrentToken()
      tokens.push({
        text: char,
        type: TOKEN_TYPES.OPERATOR.includes(char) ? "operator" : "punctuation",
      })
      continue
    }

    if (char === " " || char === "\t" || char === "\n") {
      pushCurrentToken()
      tokens.push({
        text: char,
        type: "whitespace",
      })
      continue
    }

    currentToken += char
  }

  pushCurrentToken()

  return tokens
}

initWasm()

setTimeout(setupScrollSync, 1000)