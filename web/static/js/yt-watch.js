// Player for official YouTube VODs (see watch.html). Plays one or more
// YouTube videos (e.g. one per map) in order, saves watch progress for
// logged-in profiles and, in spoiler-free mode, swaps YouTube's controls
// (whose seek bar reveals how long a map goes) for plain skip buttons.
(function () {
  var root = document.getElementById("yt-watch");
  if (!root) return;

  var data = root.dataset;
  var ids = (data.videoIds || "").split(",").filter(Boolean);
  if (!ids.length) return;

  var hideSpoilers = data.hideSpoilers === "1";
  var part = parseInt(data.initialPart || "0", 10);
  if (isNaN(part) || part < 0 || part >= ids.length) part = 0;
  var startAt = parseInt(data.initialProgress || "0", 10);
  if (isNaN(startAt) || startAt < 0) startAt = 0;

  var player = null;
  var saveTimer = null;
  var consent = document.getElementById("yt-consent");
  var controls = document.getElementById("yt-controls");
  var toggleBtn = document.getElementById("yt-toggle");
  var elapsed = document.getElementById("yt-elapsed");

  function currentTime() {
    try {
      return Math.floor(player.getCurrentTime() || 0);
    } catch (e) {
      return 0;
    }
  }

  // seconds overrides the player position (used right after switching
  // parts, when the player still reports the end of the previous one).
  function saveProgress(completed, seconds) {
    if (!data.profileId || !player) return;
    var body = {
      mediaId: data.mediaId,
      profileId: data.profileId,
      episodeId: data.episodeId,
      progressSeconds: completed ? 0 : seconds !== undefined ? seconds : currentTime(),
      partIndex: completed ? 0 : part,
      completed: !!completed,
    };
    fetch("/api/watch/progress", {
      method: "POST",
      headers: { "Content-Type": "application/json", "X-CSRF-Token": data.csrfToken },
      body: JSON.stringify(body),
      keepalive: true,
    }).catch(function () {});
  }

  function formatTime(total) {
    var h = Math.floor(total / 3600);
    var m = Math.floor((total % 3600) / 60);
    var s = total % 60;
    var mm = (h > 0 && m < 10 ? "0" : "") + m;
    return (h > 0 ? h + ":" : "") + mm + ":" + (s < 10 ? "0" : "") + s;
  }

  function updateControls() {
    if (!player || !hideSpoilers) return;
    var playing = player.getPlayerState && player.getPlayerState() === YT.PlayerState.PLAYING;
    toggleBtn.textContent = playing ? "Pause" : "Abspielen";
    // Only the elapsed time is shown, never the total length.
    elapsed.textContent = (ids.length > 1 ? "Teil " + (part + 1) + " · " : "") + formatTime(currentTime());
  }

  function onStateChange(e) {
    updateControls();
    if (e.data !== YT.PlayerState.ENDED) return;

    if (part + 1 < ids.length) {
      part += 1;
      saveProgress(false, 0);
      player.loadVideoById({ videoId: ids[part], startSeconds: 0 });
      return;
    }

    saveProgress(true);
    clearInterval(saveTimer);
    if (data.nextHref) window.location.href = data.nextHref;
  }

  // Region locks, age restrictions or a removed video make the embed fail;
  // offer the original on YouTube instead of a dead player.
  function showFallback() {
    var box = document.getElementById("yt-fallback");
    var link = document.getElementById("yt-fallback-link");
    link.href = "https://www.youtube.com/watch?v=" + encodeURIComponent(ids[part]);
    box.hidden = false;
  }

  function createPlayer() {
    player = new YT.Player("yt-player", {
      host: "https://www.youtube-nocookie.com",
      videoId: ids[part],
      playerVars: {
        autoplay: 1,
        playsinline: 1,
        rel: 0,
        iv_load_policy: 3,
        start: startAt,
        controls: hideSpoilers ? 0 : 1,
        disablekb: hideSpoilers ? 1 : 0,
      },
      events: {
        onReady: function () {
          player.playVideo();
          if (hideSpoilers) controls.hidden = false;
          updateControls();
        },
        onStateChange: onStateChange,
        onError: showFallback,
      },
    });

    // Tick every second for the elapsed-time display, save every 15 s.
    var ticks = 0;
    saveTimer = setInterval(function () {
      updateControls();
      ticks += 1;
      if (ticks % 15 === 0 && player.getPlayerState && player.getPlayerState() === YT.PlayerState.PLAYING) {
        saveProgress(false);
      }
    }, 1000);

    window.addEventListener("beforeunload", function () {
      saveProgress(false);
    });
  }

  controls.addEventListener("click", function (e) {
    var action = e.target && e.target.getAttribute("data-action");
    if (!action || !player) return;
    var skips = { back30: -30, back10: -10, fwd10: 10, fwd30: 30, fwd120: 120 };
    if (action === "toggle") {
      if (player.getPlayerState() === YT.PlayerState.PLAYING) player.pauseVideo();
      else player.playVideo();
    } else if (skips[action]) {
      player.seekTo(Math.max(0, currentTime() + skips[action]), true);
    }
    updateControls();
  });

  // Nothing is requested from YouTube until the viewer clicks play.
  consent.addEventListener("click", function (e) {
    if (e.target.tagName === "A") return;
    consent.remove();
    window.onYouTubeIframeAPIReady = createPlayer;
    var tag = document.createElement("script");
    tag.src = "https://www.youtube.com/iframe_api";
    document.head.appendChild(tag);
  });
})();
