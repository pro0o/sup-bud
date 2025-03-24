// Global state
let wasmLoaded = false
const commandHistory = []
let historyIndex = -1
const go = new Go()

// Declare evaluateOlaf as a global function, it will be defined by the WASM.
let evaluateOlaf

// Initialize the terminal UI
function initTerminal() {
  // Set up tabs (removed as no tabs are present in the current HTML)
  const tabs = document.querySelectorAll(".tab")
  const tabContents = document.querySelectorAll(".tab-content")

  if (tabs.length > 0) {
    tabs.forEach((tab) => {
      tab.addEventListener("click", () => {
        // Remove active class from all tabs and contents
        tabs.forEach((t) => t.classList.remove("active"))
        tabContents.forEach((c) => c.classList.remove("active"))

        // Add active class to clicked tab and corresponding content
        tab.classList.add("active")
        const tabId = `${tab.getAttribute("data-tab")}-tab`
        document.getElementById(tabId).classList.add("active")
      })
    })
  }

  // Update status bar time
  updateTime()
  setInterval(updateTime, 1000)
}

function updateTime() {
  const now = new Date()
  const timeString = now.toTimeString().split(" ")[0]
}

async function initWasm() {
  const loadingElement = document.getElementById("loading")
  const errorMessageElement = document.getElementById("error-message")
  const errorTextElement = document.getElementById("error-text")
  const runButton = document.getElementById("run-button")
  const clearButton = document.getElementById("clear-button")
  const codeInput = document.getElementById("code-input")
  const wasmStatusElement = document.getElementById("wasm-status")

  try {
    // Simulate a loading delay for visual effect
    await new Promise((resolve) => setTimeout(resolve, 1500))

    // Typewriter effect for loading text
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
    
    // Log the function
    console.log("evaluateOlaf:", typeof window.evaluateOlaf);
    if (typeof window.evaluateOlaf === "function") {
      evaluateOlaf = window.evaluateOlaf;
      console.log("evaluateOlaf assigned successfully.");
    } else {
      console.error("evaluateOlaf is not defined.");
    }

    // Update state and UI
    wasmLoaded = true
    runButton.disabled = false
    clearButton.disabled = false

    // Initialize UI components
    initComponents()

    // Hide loading screen with fade effect
    loadingElement.style.opacity = "0"
    setTimeout(() => {
      loadingElement.classList.add("hidden")
      loadingElement.style.opacity = "1"
    }, 500)

    // Focus the textarea after loading
    setTimeout(() => {
      codeInput.focus()
    }, 100)
  } catch (err) {
    console.error("Failed to load WASM:", err)

    // Show error message
    errorTextElement.textContent = "Failed to load the interpreter. Please check the console for details."
    errorMessageElement.classList.add("visible")
    
    // Only update wasmStatusElement if it exists
    if (wasmStatusElement) {
      wasmStatusElement.textContent = "Error"
      wasmStatusElement.style.color = "var(--terminal-red)"
    }

    // Hide loading screen
    loadingElement.classList.add("hidden")
  }
}

function initComponents() {
  const runButton = document.getElementById("run-button")
  const clearButton = document.getElementById("clear-button")
  const codeInput = document.getElementById("code-input")
  const outputElement = document.getElementById("output")

  // Fix for cursor issues - ensure the textarea is properly initialized
  fixTextareaCursor(codeInput)

  // Run button click handler
  runButton.addEventListener("click", () => {
    executeCode(codeInput.value, outputElement)
  })

  // Clear button click handler
  clearButton.addEventListener("click", () => {
    codeInput.value = ""
    outputElement.textContent = "// Output cleared"
    // Focus the textarea after clearing
    setTimeout(() => {
      codeInput.focus()
    }, 0)
  })

  // Example code click handlers
  const examples = document.querySelectorAll(".code-example")
  examples.forEach((example) => {
    example.addEventListener("click", () => {
      codeInput.value = example.getAttribute("data-code")

      // Focus the textarea and place cursor at the end
      setTimeout(() => {
        codeInput.focus()
        codeInput.setSelectionRange(codeInput.value.length, codeInput.value.length)
      }, 0)
    })
  })

  // Keyboard shortcut for running code (Ctrl/Cmd + Enter)
  codeInput.addEventListener("keydown", (e) => {
    if ((e.ctrlKey || e.metaKey) && e.key === "Enter") {
      runButton.click()
      e.preventDefault() // Prevent default behavior
    }

    // Command history navigation
    if (e.key === "ArrowUp") {
      if (historyIndex < commandHistory.length - 1) {
        historyIndex++
        codeInput.value = commandHistory[historyIndex]
        // Move cursor to end
        setTimeout(() => {
          codeInput.setSelectionRange(codeInput.value.length, codeInput.value.length)
        }, 0)
      }
      e.preventDefault()
    } else if (e.key === "ArrowDown") {
      if (historyIndex > 0) {
        historyIndex--
        codeInput.value = commandHistory[historyIndex]
        // Move cursor to end
        setTimeout(() => {
          codeInput.setSelectionRange(codeInput.value.length, codeInput.value.length)
        }, 0)
      } else if (historyIndex === 0) {
        historyIndex = -1
        codeInput.value = ""
      }
      e.preventDefault()
    }
  })
}

// Function to fix cursor issues in textarea
function fixTextareaCursor(textarea) {
  // Make sure the textarea is clickable
  textarea.style.cursor = "text"

  // Ensure the textarea is properly initialized
  textarea.setAttribute("autocomplete", "off")
  textarea.setAttribute("autocorrect", "off")
  textarea.setAttribute("autocapitalize", "off")
  textarea.setAttribute("spellcheck", "false")

  // Add click handler to ensure focus
  textarea.addEventListener("click", function (e) {
    // This ensures the cursor position is set correctly when clicking
    if (document.activeElement !== this) {
      this.focus()
    }
  })

  // Add mousedown handler to prevent any issues with selection
  textarea.addEventListener("mousedown", (e) => {
    // Allow default behavior for text selection
    e.stopPropagation()
  })
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
    // Add to command history if not already the most recent command
    if (commandHistory.length === 0 || commandHistory[0] !== trimmedCode) {
      commandHistory.unshift(trimmedCode)
      // Limit history size
      if (commandHistory.length > 50) {
        commandHistory.pop()
      }
    }
    historyIndex = -1

    // Show typing animation in output
    outputElement.textContent = ""
    const typingText = "Executing..."

    let i = 0
    const typingInterval = setInterval(() => {
      if (i < typingText.length) {
        outputElement.textContent += typingText[i]
        i++
      } else {
        clearInterval(typingInterval)

        // Execute code after typing animation
        setTimeout(() => {
          try {
            // Call the evaluateOlaf function exposed by the WASM module
            const result = evaluateOlaf(trimmedCode)

            // Typewriter effect for the result
            outputElement.textContent = ""
            const resultText = result.error ? "Error: " + result.error : result.result

            let j = 0
            const resultInterval = setInterval(() => {
              if (j < resultText.length) {
                outputElement.textContent += resultText[j]
                j++
              } else {
                clearInterval(resultInterval)
              }
            }, 10)
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

// Initialize the terminal UI
initTerminal()

// Start WASM initialization
initWasm()

// Add window resize handler for responsive design
window.addEventListener("resize", () => {
  // Adjust UI for different screen sizes if needed
})