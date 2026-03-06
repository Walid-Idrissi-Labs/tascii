# tascii

A fast, minimal task manager for the terminal.

```
  ID    STATUS  PRIO   TITLE                                   TAGS              DUE
  ──────────────────────────────────────────────────────────────────────────────────
  1     ○       !!!    Deploy to production                    #work             Dec 25 OVERDUE
  2     ◐       !!     Write documentation                     #work #docs       Dec 30
  3     ✓       !      Update dependencies                     —                 —
  ──────────────────────────────────────────────────────────────────────────────────
  3 task(s)
```

## Install

### Homebrew (macOS / Linux)
```sh
brew tap walid-idrissi-labs/tascii
brew install tascii
```

### apt (Debian / Ubuntu)
```sh
# Add the repo and install
curl -sL https://github.com/walid-idrissi-labs/tascii/releases/latest/download/tascii_linux_amd64.tar.gz | tar xz
sudo mv tascii /usr/local/bin/
```

### From source
```sh
go install github.com/walid-idrissi-labs/tascii@latest
```

## Usage

```sh
# Add tasks
tascii add "Fix the login bug"
tascii add "Deploy to production" --priority !!! --due 2024-12-25 --tag work
tascii add "Buy groceries" -p ! -d 2024-12-20 -t personal -n "Check the fridge first"

# List tasks
tascii list                        # all tasks
tascii list --tag work             # filter by tag
tascii list --sort priority        # sort by priority
tascii list --filter in-progress   # only in-progress tasks

# View & today
tascii view 3                      # full detail for task #3
tascii today                       # due today + overdue

# Update status
tascii start 3                     # mark as in-progress  ◐
tascii done 3                      # mark as done         ✓

# Edit
tascii edit 3 --title "New title"
tascii edit 3 --priority !!! --due 2025-01-10
tascii edit 3 --tag work --tag urgent
tascii edit 3 --clear-due

# Delete & cleanup
tascii delete 3                    # with confirmation
tascii delete 3 --force            # skip confirmation
tascii clear                       # remove all completed tasks

# Info
tascii summary                     # compact stats overview
tascii --version
tascii --help
tascii [command] --help
```

## Priority levels

| Symbol | Level  | Meaning        |
|--------|--------|----------------|
| `!`    | Low    | Nice to have   |
| `!!`   | Medium | Should do      |
| `!!!`  | High   | Must do        |

## Shell startup tip

Add to your `.zshrc` or `.bashrc` to see a summary every time you open a terminal:

```sh
tascii summary
```

## Data storage

Tasks are stored at `~/.local/share/tascii/tasks.json` as human-readable JSON.
You can back it up, sync it with Dropbox/iCloud, or edit it manually.

## License

MIT
