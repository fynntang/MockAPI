package server

const appJS = `
const MOCK_BASE = "/mock";
const API = "/_api/routes";

document.getElementById("mock-url").textContent = location.origin + MOCK_BASE;

async function loadRoutes() {
  const res = await fetch(API);
  const routes = await res.json();
  const el = document.getElementById("routes");
  const empty = document.getElementById("empty");

  if (routes.length === 0) {
    el.innerHTML = "";
    empty.style.display = "";
    return;
  }
  empty.style.display = "none";

  el.innerHTML = routes.map(r => {
    const methodClass = r.method.toLowerCase();
    return "<div class='route'>" +
      "<div class='route-info'>" +
        "<span class='method " + methodClass + "'>" + r.method + "</span>" +
        "<code>" + escapeHtml(MOCK_BASE + r.path) + "</code>" +
        "<span class='status'>→ " + r.status + (r.delay ? " (" + r.delay + "ms)" : "") + "</span>" +
      "</div>" +
      "<div class='route-actions'>" +
        "<button class='copy' onclick='copyCurl(" + JSON.stringify(r) + ")'>Copy curl</button>" +
        "<button class='del' onclick='deleteRoute(\"" + r.id + "\")'>✕</button>" +
      "</div>" +
    "</div>";
  }).join("");
}

function openModal() {
  document.getElementById("modal").style.display = "";
}

function closeModal() {
  document.getElementById("modal").style.display = "none";
}

async function addRoute() {
  let headers = {};
  const h = document.getElementById("f-headers").value.trim();
  if (h) {
    try { headers = JSON.parse(h); } catch(e) { alert("Invalid headers JSON"); return; }
  }

  const body = document.getElementById("f-body").value.trim() || "{}";

  const route = {
    method: document.getElementById("f-method").value,
    path: document.getElementById("f-path").value,
    status: parseInt(document.getElementById("f-status").value),
    delay: parseInt(document.getElementById("f-delay").value) || 0,
    body: body,
    headers: headers,
  };

  if (!route.path) { alert("Path is required"); return; }

  await fetch(API, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(route),
  });

  closeModal();
  loadRoutes();
}

async function deleteRoute(id) {
  await fetch(API + "?id=" + id, { method: "DELETE" });
  loadRoutes();
}

function copyCurl(r) {
  const url = location.origin + MOCK_BASE + r.path;
  const cmd = "curl -X " + r.method + " " + url;
  navigator.clipboard.writeText(cmd).then(() => {
    alert("Copied: " + cmd);
  });
}

function escapeHtml(s) {
  return s.replace(/&/g,"&amp;").replace(/</g,"&lt;").replace(/>/g,"&gt;");
}

loadRoutes();
`
