bookName = os.Args[0]

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
  	<h2>საკონტაქტო ინფორმაცია</h2>
<p>გადათარგმნილი ტექსტი მხოლოდ თქვენს აპლიკაციაში
რჩება, ჩვენთვის გამოსაგზავნად ჩამოწერეთ თქვენი აპლიკაციის
<button onclick="downloadJSON()">ლოკალური საცავის json</button> 
ფაილი და იმეილით გამოგვიგზავნეთ. დიდი მადლობა წინასწარ.
თქვენი ნათარგმნი უსასყიდლოდ გახდება ხელმისაწვდომი მთელი 
მსოფლიოსათვის წასაკითხად და გამოსაყენებლად ნებისმიერი სახით.
</p>
<p>თქვენ შეგიძლიათ ჩამოტვირთოთ ფაილი თითოეული წიგნისთვის
შესაბამისი წიგნის თავების ჩამონათვალის translate გვერდიდან
ანდა ჩამოტვირთოთ თითოეული თავისთვის ფაილი თავისი გვერდიდან</p>
<p>
naierchou@proton.me
</p>
</main><script>
function downloadJSON() {
    const storageData = {};
    for (let i = 0; i < localStorage.length; i++) {
        const key = localStorage.key(i);
        storageData[key] = localStorage.getItem(key);
    }

    const jsonString = JSON.stringify(storageData, null, 2); // 2 spaces for indentation

    const filename = "localStorage-backup.json"

    goSaveJSON(jsonString, filename)
}
</script></body></html>`)
