let content = document.querySelector("#content")
if (content){
    let tags = document.getElementsByClassName("tag")
    let divTag = document.createElement("div")
    divTag.id = "tags-container"
    let set = new Set()

    for (let it of tags){
        let vals = it.textContent
        
        vals = vals.replace(/\u00A0/, " ")
        vals.split(" ").forEach(i=>{
        set.add(i)
    })
    console.warn(it)
}
while(tags.length > 0){
    tags[0].remove()
}
let sp = document.createElement("span")
sp.className = "tags-title"
sp.id = "tags-title"
sp.textContent = "Tags:"
divTag.append(sp)
set.forEach(i=>{
    let a = document.createElement("a")
    a.href = "/?label="+i
    a.textContent = i
    a.className = "tag-item"
    divTag.append(a)
})
content.insertBefore(divTag, content.children[1])
//content.prepend(divTag)
}