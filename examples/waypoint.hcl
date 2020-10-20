project = "examples"

app "hello" {
  path = "./hello"

  build {
    use "archive" {
      sources = ["."]
      output_name = "hello.zip"
      overwrite_existing = true
    }
  }

  deploy {
    use "exec" {}
  }
}

app "ignore" {
  path = "./ignore"

  build {
    use "archive" {
      sources = ["."]
      output_name = "ignore.zip"
      overwrite_existing = true
      ignore = ["ignore", "ignore.txt"]
    }
  }

  deploy {
    use "exec" {}
  }
}

app "collapsed-top-level-folder" {
  path = "./collapsed-top-level-folder"

  build {
    use "archive" {
      sources = ["."]
      output_name = "collapsed.zip"
      overwrite_existing = true
      collapse_top_level_folder = true
    }
  }

  deploy {
    use "exec" {}
  }
}
