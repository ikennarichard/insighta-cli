# Insighta CLI

A command-line tool for interacting with the Insighta+ Labs API.

---

## Requirements

- Go 1.22 or higher
- A running instance of the Insighta backend API

---

## Installation

**Step 1 — Clone the repo:**
```bash
git clone https://github.com/ikennarichard/insighta-cli.git
cd insighta-cli
```

**Step 2 — Install globally:**
```bash
go install .
```


Then make sure `$GOPATH/bin` is in your PATH. Add this to your `~/.bashrc` or `~/.zshrc`:
```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

Reload your shell:

```bash
echo $SHELL
 # if you see this /usr/bin/bash  use  this source ~/.bashrc
source ~/.zshrc  
```

**Step 4 — Verify installation:**
```bash
ls $(go env GOPATH)/bin # use the option u see in the folder (insighta-cli or insighta)
insighta-cli --help # or insighta --help
```

## Authentication

**Login:**
```bash
insighta-cli login
```
This opens GitHub in your browser. After authenticating, the CLI saves
your credentials to `~/.insighta/credentials.json` automatically.

**Check who you are:**
```bash
insighta-cli whoami
```

**Logout:**
```bash
insight-cli logout
```

---

## Profiles

**List all profiles:**
```bash
insighta-cli profiles list
```

**Filter profiles:**
```bash
insighta-cli profiles list --gender male
insighta-cli profiles list --country NG
insighta-cli profiles list --age-group adult
insighta-cli profiles list --min-age 25 --max-age 40
```

**Sort profiles:**
```bash
insighta-cli profiles list --sort-by age --order desc
```

**Paginate:**
```bash
insighta-cli profiles list --page 2 --limit 20
```

**Get a single profile:**
```bash
insighta-cli profiles get <id>
```

**Search using natural language:**
```bash
insighta-cli profiles search "young males from nigeria"
```

**Create a profile:**
```bash
insighta-cli profiles create --name "Harriet Tubman"
```

**Export to CSV:**
```bash
insighta-cli profiles export --format csv
insighta-cli profiles export --format csv --gender male --country NG
```
The CSV file is saved to your current working directory.

---

## Pointing to a different backend

By default the CLI talks to `http://localhost:8080`.
To use a different backend pass the `--api-url` flag:

```bash
insighta --api-url https://api.yourserver.com profiles list
```

---

## Credentials

Credentials are stored at:

```bash
~/.insighta/credentials.json   # Mac/Linux
C:\Users<you>.insighta\credentials.json  # Windows
```