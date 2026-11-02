# Catalog fixture

## Deploy

### Existing resource conflict

**Symptoms:**

* A deploy fails with:

  ```text
  helm.go:75: [debug] existing resource conflict
  ```

**Root cause:**

Helm refuses to adopt an untracked object.

**How to check:**

1. Read the Helm release annotation.

**How to fix:**

1. Adopt the object.

### Other case

**Symptoms:**

* Unrelated second symptom.

**Root cause:**

Something else.
