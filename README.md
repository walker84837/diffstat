# diffstat 

> Compare how different two Git branches are.

[![Build Source Code](https://github.com/walker84837/diffstat/actions/workflows/build.yml/badge.svg)](https://github.com/walker84837/diffstat/actions/workflows/build.yml)
[![License: MPL-2.0](https://img.shields.io/badge/License-MPL--2.0-blue.svg)](LICENSE)

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [Contributing](#contributing)
  - [Roadmap](#roadmap)
- [License](#license)

## Installation

Start by cloning the repository and building the project:

```bash
git clone https://github.com/walker84837/diffstat.git
cd diffstat
go build
```

Once built, run the tool by specifying the two branches you want to compare:

```bash
./diffstat <branch1> <branch2>
```

`branch1` is usually the base (or main) branch and `branch2` is usually the feature branch.

## Usage

Simply invoke the command above to see a summary of how your Git branches differ. It's that easy!

## Contributing

Contributions are more than welcome. Whether it's a bug fix, feature enhancement, or just ideas, feel free to open an issue or submit a pull request.

### Roadmap

- [ ] **Dynamic Text Colors:** Update text color based on the difference magnitude 🎨
- [ ] **Performance Boost:** Optimize for large repositories ⚡

Your feedback and contributions help shape the project. Check out our [issues](https://github.com/walker84837/diffstat/issues) and feel free to go through the source code for more details.
