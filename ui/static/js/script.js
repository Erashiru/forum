button = document.getElementById("home-button-error");

button.addEventListener("click", function() {
  window.location = "/";
})

function toggleMenu() {
  const menu = document.getElementById("userOptionsMenu");
  if (menu.style.right === "0px") {
      menu.style.right = "-250px";
  } else {
      menu.style.right = "0px";
  }
}