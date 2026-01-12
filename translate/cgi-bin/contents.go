bookName = os.Args[0]
dbpath = "/usr/share/treadlate/database/" + bookName + ".sqlite3"
title = os.Exec("sqlite3", dbpath, `
select group_concat(item_interlinear, ' ')
from words
where chapter = 1 AND paragraph = 1 AND sentence = 1 AND item_type = 'h1'
`)
abbreviation = os.Exec("sqlite3", dbpath, `
select value
from info
where key = 'abbreviation'
`)
abbreviation = abbreviation[:len(abbreviation)-1]
writer.Write(`
<!DOCTYPE html>
<html>` + os.ScriptEval("layouts", "head", title) + `
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
<h2>` + title + `- სათარგმნი აპლიკაცია</h2>
    <p>გადათარგმნილი ტექსტი მხოლოდ თქვენს აპლიკაციაში
    რჩება, ჩვენთვის გამოსაგზავნად ჩამოწერეთ თქვენი აპლიკაციის 
    <button onclick="downloadJSON()">ლოკალური საცავის json</button> 
    ფაილი მოცემული წიგნისთვის და იმეილით გამოგვიგზავნეთ. დიდი
    მადლობა წინასწარ.
    თქვენი ნათარგმნი უსასყიდლოდ გახდება ხელმისაწვდომი მთელი 
    მსოფლიოსათვის წასაკითხად და გამოსაყენებლად ნებისმიერი სახით.
    </p>
    <p>თქვენ შეგიძლიათ ჩამოტვირთოთ ფაილი მთელი აპლიკაციისთვის
    საკონტაქტო (contact) გვერდიდან ანდა ჩამოტვირთოთ თითოეული თავისთვის
    ფაილი თავისი გვერდიდან</p>
    <p>	
   	naierchou@proton.me
   	</p>

    <ul>` + os.Exec("sqlite3", dbpath, `
select '<li><a href="/translate/chapter/`+bookName+`/'||
chapter ||
'">' ||
     chapter ||
  		 ' - ' ||
group_concat(item, ' ') ||
'</a></li>'
from words
where paragraph = 1 AND item_type = 'h1'
group by chapter, paragraph, sentence, item_type
`) + `</ul></main><script>
   async function downloadJSON() {
    const storageData = {};
    for (let i = 0; i < localStorage.length; i++) {
        const key = localStorage.key(i);
        if (key.startsWith("` + abbreviation + `_")) {
            storageData[key] = localStorage.getItem(key);
        }
    }

    const jsonString = JSON.stringify(storageData, null, 2);

    const filename = "localStorage-book-` + bookName + `.json"

    await goSaveJSON(jsonString, filename)
}
</script></body></html>`)
