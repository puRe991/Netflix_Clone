(function () {
  var input = document.getElementById("search-input");
  if (!input) return;

  var cards = Array.prototype.slice.call(document.querySelectorAll(".card[data-title]"));

  input.addEventListener("input", function () {
    var q = input.value.trim().toLowerCase();
    cards.forEach(function (card) {
      var title = (card.getAttribute("data-title") || "").toLowerCase();
      card.style.display = q === "" || title.indexOf(q) !== -1 ? "" : "none";
    });
  });
})();
