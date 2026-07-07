(function () {
  var el = document.getElementById("player");
  if (!el) return;

  var data = el.dataset;
  var initialProgress = parseInt(data.initialProgress || "0", 10);

  if (initialProgress > 0 && !isNaN(initialProgress)) {
    el.addEventListener("loadedmetadata", function () {
      try {
        el.currentTime = initialProgress;
      } catch (e) {}
    }, { once: true });
  }

  function saveProgress(completed) {
    if (el.paused && !completed) return;
    var body = {
      mediaId: data.mediaId,
      profileId: data.profileId,
      progressSeconds: Math.floor(el.currentTime || 0),
      completed: !!completed,
    };
    if (data.episodeId) body.episodeId = data.episodeId;

    fetch("/api/watch/progress", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "X-CSRF-Token": data.csrfToken,
      },
      body: JSON.stringify(body),
      keepalive: true,
    }).catch(function () {});
  }

  var interval = setInterval(function () {
    saveProgress(false);
  }, 15000);

  el.addEventListener("ended", function () {
    saveProgress(true);
    clearInterval(interval);
    if (data.nextHref) {
      window.location.href = data.nextHref;
    }
  });

  window.addEventListener("beforeunload", function () {
    saveProgress(false);
  });
})();
