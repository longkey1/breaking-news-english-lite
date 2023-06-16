<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>{{ .Title }}</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.2.1/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-iYQeCzEYFbKjA/T2uDLTpkwGzCiq6soy8tYaI1GyVh/UjpbCx/TYkiZhlZB6+fzT" crossorigin="anonymous">
  </head>
  <body>
    <div class="container">
      <div class="mt-5">
        <h1>{{ .Title }}</h1>
        <ul class="list-group list-group-flush text-left">
        {{ range $c := .Contents }}
            <li class="list-group-item"><a href="{{ $c.ListeningPage }}">{{ $c.Date.Format "2006/01/02" }} - {{ $c.Title }}</a></li>
        {{ end }}
        </ul>
      </div>
      <div>
        <a href="{{ .Feed }}">Feed</a>
      </div>
      <div class="mt-5">
        <p class="text-end">Updated at {{ .UpdatedAt.Format "2006-01-02 15:04:05" }}</p>
      </div>
    </div>
    <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.2.1/dist/js/bootstrap.bundle.min.js" integrity="sha384-u1OknCvxWvY5kfmNBILK2hRnQC3Pr17a+RTT6rIHI7NnikvbZlHgTPOOmMi466C8" crossorigin="anonymous"></script>
  </body>
</html>
