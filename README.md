# Waypoint Plugin Archive

waypoint-plugin-archive is a builder plugin for [Waypoint](https://github.com/hashicorp/waypoint). It allows you to
archive source files or directories.

**The plugin is working as expected for my use case but is still missing some features, please open an issue for any
feedback, issues or missing features..**

# Install

To install the plugin, run the following command:

```bash
make install # Installs the plugin in `${HOME}/.config/waypoint/plugins/`
```

# Configure

```hcl
project = "project"

app "webapp" {
  path = "./webapp"

  build {
    use "archive" {}
  }
}
```
