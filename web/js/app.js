// Initialize the Go WASM runtime
const go = new Go();

// WASM module loading
async function initWasm() {
    try {
        const result = await WebAssembly.instantiateStreaming(
            fetch("olaf.wasm"), 
            go.importObject
        );
        go.run(result.instance);
        console.log("WASM loaded successfully");
        
        // Initialize components after WASM is loaded
        initComponents();
    } catch (err) {
        console.error("Failed to load WASM:", err);
        document.getElementById("output").textContent = 
            "Failed to load the interpreter. Please check the console for details.";
    }
}

// Initialize UI components and event handlers
function initComponents() {
    // Editor component handlers
    const runButton = document.getElementById("run-button");
    const clearButton = document.getElementById("clear-button");
    const codeInput = document.getElementById("code-input");
    const outputElement = document.getElementById("output");
    
    // Run button click handler
    runButton.addEventListener("click", () => {
        const code = codeInput.value.trim();
        
        if (!code) {
            outputElement.textContent = "Error: No code to execute";
            return;
        }
        
        try {
            // Show loading indicator
            outputElement.textContent = "Running...";
            
            // Execute with a small delay to allow UI to update
            setTimeout(() => {
                try {
                    const result = evaluateOlaf(code);
                    
                    if (result.error) {
                        outputElement.textContent = "Error: " + result.error;
                    } else {
                        outputElement.textContent = result.result;
                    }
                } catch (error) {
                    outputElement.textContent = "Error: " + error.message;
                    console.error(error);
                }
            }, 10);
        } catch (error) {
            outputElement.textContent = "Error: " + error.message;
            console.error(error);
        }
    });
    
    // Clear button click handler
    clearButton.addEventListener("click", () => {
        codeInput.value = "";
        outputElement.textContent = "// Result will appear here";
    });
    
    // Examples component handlers
    const examples = document.querySelectorAll(".code-example");
    examples.forEach(example => {
        example.addEventListener("click", () => {
            codeInput.value = example.getAttribute("data-code");
            // Scroll to editor on mobile
            if (window.innerWidth < 1024) {
                document.getElementById("code-input").scrollIntoView({ 
                    behavior: 'smooth' 
                });
            }
        });
    });
    
    // Add keyboard shortcut (Ctrl+Enter or Cmd+Enter to run code)
    codeInput.addEventListener("keydown", (e) => {
        if ((e.ctrlKey || e.metaKey) && e.key === "Enter") {
            runButton.click();
        }
    });
}

// Start loading the WASM module
initWasm();