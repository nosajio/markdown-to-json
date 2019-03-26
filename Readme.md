# Markdown To JSON
Download a repository of .md files, parse them, and output them as a
structured JSON string.

## API

### mdtojson.ProcessRepo(repoURL string, repoDIR string) (json string, err error)
Download a repo containing markdown files to the specified location. In order
to be able to parse them properly, filenames must be formatted using only
letters and numbers, with dashes instead of spaces, and with the publish date
appended to the file in the format `YYYY-MM-DD`.


For example:
```
Will work:
foo-bar-2019-03-21.md
my-awesome-post-2018-12-03.md
my-awesome-post-2-2018-12-04.md

Won't work:
my_file.md
fizz buzz-2019-02-09.md
my-post-2019-2-9.md
```