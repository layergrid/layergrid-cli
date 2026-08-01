# CI Integration

LayerGrid is a local static scanner. CI jobs install the CLI, run `layergrid scan`, and keep results inside the workflow unless you explicitly upload SARIF to your code scanning provider.

## GitHub Actions

Use the published composite action:

```yaml
name: LayerGrid

on:
  pull_request:
  push:
    branches: [main]

jobs:
  scan:
    runs-on: ubuntu-latest
    permissions:
      contents: read
      pull-requests: write
    steps:
      - uses: actions/checkout@v4
      - uses: layergrid/scan-action@v1
        with:
          fail-on: high
          comment: true
```

For SARIF upload:

```yaml
permissions:
  contents: read
  pull-requests: write
  security-events: write

steps:
  - uses: actions/checkout@v4
  - uses: layergrid/scan-action@v1
    with:
      fail-on: high
      comment: true
      sarif-upload: true
```

Useful inputs:

```yaml
- uses: layergrid/scan-action@v1
  with:
    path: ./agents
    config: .layergrid.yaml
    fail-on: critical
    comment: false
```

## GitLab CI

Install the release binary with the curl installer:

```yaml
layergrid:
  image: alpine:3.20
  before_script:
    - apk add --no-cache bash curl ca-certificates
    - curl -sSL https://layergrid.github.io/layergrid-cli/install.sh | bash
    - export PATH="$HOME/.local/bin:/usr/local/bin:$PATH"
  script:
    - layergrid scan . --fail-on high
```

For JSON artifacts:

```yaml
layergrid:
  image: alpine:3.20
  before_script:
    - apk add --no-cache bash curl ca-certificates
    - curl -sSL https://layergrid.github.io/layergrid-cli/install.sh | bash
    - export PATH="$HOME/.local/bin:/usr/local/bin:$PATH"
  script:
    - layergrid scan . --format json --output layergrid.json --fail-on high
  artifacts:
    when: always
    paths:
      - layergrid.json
```

## Jenkins

Install once on the agent image, or install during the pipeline:

```groovy
pipeline {
  agent any

  stages {
    stage('LayerGrid') {
      steps {
        sh '''
          curl -sSL https://layergrid.github.io/layergrid-cli/install.sh | bash
          export PATH="$HOME/.local/bin:/usr/local/bin:$PATH"
          layergrid scan . --fail-on high
        '''
      }
    }
  }
}
```

To archive a report:

```groovy
stage('LayerGrid') {
  steps {
    sh '''
      curl -sSL https://layergrid.github.io/layergrid-cli/install.sh | bash
      export PATH="$HOME/.local/bin:/usr/local/bin:$PATH"
      layergrid scan . --format markdown --output layergrid.md --fail-on never
      layergrid scan . --fail-on high
    '''
    archiveArtifacts artifacts: 'layergrid.md', allowEmptyArchive: false
  }
}
```

## Pre-commit Hook

Use LayerGrid as a local guard before pushing:

```sh
cat > .git/hooks/pre-commit <<'SH'
#!/usr/bin/env bash
set -euo pipefail
layergrid scan . --fail-on critical
SH
chmod +x .git/hooks/pre-commit
```

Keep this threshold conservative. Developers should not need to fight medium-severity exploratory findings on every commit.

## Baseline Drift

For mature repos, commit a LayerGrid baseline and fail only when trust boundaries widen:

```sh
layergrid baseline save
git add .layergrid/baseline.json
```

Then run this in CI:

```yaml
- uses: layergrid/scan-action@v1
  with:
    fail-on: never
    comment: true

- name: Compare LayerGrid baseline
  run: layergrid baseline compare --fail-on scope-widening
```

This pattern lets existing findings remain visible while blocking new tool scope expansion.
