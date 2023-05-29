<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>{{ .Title }}</title>
</head>
<body>
<h1>{{ .Title }}</h1>
<ul>
{{ range $c := .Contents }}
    <li><a href="{{ $c.ListeningPage }}">{{ $c.Date.Format "2006/01/02" }} - {{ $c.Title }}</a></li>
{{ end }}
</ul>
<p><a href="{{ .Feed }}">Feeds</a></p>
</body>
</html>
