project = "examples"

app "hello" {
  path = "./hello"

  build {
    use "archive" {}
  }

  deploy {
    use "exec" {}
  }
}

app "ignore" {
  path = "./ignore"

  build {
    use "archive" {
      ignore = ["ignore", "ignore.txt"]
    }
  }

  deploy {
    use "exec" {}
  }
}

app "include-top-level-directory" {
  path = "./include-top-level-directory"

  build {
    use "archive" {
      include_top_level_directory = true
    }
  }

  deploy {
    use "exec" {}
  }
}
