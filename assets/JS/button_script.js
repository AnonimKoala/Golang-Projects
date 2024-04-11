function show_nav() {
    var element = document.getElementById("main_nav");
    element.classList.add("nav_clicked");
}

function hide_nav() {
    var element = document.getElementById("main_nav");
    element.classList.remove("nav_clicked");
}