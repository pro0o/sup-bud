let wasmLoaded = false;
const go = new Go();

let evaluateSupBud;

async function initWasm() {
  const loadingElement = document.getElementById("loading");
  const errorMessageElement = document.getElementById("error-message");
  const errorTextElement = document.getElementById("error-text");
  const codeInput = document.getElementById("code-input");

  try {
    await new Promise((resolve) => setTimeout(resolve, 500));
    console.log("Initializing WASM...");
    const result = await WebAssembly.instantiateStreaming(fetch("./sup-bud.wasm"), go.importObject);
    console.log("WASM Instance:", result.instance);

    go.run(result.instance);
    console.log("WASM Started");

    console.log("evaluateSupBud:", typeof window.evaluateSupBud);
    if (typeof window.evaluateSupBud === "function") {
      evaluateSupBud = window.evaluateSupBud;
      console.log("evaluateSupBud assigned successfully.");
    } else {
      console.error("evaluateSupBud is not defined.");
      throw new Error("evaluateSupBud function not found");
    }

    wasmLoaded = true;
    initComponents();

    loadingElement.style.opacity = "0";
    setTimeout(() => {
      loadingElement.classList.add("hidden");
      loadingElement.style.opacity = "1";
    }, 500);

    setTimeout(() => {
      codeInput.focus();
    }, 100);
  } catch (err) {
    console.error("Failed to load WASM:", err);

    errorTextElement.textContent = "Failed to load the interpreter. Please check the console for details.";
    errorMessageElement.classList.remove("hidden");
    errorMessageElement.classList.add("visible");

    loadingElement.classList.add("hidden");
  }
}

function applySyntaxHighlighting(textarea) {
  const code = textarea.value;
  const tokens = tokenize(code);

  let highlightDiv = document.getElementById("syntax-highlight-container");
  if (!highlightDiv) {
    highlightDiv = document.createElement("div");
    highlightDiv.id = "syntax-highlight-container";
    highlightDiv.className = "syntax-highlight-container";
    textarea.parentNode.insertBefore(highlightDiv, textarea);
  }

  highlightDiv.innerHTML = "";

  tokens.forEach((token) => {
    const span = document.createElement("span");
    span.textContent = token.text;
    span.className = `token-${token.type}`;
    highlightDiv.appendChild(span);
  });

  highlightDiv.scrollTop = textarea.scrollTop;
  highlightDiv.scrollLeft = textarea.scrollLeft;
}

function setupScrollSync() {
  const textarea = document.getElementById("code-input");
  const highlightDiv = document.getElementById("syntax-highlight-container");

  if (textarea && highlightDiv) {
    textarea.addEventListener("scroll", () => {
      highlightDiv.scrollTop = textarea.scrollTop;
      highlightDiv.scrollLeft = textarea.scrollLeft;
    });

    window.addEventListener("resize", () => {
      highlightDiv.scrollTop = textarea.scrollTop;
      highlightDiv.scrollLeft = textarea.scrollLeft;
    });

    highlightDiv.scrollTop = textarea.scrollTop;
    highlightDiv.scrollLeft = textarea.scrollLeft;
  }
}
function initComponents() {
  const codeInput = document.getElementById("code-input");
  const outputElement = document.getElementById("output");
  const snippetButtons = document.querySelectorAll(".code-example");
  const MAX_CODE_LENGTH = 200;
  const executeShortcut = document.querySelector(".shortcut-hint");

  fixTextareaCursor(codeInput);

  executeShortcut.innerHTML = '<span class="clickable-shortcut">CTRL+ENTER</span> to execute | <span class="clickable-shortcut">CTRL+D</span> to clear';
  
  executeShortcut.addEventListener("click", (event) => {
    if (event.target.classList.contains("clickable-shortcut")) {
      const text = event.target.textContent;
      
      if (text === "CTRL+ENTER") {
        executeCode(codeInput.value, outputElement, MAX_CODE_LENGTH);
      } else if (text === "CTRL+D") {
        clearCode(codeInput, outputElement);
      }
    }
  });

  codeInput.addEventListener("keydown", (event) => {
    if ((event.ctrlKey || event.metaKey) && event.key === "Enter") {
      event.preventDefault();
      executeCode(codeInput.value, outputElement, MAX_CODE_LENGTH);
    }
    
    if ((event.ctrlKey || event.metaKey) && event.key === "d") {
      event.preventDefault();
      clearCode(codeInput, outputElement);
    }
  });

  codeInput.addEventListener("input", () => {
    applySyntaxHighlighting(codeInput);
  });

  applySyntaxHighlighting(codeInput);

  snippetButtons.forEach((button) => {
    button.addEventListener("click", () => {
      codeInput.value = button.getAttribute("data-code");
      applySyntaxHighlighting(codeInput);
      codeInput.focus();
    });
  });
}
function clearCode(codeInput, outputElement) {
  codeInput.value = "";
  outputElement.textContent = "// Output cleared";
  applySyntaxHighlighting(codeInput);
  setTimeout(() => {
    codeInput.focus();
  }, 0);
}

function fixTextareaCursor(textarea) {
  textarea.style.cursor = "text";

  textarea.setAttribute("autocomplete", "off");
  textarea.setAttribute("autocorrect", "off");
  textarea.setAttribute("autocapitalize", "off");
  textarea.setAttribute("spellcheck", "false");
}

function executeCode(code, outputElement, maxLength = 1000) {
  if (!wasmLoaded) {
    outputElement.textContent = "Error: Interpreter not loaded yet";
    return;
  }

  const trimmedCode = code.trim();

  if (trimmedCode.length > maxLength) {
    outputElement.textContent = `Error: Code exceeds maximum length of ${maxLength} characters`;
    return;
  }

  if (!trimmedCode) {
    outputElement.textContent = "Error: No code to execute";
    return;
  }

  try {
    outputElement.textContent = "";
    const typingText = "Executing...";

    let i = 0;
    const typingInterval = setInterval(() => {
      if (i < typingText.length) {
        outputElement.textContent += typingText[i];
        i++;
      } else {
        clearInterval(typingInterval);

        setTimeout(() => {
          try {
            const result = evaluateSupBud(trimmedCode);
            outputElement.textContent = result.error ? "Error: " + result.error : result.result;
          } catch (error) {
            outputElement.textContent = "Error: " + error.message;
            console.error(error);
          }
        }, 300);
      }
    }, 50);
  } catch (error) {
    outputElement.textContent = "Error: " + error.message;
    console.error(error);
  }
}

const TOKEN_TYPES = {
  KEYWORD: ["sup", "bud", "true", "false", "if", "else", "return"],
  OPERATOR: ["+", "-", "*", "/", "=", "==", "!=", "<", ">", "<=", ">="],
  PUNCTUATION: [",", ";", "(", ")", "{", "}"],
  NUMBER: /^\d+$/,
  IDENTIFIER: /^[a-zA-Z_][a-zA-Z0-9_]*$/,
};

function tokenize(code) {
  const tokens = [];
  let currentToken = "";

  function pushCurrentToken() {
    if (currentToken) {
      let type = "identifier";

      if (TOKEN_TYPES.KEYWORD.includes(currentToken)) {
        type = "keyword";
      } else if (currentToken.match(TOKEN_TYPES.NUMBER)) {
        type = "number";
      }

      tokens.push({
        text: currentToken,
        type: type,
      });

      currentToken = "";
    }
  }

  for (let i = 0; i < code.length; i++) {
    const char = code[i];

    if (i + 1 < code.length) {
      const twoChars = char + code[i + 1];
      if (TOKEN_TYPES.OPERATOR.includes(twoChars)) {
        pushCurrentToken();
        tokens.push({
          text: twoChars,
          type: "operator",
        });
        i++;
        continue;
      }
    }

    if (TOKEN_TYPES.OPERATOR.includes(char) || TOKEN_TYPES.PUNCTUATION.includes(char)) {
      pushCurrentToken();
      tokens.push({
        text: char,
        type: TOKEN_TYPES.OPERATOR.includes(char) ? "operator" : "punctuation",
      });
      continue;
    }

    if (char === " " || char === "\t" || char === "\n") {
      pushCurrentToken();
      tokens.push({
        text: char,
        type: "whitespace",
      });
      continue;
    }

    currentToken += char;
  }

  pushCurrentToken();

  return tokens;
}

document.addEventListener("DOMContentLoaded", () => {
  initWasm();
  setTimeout(setupScrollSync, 1000);
});