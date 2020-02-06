<html>
<head>
<title>{{.PageTitle}}: {{.Path}}</title>
<style type="text/css">
table th, table td {
  padding-left: 30px;
  padding-right: 30px;
}
</style>
</head>
<body>
  <h1>{{.PageTitle}}</h1>
  {{if ne "/" .Path}}
  <a href="..">Top</a>
  {{end}}
  <table>
    <tr>
      <th>size</th>
      <th>last modified</th>
      <th>name</th>
    </tr>
    {{range .Items}}
    <tr>
      <td>{{.Size}}</td>
      <td>{{.ModTime}}</td>
      <td><a href="{{.URL}}">{{.Name}}</a></td>
    </tr>
    {{end}}
  </table>
</body>
</html>
