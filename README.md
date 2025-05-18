# 🪄 qik - AI-Powered Text Assistant 🚀

**Your intelligent command-line companion for text refinement, explanation, and Q&A, powered by Google Gemini.**


`qik` is a versatile CLI tool designed to streamline your text-based tasks. Whether you need to polish your writing, understand complex text, or get quick answers to your questions, `qik` leverages the power of Google's Gemini AI to assist you directly from your terminal.

---

## ✨ Features

*   📝 **Fix Text (`qik fix`)**:
    *   Corrects spelling and grammar.
    *   Improves text flow and clarity.
    *   Adjusts text to various moods/tones (e.g., professional, casual, funny).
    *   Supports multiple languages (default: Norwegian, configurable).
    *   Copies the fixed text directly to your clipboard.
*   💡 **Explain Text (`qik explain`)**:
    *   Provides simple, concise explanations of complex text.
    *   Automatically attempts to explain in the language of the input text.
    *   Outputs explanation to the terminal.
    *   Option to copy explanation to clipboard.
*   ❓ **Answer Questions (`qik answer`)**:
    *   Get clear answers to your questions.
    *   Supports different languages and moods for the answer.
    *   Outputs answer to the terminal.
    *   Option to copy answer to clipboard.
*   ⚙️ **Configurable**:
    *   Uses a simple YAML configuration file (`~/.config/qik/config.yaml`).
    *   Customize default language, editor, AI model, prompts, and moods.
    *   Secure API key handling (environment variable or `pass` recommended).
*   📋 **Clipboard Integration**: Easily paste processed text.
*   ✏️ **Editor Integration**: Uses your preferred command-line editor (`nvim`, `vim`, `nano`, etc.) for text input.
*   💨 **Fast & Efficient**: Leverages models like `gemini-1.5-flash-latest` for quick responses.

---

## 📦 Installation

### One-Liner (Recommended for Linux/macOS)

You can install the latest release of `qik` with this command:

```bash
curl -sfL https://raw.githubusercontent.com/Andr0id88/qik/main/install.sh | sh -s
```

This script will:
1. Detect your OS and architecture.
2. Download the appropriate binary from the latest GitHub release.
3. Make it executable.
4. Attempt to install it into a directory in your PATH (e.g., $HOME/.local/bin, /usr/local/bin). It will ask for sudo if needed for system-wide directories.

### Manual Installation

1. Go to the [Releases page](https://github.com/Andr0id88/qik/releases)
2. Download the appropriate archive (.tar.gz or .zip) for your operating system and architecture.
3. Extract the qik binary from the archive.
4. Move the qik binary to a directory in your system's PATH (e.g., /usr/local/bin, $HOME/.local/bin, or $HOME/bin).
5. Make sure the binary is executable: chmod +x /path/to/qik.

### From Source (for Developers)

1.Clone the repository:
```bash
git clone https://github.com/Andr0id88/qik.git
cd qik
```
2. Build the binary
```bash
go build -o qik .
```
3. Move the qik binary to your desired location in PATH.


## 🚀 Usage

qik uses your configured editor (default: nvim) to open a temporary file where you can type or paste your text. After you save and close the editor, qik processes the text.

### 🛠️ Fixing Text: `qik fix`
Correct spelling, grammar, improve flow, and adjust tone.
```bash
qik fix
```

- Specify Language:
```bash
qik fix -l English  # or --language English
qik fix -e          # Shorthand for English fix
```

- Specify Mood/Tone:
```bash
qik fix -m professional # or --mood professional
qik fix -m funny -l Norwegian
```
*Run qik list-moods to see available moods.*

- Specify AI Prompt Template (Advanced):

```bash
qik fix -p english_fix_only # Uses the 'english_fix_only' prompt from config
```

### 🧐 Explaining Text: `qik explain`

Get a simple explanation of a piece of text.

```bash
qik explain
```

- Specify Language for Explanation:

```bash
qik explain -l English # Get explanation in English
```
(If no language is specified, qik attempts to match the input text's language.)

- Copy Explanation to Clipboard:

```bash
qik explain -c # or --copy
```

### 💬 Answering Questions: `qik answer`
Ask a question and get an AI-generated answer.

```bash
qik answer
```
- Specify Language for Answer:

```bash
qik answer -l English

```
- Specify Mood/Tone for Answer:

```bash
qik answer -m concise
```
- Copy Answer to Clipboard:

```bash
qik answer -c # or --copy
```
### ℹ️ General Options

* -v, --verbose: Enable verbose output for more details.
* --config /path/to/config.yaml: Specify a custom configuration file.
* --help: Show help for qik or any subcommand.

### 📋 Listing Options

* qik list-models: Shows available Gemini models with descriptions.
* qik list-moods: Shows moods defined in your configuration.

## ⚙️ Configuration

qik looks for a configuration file in the following order:

1. Path specified by the --config flag.
2. $HOME/.config/qik/config.yaml (on Linux/macOS) or platform equivalent.
3. ./config.yaml (in the current directory).

If no config file is found, qik will attempt to create a default one at $HOME/.config/qik/config.yaml.

You can customize:
* defaultLanguage: e.g., "Norwegian", "English"
* editor: e.g., "nvim", "vim", "nano", "code --wait"
* geminiModel: e.g., "gemini-1.5-flash-latest"
* defaultMood: e.g., "neutral", "professional"
* prompts: Customize the instructions given to the AI for fix, explain, and answer tasks.
* moods: Define custom moods with their descriptions and AI instructions.

See the config.example.yaml in this repository for a full example and all available options.

**API Key:**
It is **highly recommended** to set your Gemini API key using:

1. The pass password manager: pass insert gemini_api_key
2. An environment variable: export GEMINI_API_KEY="YOUR_API_KEY"

You can also set geminiApiKey in the config file, but this is less secure.

## 🤝 Contributing (Optional)

Contributions are welcome! If you have ideas for improvements or find a bug, please feel free to:
1. Fork the repository.
2. Create a new feature branch (git checkout -b feature/AmazingFeature).
3. Make your changes.
4. Commit your changes (git commit -m 'Add some AmazingFeature').
5. Push to the branch (git push origin feature/AmazingFeature).
6. Open a Pull Request.

## 📄 License
This project is licensed under the [MIT License](https://github.com/Andr0id88/qik/).


































```bash
```
```bash
```







