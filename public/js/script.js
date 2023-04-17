const menuToggle = document.querySelector('.menu-toggle input');
const nav = document.querySelector('.nav-left')

menuToggle.addEventListener('click' , function() {
    nav.classList.toggle('slide')
})


function getData() {
    let name = document.getElementById("name").value
    let email = document.getElementById("email").value
    let phone = document.getElementById("phone").value
    let subject = document.getElementById("subject").value
    let message = document.getElementById("message").value

    //validation
    if (name == "" || email == "" || phone == "" || subject == "" || message == "") {
        return alert("Semua Field harus diisi")
    }

    let destination = "skadskuds.f@gmail.com"
    const contact = document.createElement("a")
    contact.setAttribute('target', '_blank');
    contact.href = `https://mail.google.com/mail/?view=cm&fs=1&to=${destination}&su=${subject}&body=Hallo nama saya ${name} saya ingin ${message} ,bila berkenan hubungi saya di ${phone}`
    contact.click()
}

const listproject = []

function getProject(event) {
    event.preventDefault()

    let project = document.getElementById("project-name").value
    let start = document.getElementById("project-start").value
    let end = document.getElementById("project-end").value
    let description = document.getElementById("project-desc").value
    let img = document.getElementById("img-file").files
    let node = document.getElementById("node").checked
    let react = document.getElementById("react").checked
    let next = document.getElementById("next").checked
    let typescript = document.getElementById("typescript").checked

    //validation
    if (img.length !== 0) {
        img = URL.createObjectURL(img[0])
    }
    if (!node && !react && !next && !typescript) {
        alert("harap pilih minimal 1 teknologies")
    }

    let blog = {
        project,
        start,
        end,
        description,
        postedAt:new Date(),
        img,
        node,
        react,
        next,
        typescript
    }
    listproject.push(blog);
    console.log(listproject);
    renderProject();
}

function renderProject() {
    document.getElementById("project-list").innerHTML = ""
    if(listproject.length == 0) {
        document.getElementById("project-list").innerHTML +=`<div style="margin:auto;">
        <h1>No Project Yet</h1> </div>`
    }else{
        for (let i = 0; i < listproject.length; i++) {
            document.getElementById("project-list").innerHTML +=
                `<div class="myproject-card">
                <img src="${listproject[i].img}" style="height: 200px;" alt="posted">
                <a style="font-weight: bold; cursor: pointer; text-decoratui" href="detail-project.html">${listproject[i].project}</a>
                <p style="color: #8d8d8d;">durasi : ${createDuration(listproject[i].start, listproject[i].end)}</p>
                <div class="desc-project">
                <article>${listproject[i].description}</article>
                </div>
                <div class="tech-used">
                ${listproject[i].node ? "<img src='assets/images/atom.png'>" : ""}
                ${listproject[i].react ? "<img src='assets/images/nodejs.png'>" : ""}
                ${listproject[i].next ? "<img src='assets/images/next-js.png'>" : ""}
                ${listproject[i].typescript ? "<img src='assets/images/typescript.png' >" : ""}
                </div>
                <div style="display: flex; padding: 20px; ">
                        <button class="action">edit</button>
                        <button class="action">delete</button>
                    </div>
                    <div style="float: right; padding: 10px;">
                        <span>Posted at :</span><p style="color: #8d8d8d;"> ${getPosted(listproject[i].postedAt)}</p>
                    </div>
            </div>  
        </div>`
        }
    }
}

function createDuration(start, end) {
    let distance = new Date(end) - new Date(start)
    // validation
    const dayDistance = Math.floor(distance / ( 24 * 60 * 60 * 1000 ))
    const monthDistance = Math.floor(distance /( 30 * 24 * 60 * 60 * 1000 ))
    const yearDistance = Math.floor (distance /( 12 * 30 * 24 * 60 * 60 * 1000 ))
    if( dayDistance > 0 && dayDistance <= 31 ){
        return `${dayDistance} hari `
    }else if(dayDistance > 31 && monthDistance <= 12){ 
            const newDay = Math.floor(dayDistance % 31)
            return`${newDay} hari ${monthDistance} bulan`
        } else if(monthDistance > 12){
        const newMonth = Math.floor( monthDistance % 12 )
        return `${newMonth} bulan ${yearDistance} tahun`
    }
}

function getPosted(time){
    const posted = new Date() - new Date(time)
    const dayPosted = Math.floor(posted / ( 24 * 60 * 60 * 1000 ))
    if(dayPosted > 0) {
        return dayPosted + " day ago"
    }else {
        const hourPosted = Math.floor(posted/ ( 60* 60* 1000))
        if(hourPosted > 0 ){
            return hourPosted + " hour ago"
        }else {
            const minutesPosted = Math.floor(posted / ( 60 * 1000))
            if(minutesPosted > 0 ){
                return minutesPosted + " minutes ago"
            } else{
                const secondPosted = Math.floor(posted/1000)
                if(secondPosted > 0) {
                    return secondPosted + " second ago"
                }
            }
        }
    }
}

setInterval(renderProject, 1000)