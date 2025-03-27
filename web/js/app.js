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
    console.log("Initializing WASM...");
    const result = await WebAssembly.instantiateStreaming(fetch("olaf.wasm"), go.importObject);
    console.log("WASM Instance:", result.instance);
    
    go.run(result.instance);
    console.log("WASM Started");
    
    console.log("evaluateOlaf:", typeof window.evaluateOlaf);
    if (typeof window.evaluateOlaf === "function") {
      evaluateOlaf = window.evaluateOlaf;
      console.log("evaluateOlaf assigned successfully.");
    } else {
      console.error("evaluateOlaf is not defined.");
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

function initComponents() {
  const runButton = document.getElementById("run-button")
  const clearButton = document.getElementById("clear-button")
  const codeInput = document.getElementById("code-input")
  const outputElement = document.getElementById("output")

  fixTextareaCursor(codeInput)

  runButton.addEventListener("click", () => {
    executeCode(codeInput.value, outputElement)
  })

  clearButton.addEventListener("click", () => {
    codeInput.value = ""
    outputElement.textContent = "// Output cleared"
    setTimeout(() => {
      codeInput.focus()
    }, 0)
  })
}

function fixTextareaCursor(textarea) {
  textarea.style.cursor = "text"

  textarea.setAttribute("autocomplete", "off")
  textarea.setAttribute("autocorrect", "off")
  textarea.setAttribute("autocapitalize", "off")
  textarea.setAttribute("spellcheck", "false")
}

function executeCode(code, outputElement) {
  if (!wasmLoaded) {
    outputElement.textContent = "Error: Interpreter not loaded yet"
    return
  }

  const trimmedCode = code.trim()

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

initWasm()
