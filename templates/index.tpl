<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>Breaking News English Lite</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.2.1/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-iYQeCzEYFbKjA/T2uDLTpkwGzCiq6soy8tYaI1GyVh/UjpbCx/TYkiZhlZB6+fzT" crossorigin="anonymous">
  </head>
  <body>
    <div class="container">
      <div class="mt-5">
        <h1>Breaking News English Lite</h1>
        <p class="lead">This page provides simple index and customized feeds of <a href="https://breakingnewsenglish.com/">Breaking News English</a></p>
        <ul class="list-group list-group-flush text-left">
          <li class="list-group-item"><a href="level0.html">Level0</a></li>
          <li class="list-group-item"><a href="level1.html">Level1</a></li>
          <li class="list-group-item"><a href="level2.html">Level2</a></li>
          <li class="list-group-item"><a href="level3.html">Level3</a></li>
          <li class="list-group-item"><a href="level4.html">Level4</a></li>
          <li class="list-group-item"><a href="level5.html">Level5</a></li>
          <li class="list-group-item"><a href="level6.html">Level6</a></li>
        </ul>
      </div>
      <div class="mt-5">
        <p class="text-end">Updated at {{ .UpdatedAt.Format "2006-01-02 15:04:05" }}</p>
      </div>
    </div>
    <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.2.1/dist/js/bootstrap.bundle.min.js" integrity="sha384-u1OknCvxWvY5kfmNBILK2hRnQC3Pr17a+RTT6rIHI7NnikvbZlHgTPOOmMi466C8" crossorigin="anonymous"></script>
  </body>
</html>
