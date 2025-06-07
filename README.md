
# 📂 FileFlow AI

**FileFlow AI** is a simple and intelligent file organizer built with Golang.

---

## 💼 How it works
**Instructions:** ¹Insert files into _files/_ (create it if necessary) → ²Run the program _(explained below)_ → ³It will output the result at _result/_, and the remaining empty folders or files that didn't moved right will be placed in _trash/_, so analyze the folder and delete it when necessary. → Organization done!

## 📋 Project Structure

**cmd/** - Entry Point

**internal/** - Program Actions
- **ai/** - AI functions
- **fileutils/** - File handling
- **folderutils/** - Folder creation

_extras_ ---------------------------------

**files/** - Files input

**result/** - Program output

**trash/** - Remaining empty folders | files moved wrong

---

## 💻 Installation

1. **Clone the repository:**
```bash
git clone https://github.com/vMatheusCouto/fileflow-ai.git
cd fileflow-ai
```

2. **Initialize Go modules:**
```bash
go mod tidy
```

3. **Run the program:**

```bash
go run ./cmd
```

---
