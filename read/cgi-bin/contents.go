bookName = os.Args[0]
dbpath = "/usr/share/treadlate/database/" + bookName + ".sqlite3"
title = os.Exec("sqlite3", dbpath, `
select group_concat(item_interlinear, ' ')
from words
where chapter = 1 AND paragraph = 1 AND sentence = 1 AND item_type = 'h1'
`)
writer.Write(`<!DOCTYPE html>
<html>
` + os.ScriptEval("layouts", "head", title) + `
<body>
<nav>` + `
<ul>
<li><a href="/translate/contents/` + bookName + `">Translate</a></li>
<li><a href="/read/contents/` + bookName + `">Read</a></li>
<li><a href="/layouts/contact/` + bookName + `">Contact</a>
</ul>
` + `
  </nav>
  <main>
<h2>` + title + `- წასაკითხი</h2>
     <ul>` + os.Exec("sqlite3", dbpath, `
select '<li><a href="/read/chapter/`+bookName+`/'||
chapter ||
'">' ||
     chapter ||
  		 ' - ' ||
group_concat(item, ' ') ||
'</a></li>'
from words
where paragraph = 1 AND item_type = 'h1'
group by chapter, paragraph, sentence, item_type
`) + `</ul></main></body></html>`)
