# Waypoint Plugin Archive [**WIP**]

waypoint-plugin-archive is a builder plugin for [Waypoint](https://github.com/hashicorp/waypoint). 
It allows you to archive source files or directories.

**This plugin is still work in progress, please open an issue for any feedback or issues.**

# Install
To install the plugin, run the following command:

````bash
make install # Installs the plugin in `${HOME}/.config/waypoint/plugins/`
````

# Configure
```hcl
project = "project"

app "webapp" {
  path = "./webapp"

  build {
    use "archive" {
      sources = ["src/", "public/", "package.json"] # Sources are relative to /path/to/project/webapp/
      output_name = "webapp.zip"
      overwrite_existing = true
    }
  }
}
```

