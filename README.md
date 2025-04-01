# tmgr - Terminal Task Manager & PA

A simple, efficient command-line task manager and personal assistant built with Go and SQLite. Manage tasks, store quick notes, and get automatic reminders directly in your terminal.

## Features

*   **Task Management:** Add, list, mark as done, and remove tasks.
*   **Note Taking:** Quickly add, list, and remove simple notes.
*   **Flexible Due Dates:** Set due dates for tasks using absolute formats (`YYYY-MM-DD`, `YYYY-MM-DD HH:MM:SS`) or convenient relative durations (`2d`, `3h30m`, `1w`, `6M`).
*   **Automatic Reminders:** Integrated with shell startup (PowerShell, Bash, Zsh) to automatically check for overdue and upcoming tasks.
*   **Local Storage:** Uses a simple SQLite database stored locally in your user's standard configuration directory.
*   **Cross-Platform:** Built with Go, works on Windows, Linux, and macOS.

## Installation (Building from Source)

**Prerequisites:**
*   Go (version 1.18 or later recommended) installed and configured.
*   Git installed.
*   A C compiler (like `gcc` or `clang`) for the SQLite driver (`mattn/go-sqlite3`). Most Linux/macOS systems with development tools have this. On Windows, `gcc` (via MinGW/MSYS2/etc.) might be needed if not already set up.

**Steps:**

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/57Ajay/go-taskmanagerCLI.git
    cd go-taskmanagerCLI
    ```

2.  **Build the executable:**
    *   **Linux/macOS:**
        ```bash
        go build -o tmgr
        ```
    *   **Windows:**
        ```bash
        go build -o tmgr.exe
        ```
    *(This creates the executable `tmgr` or `tmgr.exe` in the current directory)*

3.  **Make the executable accessible (Recommended):**

    *   **Linux/macOS:**
        1.  Choose a directory in your system's `$PATH`. Common choices are `~/bin` (you might need to create it: `mkdir ~/bin`) or `/usr/local/bin` (may require `sudo`).
        2.  Move the `tmgr` executable there:
            ```bash
            # Example for ~/bin
            mv tmgr ~/bin/
            # Or for /usr/local/bin
            sudo mv tmgr /usr/local/bin/
            ```
        3.  Ensure the chosen directory is in your `$PATH`. For `~/bin`, you might need to add `export PATH="$HOME/bin:$PATH"` to your `~/.bashrc`, `~/.zshrc`, or `~/.profile` and restart your shell. `/usr/local/bin` is usually in the PATH by default.
        4.  Verify by opening a *new* terminal tab/window and typing `tmgr version`.

    *   **Windows:**
        1.  Create a directory if needed (e.g., `C:\Users\<YourUsername>\bin` or `C:\Tools`).
        2.  Copy `tmgr.exe` to that directory.
        3.  Add this directory to your Windows **PATH** environment variable (User variables -> Path -> Edit -> New).
        4.  **Important:** Close and reopen PowerShell/CMD for the PATH change to take effect.
        5.  Verify by opening a new PowerShell and typing `tmgr --version`.

## Usage

All commands are accessed via the main `tmgr` executable (or `tmgr.exe` on Windows).

### General Commands

*   **Show Version:**
    ```bash
    tmgr version
    ```
*   **Show Help:**
    ```bash
    tmgr --help
    tmgr [command] --help # Help for a specific command (e.g., tmgr task add --help)
    ```

### Task Management (`tmgr task ...`)

*   **Add a Task:**
    ```bash
    tmgr task add <description> [--due <date_or_duration>]
    ```
    *   `<description>`: The text of the task (required).
    *   `--due` / `-d` (Optional): Sets the due date. Accepts:
        *   Absolute Dates: `YYYY-MM-DD` (e.g., `2025-12-31`)
        *   Absolute DateTimes: `YYYY-MM-DD HH:MM:SS` (e.g., `"2025-04-15 09:30:00"`)
        *   Relative Durations: Combinations of `y`(year), `M`(month), `w`(week), `d`(day), `h`(hour), `m`(minute), `s`(second) (e.g., `2d`, `1w`, `3h30m`, `"1M 2d"`). Also supports standard Go durations like `1.5h`.
    *   **Examples:**
        ```bash
        tmgr task add "Submit project report"
        tmgr task add "Prepare presentation slides" --due 2d
        tmgr task add "Team meeting" -d "2025-05-01 10:00:00"
        tmgr task add "Follow up with client" --due 1w
        ```

*   **List Tasks:**
    ```bash
    tmgr task list
    ```
    *   Displays all tasks with ID, status ([pending]/[DONE]), description, and due date.
    *   Marks overdue tasks with `[OVERDUE!]`.
    *   **Example Output:** (Format may vary slightly based on terminal)
        ```
        Your Tasks:
        -----------
        1. [pending] Submit project report (Due: N/A)
        2. [pending] Prepare presentation slides (Due: Sat, 29 Mar - 16:00)
        3. [DONE]    Team meeting (Due: Tue, 01 May - 10:00)
        4. [pending] Follow up with client (Due: Sat, 05 Apr - 16:00) [OVERDUE!]
        -----------
        ```

*   **Mark Task as Done:**
    ```bash
    tmgr task done <task_id>
    ```
    *   `<task_id>`: The numerical ID of the task to mark as completed.
    *   **Example:** `tmgr task done 1`

*   **Remove a Task:**
    ```bash
    tmgr task remove <task_id>
    ```
    *   `<task_id>`: The numerical ID of the task to permanently remove.
    *   **Example:** `tmgr task remove 3`

*   **Check for Important Tasks:**
    ```bash
    tmgr task check [--silent-if-clear]
    ```
    *   Scans tasks and displays separate lists for "Overdue" and "Upcoming (due within 24 hours)".
    *   `--silent-if-clear` / `-s` (Optional): Suppresses output if no tasks need attention (useful for startup scripts).
    *   **Example:** `tmgr task check` or `tmgr task check -s`

### Note Management (`tmgr note ...`)

*   **Add a Note:**
    ```bash
    tmgr note add <content>
    ```
    *   `<content>`: The text of the note (required).
    *   **Example:** `tmgr note add "Remember API key location"`

*   **List Notes:**
    ```bash
    tmgr note list
    # OR just:
    tmgr note
    ```
    *   Displays all notes with ID, content, and creation timestamp.

*   **Remove a Note:**
    ```bash
    tmgr note remove <note_id>
    ```
    *   `<note_id>`: The numerical ID of the note to remove.
    *   **Example:** `tmgr note remove 1`

## Configuration

### Database Location

`tmgr` stores its data in a SQLite database file located in the standard user configuration directory:
*   **Linux:** `~/.config/TaskManagerCLI/tasks.db` (or your chosen app name)
*   **macOS:** `~/Library/Application Support/TaskManagerCLI/tasks.db`
*   **Windows:** `C:\Users\<YourUsername>\AppData\Roaming\TaskManagerCLI\tasks.db`

*(Path based on `os.UserConfigDir()` in Go)*.

### Shell Startup Integration (Automatic Reminders)

To get automatic reminders when opening a new terminal/shell session:

1.  Ensure the `tmgr` executable is built and accessible via your system PATH (see Installation).
2.  Edit your shell's startup configuration file:
    *   **Bash:** Edit `~/.bashrc`
    *   **Zsh:** Edit `~/.zshrc`
    *   **PowerShell:** Edit the file path given by `$PROFILE` (run `notepad $PROFILE` to open it; create if needed with `New-Item -Path $PROFILE -ItemType File -Force`).
    *   *(Other shells like Fish have different config files, e.g., `~/.config/fish/config.fish`)*
3.  Add the following line to the **end** of the file:
    ```bash
    # For Bash/Zsh/Fish etc. on Linux/macOS
    tmgr task check --silent-if-clear

    # For PowerShell on Windows
    tmgr task check --silent-if-clear
    # OR if not in PATH: & "C:\path\to\your\tmgr.exe" task check --silent-if-clear
    ```
4.  Save the configuration file.
5.  **Important:** Close and reopen your terminal/shell completely for the changes to take effect. Reminders should now appear on startup if needed, and remain silent otherwise.

*(Note for PowerShell: If scripts are disabled, you might need `Set-ExecutionPolicy RemoteSigned -Scope CurrentUser`.)*
*(Note for Linux/macOS: Ensure the check command runs quickly and doesn't produce excessive output if using `-s`)*

## Contributing

Contributions are welcome! Whether it's reporting a bug, suggesting a feature, or submitting a pull request, your input is valued.

### Ways to Contribute

*   **Report Bugs:** If you encounter a bug, please open an issue on GitHub detailing the problem, steps to reproduce it, your operating system, and the version of `tmgr` you are using.
*   **Suggest Enhancements:** Have an idea for a new feature or an improvement to an existing one? Open an issue to discuss it. This allows for feedback before significant development work begins.
*   **Submit Pull Requests:** If you want to contribute code directly, please follow the workflow below.

### Development Setup

1.  **Fork the repository:** Click the "Fork" button on the top right of the GitHub repository page.
2.  **Clone your fork:**
    ```bash
    git clone https://github.com/57Ajay/go-taskmanagerCLI.git
    cd go-taskmanagerCLI
    ```
3.  **Set up prerequisites:** Ensure you have Go installed (see Installation section) and any necessary C compilers for the SQLite driver.
4.  **Build the project:**
    ```bash
    go build -o tmgr # (or tmgr.exe on Windows)
    ```
    You should now be able to run `./tmgr` (or `.\tmgr.exe`) from the project directory.

### Pull Request Workflow

1.  **Create a Branch:** Before making changes, create a new branch off the `main` branch (or the primary development branch):
    ```bash
    git checkout main
    git pull origin main # Ensure you have the latest changes from upstream
    git checkout -b feature/your-feature-name # e.g., feature/add-task-priority
    ```
2.  **Make Changes:** Write your code. Ensure it follows standard Go practices:
    *   Run `go fmt ./...` to format your code.
    *   Run `go vet ./...` to check for suspicious constructs.
    *   Add comments where necessary.
    *   If adding a new command or flag, update this `README.md` accordingly.
3.  **Commit Changes:** Make clear, concise commit messages.
    ```bash
    git add .
    git commit -m "feat: Add feature X" -m "Detailed description of the changes made."
    # Or "fix: Fix bug Y", "docs: Update README", etc.
    ```
4.  **Push Branch:** Push your feature branch to *your* fork:
    ```bash
    git push origin feature/your-feature-name
    ```
5.  **Open a Pull Request (PR):** Go to the original `tmgr` repository on GitHub. You should see a prompt to open a pull request from your recently pushed branch.
    *   Ensure the PR targets the `main` branch of the original repository.
    *   Provide a clear title and description for your PR, explaining the "what" and "why" of your changes. Link to any relevant issues (e.g., "Closes #123").
6.  **Code Review:** Be prepared to discuss your changes and make adjustments based on feedback during the code review process.

### Coding Style

*   Follow standard Go formatting (`go fmt`).
*   Adhere to Go best practices (effective Go, error handling, etc.).
*   Keep changes focused – one feature or bug fix per pull request is ideal.

Thank you for contributing!

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
