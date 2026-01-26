bookName = os.Args[0]
chapter = os.Args[1]
dbpath = "/usr/share/treadlate/database/" + bookName + ".sqlite3"
dictpath = "/usr/share/treadlate/database/" + bookName + ".SQLite3"
bookTitle = os.Exec("sqlite3", dbpath, `
select group_concat(item, ' ')
from words
where chapter = 1 AND paragraph = 1 AND sentence = 1 AND item_type = 'h1'
`)

title = os.Exec("sqlite3", dbpath, `
select group_concat(item, ' ')
from words
where chapter = `+chapter+` AND paragraph = 1  AND item_type = 'h1'
group by paragraph, sentence, item_type
`)

writer.Write(`<!DOCTYPE html>
<html>
` + os.ScriptEval("layouts", "head", "contact") + `
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
<a href="/read/contents/` + bookName + `">` + bookTitle + `</a>
<h2>წასაკითხი ვებ აპლიკაცია - თავი ` + chapter + `</h2>
<div id="content">` + os.Exec("sqlite3", dbpath, `
WITH initial_items AS (
  SELECT DISTINCT lowercase_item
  FROM words
  WHERE chapter = `+chapter+`
),
ordered_words AS (
  SELECT
    *,
    ROW_NUMBER() OVER (ORDER BY paragraph, item_per_paragraph) AS rn
  FROM words
  WHERE chapter = `+chapter+`
),
datalist AS (
  SELECT lowercase_item
  FROM (
    SELECT DISTINCT lowercase_item, item_interlinear
    FROM words
  )
  WHERE lowercase_item IN (SELECT lowercase_item FROM initial_items)
  GROUP BY lowercase_item
),
annotated_words AS (
  SELECT
    w.*,
    pw.end_tags AS prev_end_tags
  FROM ordered_words w
  LEFT JOIN ordered_words pw ON pw.rn = w.rn - 1
  LEFT JOIN datalist d ON d.lowercase_item = w.lowercase_item
)

SELECT
  w.start_tags ||
  CASE
    WHEN w.start_tags = '' AND (w.prev_end_tags IS NULL OR w.prev_end_tags = '') THEN ' '
    ELSE ''
  END ||
  '<span i="' || w.lowercase_item || '">' ||
  	w.start_punctuation ||
w.item ||
w.end_punctuation || 
  CASE 
    WHEN w.word_per_paragraph = 0 THEN ''
    ELSE
  '<sup class="inter">' ||
  w.start_punctuation_interlinear ||
  w.item_interlinear ||
  w.end_punctuation_interlinear||'</sup>'
  END ||
  '</span>' ||
  w.end_tags AS html_output
FROM annotated_words w
ORDER BY w.paragraph, w.item_per_paragraph;
`) + `</div><div id="definition"></div></main><script>
` + os.Exec("sqlite3", dbpath, `
ATTACH DATABASE '`+dictpath+`' AS dictdb;

WITH chapter_terms AS (
  SELECT DISTINCT lowercase_item AS term
  FROM words
  WHERE chapter = `+chapter+` AND word_per_paragraph != 0
),
term_lemmas AS (
  SELECT term, lemma
  FROM lemmas
  WHERE term IN (SELECT term FROM chapter_terms)
),
lemmap_json AS (
  SELECT json_group_object(term, json(lemmas)) AS lemmap
  FROM (
    SELECT term, json_group_array(lemma) AS lemmas
    FROM term_lemmas
    GROUP BY term
  )
),
first_defs AS (
  SELECT
    re.related AS lemma,
    CASE WHEN dict.topic != re.related THEN dict.topic || ' ' || dict.definition
         ELSE dict.definition
    END AS definition
  FROM dictdb.dictionary dict
  JOIN dictdb.relationships re ON dict.topic = re.topic
  WHERE re.related IN (SELECT DISTINCT lemma FROM term_lemmas)
),
found_lemmas AS (
  SELECT DISTINCT lemma FROM first_defs
),
missing_lemmas AS (
  SELECT DISTINCT lemma
  FROM term_lemmas
  WHERE lemma NOT IN (SELECT lemma FROM found_lemmas)
),
fallback_defs AS (
  SELECT
    ae.related AS lemma,
    IFNULL(dict.topic || ' ' || dict.definition, ae.topic) AS definition
  FROM dictdb.ascetic_experiences ae
  LEFT JOIN dictdb.dictionary dict ON dict.topic = ae.topic
  WHERE ae.related IN (SELECT lemma FROM missing_lemmas)
),
combined_defs AS (
  SELECT * FROM first_defs
  UNION ALL
  SELECT * FROM fallback_defs
),
dictionary_json AS (
  SELECT json_group_object(lemma, definition) AS dictionary
  FROM combined_defs
)
SELECT  (SELECT 'const lemmap = ' ||lemmap||';' FROM lemmap_json) || '
' || (SELECT 'const dictionary = ' ||dictionary||';' FROM dictionary_json);
`) + `
document.querySelectorAll('span[i]').forEach(span => {

    span.addEventListener('click', event => {
key = event.target.getAttribute('i');
        const items = lemmap[key] || [];

        const html = items
          .filter(item => dictionary[item])
          .map(item => ` + "`" + `<b>${item}</b> ${dictionary[item]}` + "`" + `) 
          .join('<br/>');

        document.getElementById('definition').innerHTML = html;
    });
});
</script></body></html>`)
