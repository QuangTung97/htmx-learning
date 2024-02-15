document.body.addEventListener("htmx:sendError", function (e) {
  alert("Error: Your network is disconnected");
});

function getCookie(name) {
  const cookies = document.cookie.split(',');
  for (let i = 0; i < cookies.length; i++) {
    const c = cookies[i];
    const cookieKey = name + "=";
    if (c.startsWith(cookieKey)) {
      return c.slice(cookieKey.length);
    }
  }
  return "";
}

document.body.addEventListener("htmx:beforeRequest", function (e) {
  if (e.detail.requestConfig.verb !== "get") {
    e.detail.xhr.setRequestHeader("X-Csrf-Token", getCookie("csrf_token"));
  }
});
