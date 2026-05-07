# Git Workflow Rules

## Branches

main → stable code

develop → integration branch

feature/* → feature branches

---

## Rules

- Never push directly to main
- Pull latest changes before starting work
- Create a feature branch for every task
- Open pull requests before merging

---

## Example Workflow

```bash
git checkout develop
git pull
git checkout -b feature/master-api
```