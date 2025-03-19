// Global state
let wasmLoaded = false;

// Initialize the WebAssembly module
const go = new Go();

async function initWasm() {
    const loadingElement = document.getElementById("loading");
    const errorMessageElement = document.getElementById("error-message");
    const errorTextElement = document.getElementById("error-text");
    const runButton = document.getElementById("run-button");
    const clearButton = document.getElementById("clear-button");
    const codeInput = document.getElementById("code-input");
    
    try {
        // Show loading screen
        loadingElement.classList.remove("hidden");
        
        // Fetch and instantiate the WebAssembly module
        const result = await WebAssembly.instantiateStreaming(
            fetch("olaf.wasm"), 
            go.importObject
        );
        
        // Run the Go WASM instance
        go.run(result.instance);
        console.log("WASM loaded successfully");
        
        // Update state and UI
        wasmLoaded = true;
        runButton.disabled = false;
        clearButton.disabled = false;
        
        // Initialize UI components
        initComponents();
        
        // Hide loading screen
        loadingElement.classList.add("hidden");
        
        // Focus the textarea after loading
        setTimeout(() => {
            codeInput.focus();
        }, 100);
    } catch (err) {
        console.error("Failed to load WASM:", err);
        
        // Show error message
        errorTextElement.textContent = "Failed to load the interpreter. Please check the console for details.";
        errorMessageElement.classList.add("visible");
        
        // Hide loading screen
        loadingElement.classList.add("hidden");
    }
}

function initComponents() {
    const runButton = document.getElementById("run-button");
    const clearButton = document.getElementById("clear-button");
    const codeInput = document.getElementById("code-input");
    const outputElement = document.getElementById("output");
    
    // Fix for cursor issues - ensure the textarea is properly initialized
    fixTextareaCursor(codeInput);
    
    // Run button click handler
    runButton.addEventListener("click", () => {
        executeCode(codeInput.value, outputElement);
    });
    
    // Clear button click handler
    clearButton.addEventListener("click", () => {
        codeInput.value = "";
        outputElement.textContent = "// Result will appear here";
        // Focus the textarea after clearing
        setTimeout(() => {
            codeInput.focus();
        }, 0);
    });
    
    // Example code click handlers
    const examples = document.querySelectorAll(".code-example");
    examples.forEach(example => {
        example.addEventListener("click", () => {
            codeInput.value = example.getAttribute("data-code");
            
            // Focus the textarea and place cursor at the end
            setTimeout(() => {
                codeInput.focus();
                codeInput.setSelectionRange(codeInput.value.length, codeInput.value.length);
            }, 0);
            
            if (window.innerWidth < 1024) {
                codeInput.scrollIntoView({ 
                    behavior: 'smooth' 
                });
            }
        });
    });
    
    // Keyboard shortcut for running code (Ctrl/Cmd + Enter)
    codeInput.addEventListener("keydown", (e) => {
        if ((e.ctrlKey || e.metaKey) && e.key === "Enter") {
            runButton.click();
            e.preventDefault(); // Prevent default behavior
        }
    });
}

// Function to fix cursor issues in textarea
function fixTextareaCursor(textarea) {
    // Make sure the textarea is clickable
    textarea.style.cursor = "text";
    
    // Ensure the textarea is properly initialized
    textarea.setAttribute("autocomplete", "off");
    textarea.setAttribute("autocorrect", "off");
    textarea.setAttribute("autocapitalize", "off");
    textarea.setAttribute("spellcheck", "false");
    
    // Add click handler to ensure focus
    textarea.addEventListener("click", function(e) {
        // This ensures the cursor position is set correctly when clicking
        if (document.activeElement !== this) {
            this.focus();
        }
    });
    
    // Add mousedown handler to prevent any issues with selection
    textarea.addEventListener("mousedown", function(e) {
        // Allow default behavior for text selection
        e.stopPropagation();
    });
    
    // Add a click handler to the parent container to ensure clicks propagate to the textarea
    const editorBox = textarea.closest('.editor-box');
    if (editorBox) {
        editorBox.addEventListener("click", function(e) {
            // If the click is on the editor box but not directly on the textarea,
            // and the textarea is not already focused, focus it
            if (e.target !== textarea && document.activeElement !== textarea) {
                textarea.focus();
                
                // Place cursor at end of text
                textarea.setSelectionRange(textarea.value.length, textarea.value.length);
            }
        });
    }
}

function executeCode(code, outputElement) {
    if (!wasmLoaded) {
        outputElement.textContent = "Error: Interpreter not loaded yet";
        return;
    }
    
    const trimmedCode = code.trim();
    
    if (!trimmedCode) {
        outputElement.textContent = "Error: No code to execute";
        return;
    }
    
    try {
        outputElement.textContent = "Running...";
        
        // Use setTimeout to allow the UI to update before running the code
        setTimeout(() => {
            try {
                // Call the evaluateOlaf function exposed by the WASM module
                const result = evaluateOlaf(trimmedCode);
                
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
}

// Start initialization
initWasm();

// Add window resize handler for responsive design
window.addEventListener('resize', () => {
    // Add any responsive adjustments here if needed
});