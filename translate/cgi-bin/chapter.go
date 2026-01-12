bookName = os.Args[0]
chapter = os.Args[1]
dbpath = "/usr/share/treadlate/database/" + bookName + ".sqlite3"
dictpath = "/usr/share/treadlate/database/" + bookName + ".SQLite3"
bookTitle = os.Exec("sqlite3", dbpath, `
select group_concat(item, ' ')
from words
where chapter = 1 AND paragraph = 1 AND sentence = 1 AND item_type = 'h1'
`)

abbreviation = os.Exec("sqlite3", dbpath, `
select value
from info
where key = 'abbreviation'
`)
abbreviation = abbreviation[:len(abbreviation)-1]
title = os.Exec("sqlite3", dbpath, `
select  group_concat(item, ' ') 
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
<a href="/translate/contents/` + bookName + `">` + bookTitle + `</a><h2>სათარგმნი ვებ აპლიკაცია - თავი ` + chapter + `</h2>
<p>აპლიკაციაში გადათარგმნილი ტექსტი მხოლოდ თქვენს ბრაუზერში
    რჩება, ჩვენთვის გამოსაგზავნად ჩამოწერეთ თქვენი ბრაუზერის 
    <button onclick="downloadJSON()">ლოკალური საცავის json</button> 
    ფაილი მოცემული თავისათვის და იმეილით გამოგვიგზავნეთ.
    დიდი მადლობა წინასწარ.
    თქვენი ნათარგმნი უსასყიდლოდ გახდება ხელმისაწვდომი მთელი 
    მსოფლიოსათვის წასაკითხად და გამოსაყენებლად ნებისმიერი სახით.
    </p>
    <p>თქვენ შეგიძლიათ ჩამოტვირთოთ ფაილი თითოეული წიგნისთვის
    შესაბამისი წიგნის თავების ჩამონათვალის "content" გვერდიდან
    ანდა ჩამოტვირთოთ მთელი საიტის მასშტაბით ფაილი საკონტაქტო გვერდიდან</p>
<a href="mailto:naierchou@proton.me?subject=` + bookName + "_თავი-" + chapter + `%20ნათარგმნი%20ლოკალური%20საცავი">
   naierchou@proton.me</a>
  <p>[ CTRL + p ] - სიტყვის ლექსიკონიდან ნათარგმნის ჩვენება ვებგვერდის ქვედა ნაწილში<br>
     [ CTRL + ENTER ] - ჩაწერილი სიტყვისათვის სასვენი ნიშნების დამატება<br>
     [ TAB ] - შემდეგ უჯრაზე გადასვლა<br>
     [ SHIFT + TAB ] - წინა უჯრაზე გადასვლა</p>
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
  SELECT 
    lowercase_item,
    GROUP_CONCAT('<option value="' || item_interlinear || '">', '') AS options
  FROM (
    SELECT DISTINCT lowercase_item, item_interlinear
    FROM words
    WHERE item_interlinear != ''
  )
  WHERE lowercase_item IN (SELECT lowercase_item FROM initial_items)
  GROUP BY lowercase_item
),
annotated_words AS (
  SELECT
    w.*,
    pw.end_tags AS prev_end_tags,
    d.options AS datalist_options
  FROM ordered_words w
  LEFT JOIN ordered_words pw ON pw.rn = w.rn - 1
  LEFT JOIN datalist d ON d.lowercase_item = w.lowercase_item
)
SELECT
  GROUP_CONCAT(
  REPLACE(w.start_tags, '<a', '<span') ||
  CASE
    WHEN w.start_tags = '' AND (w.prev_end_tags IS NULL OR w.prev_end_tags = '') THEN '<br/>
'
    ELSE ''
  END ||
  '<span i="' || w.lowercase_item || '">' || w.start_punctuation || w.item || w.end_punctuation || '' ||
  CASE 
    WHEN w.word_per_paragraph = 0 OR instr(w.start_tags, ' t ') THEN ''
    ELSE
      ' <input list="' || w.chapter || '_' || w.paragraph || '_' || w.item_per_paragraph || 
      '" s=''' || w.start_punctuation || '''  e=''' || w.end_punctuation || ''' type="text" value="' || 
      w.start_punctuation_interlinear || w.item_interlinear || w.end_punctuation_interlinear || '">' ||
      '<datalist id="' || w.chapter || '_' || w.paragraph || '_' || w.item_per_paragraph || '">' ||
      COALESCE(w.datalist_options, '') ||
      '</datalist>'
    END ||'</span>' ||
  REPLACE(w.end_tags, '</a>', '</span>'), '')
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
function normalizePunctuation(text, isPrefix) {
  const replacements = isPrefix
    ? { '"': '„', '«': '„', '“': '„' }
    : { '"': '“', '»': '“', '”': '“' };

  return Object.entries(replacements).reduce(
    (acc, [target, replacement]) => acc.replace(target, replacement),
    text
  );
}
document.querySelectorAll('span[i]').forEach(span => {
  const input = span.querySelector('input');
  const key = span.getAttribute('i');
    // Keydown event on the input inside the span
  if (input) {
  const listId = input?.getAttribute('list');
let val = localStorage.getItem(listId);
if (val){
input.value = val
}
    input.addEventListener('keydown', event => {

const isAltOrCtrl = event.altKey || event.ctrlKey;
  const isArrow = event.code === 'ArrowUp' || event.code === 'ArrowDown';

  if (isAltOrCtrl && isArrow) {
    event.preventDefault();

    const inputs = Array.from(document.querySelectorAll('input'))
      .filter(el => el.offsetParent !== null && !el.disabled);

    const current = document.activeElement;
    const index = inputs.indexOf(current);
    if (index === -1) return;

    const direction = event.code === 'ArrowDown' ? 1 : -1;
    const nextIndex = index + direction;

    if (nextIndex >= 0 && nextIndex < inputs.length) {
      inputs[nextIndex].focus();
    }
  }
      if (event.ctrlKey && (event.key === 'p' ||
        event.key === 'პ' )) {
        event.preventDefault();
        event.preventDefault();
        console.log('Ctrl + E was pressed!');

        const items = lemmap[key] || [];

        const html = items
          .filter(item => dictionary[item])
          .map(item => ` + "`" + `<b>${item}</b> ${dictionary[item]}` + "`" + `)
          .join('<br/>');

        document.getElementById('definition').innerHTML = html;
      }
      if (event.ctrlKey && event.key === 'Enter') {
        let prefix = input.getAttribute('s');
        let suffix = input.getAttribute('e');
        prefix = normalizePunctuation(prefix, true);
        suffix = normalizePunctuation(suffix, false);

        input.value = prefix + input.value.trim() + suffix;

localStorage.setItem("` + abbreviation + `_" + listId, input.value);
      }
    });
    input.addEventListener('blur', () => {
      const value = input.value;
localStorage.setItem("` + abbreviation + `_" + listId, value);
    });
    input.addEventListener('focus', () => {
      input.scrollIntoView({ behavior: 'smooth', block: 'center', inline: 'center' });
      input.focus(); // optional, if it's focusable
    });
  }
});
   function downloadJSON() {
    	const storageData = {};
    	for (let i = 0; i < localStorage.length; i++) {
        	const key = localStorage.key(i);
          if (key.startsWith("` + abbreviation + "_" + chapter + `_")){
        		storageData[key] = localStorage.getItem(key);
          }
    	}

    	const jsonString = JSON.stringify(storageData, null, 2); // 2 spaces for indentation

    	const filename = "localStorage-` + bookName + "-" + chapter + `.json"

      goSaveJSON(jsonString, filename)
    }

</script></body></html>`)
