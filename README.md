# tascii

A fast, minimal task manager for the terminal built in GO.

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
tascii list                        
tascii list --tag work             
tascii list --sort priority        
tascii list --filter in-progress   

# View & today
tascii view 3                      
tascii today                       

# Update status
tascii start 3                     
tascii done 3                      

# Edit
tascii edit 3 --title "New title"
tascii edit 3 --priority !!! --due 2025-01-10
tascii edit 3 --tag work --tag urgent
tascii edit 3 --clear-due

# Delete & cleanup
tascii delete 3                    
tascii delete 3 --force            
tascii clear                       

# Info
tascii summary                     
tascii --version
tascii --help
tascii [command] --help
```

## Priority levels

- `!` Low Priority
- `!!` Medium Priority
- `!!!` High Priority


## Shell startup option

Add to your `.zshrc` or `.bashrc` to see a summary every time you open a terminal:

```sh
tascii summary
```

## Data storage

Tasks are stored at `~/.local/share/tascii/tasks.json` as human-readable JSON.
You can back it up, sync it with Dropbox/iCloud, or edit it manually.

## License

MIT
