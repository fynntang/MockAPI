package server

const appJS = `
const MOCK_BASE = "/mock";
const API = "/_api/routes";
const CONFIG_API = "/_api/config";

document.getElementById("mock-url").textContent = location.origin + MOCK_BASE;

// --- Tabs ---
document.querySelectorAll(".tab").forEach(tab => {
  tab.addEventListener("click", () => {
    document.querySelectorAll(".tab").forEach(t => t.classList.remove("active"));
    document.querySelectorAll(".tab-content").forEach(c => c.classList.remove("active"));
    tab.classList.add("active");
    document.getElementById("tab-" + tab.dataset.tab).classList.add("active");
    if (tab.dataset.tab === "logs") loadLogs();
    if (tab.dataset.tab === "settings") loadSettings();
  });
});

// --- Routes ---
let editingId = null;

async function loadRoutes() {
  const res = await fetch(API);
  const routes = await res.json();
  const el = document.getElementById("routes");
  const empty = document.getElementById("empty");
  const filter = (document.getElementById("route-search")?.value || "").toLowerCase();

  const filtered = routes.filter(r =>
    !filter || r.method.toLowerCase().includes(filter) ||
    r.path.toLowerCase().includes(filter) ||
    (r.description || "").toLowerCase().includes(filter)
  );

  if (filtered.length === 0) {
    el.innerHTML = "";
    empty.style.display = "";
    return;
  }
  empty.style.display = "none";

  el.innerHTML = filtered.map(r => {
    const mc = r.method.toLowerCase();
    const desc = r.description ? '<span class="desc">' + escapeHtml(r.description) + '</span>' : '';
    const cond = (r.match_headers || r.match_body) ? '<span class="condition-badge" title="Conditional">⚡</span>' : '';
    return '<div class="route">' +
      '<div class="route-info">' +
        '<span class="method ' + mc + '">' + r.method + '</span>' +
        '<code>' + escapeHtml(r.path) + '</code>' +
        cond + desc +
        '<span class="status">' + r.status + (r.delay ? ' · ' + r.delay + 'ms' : '') + '</span>' +
      '</div>' +
      '<div class="route-actions">' +
        '<button class="copy" onclick="copyCurl(' + JSON.stringify(r) + ')" title="Copy curl">📋</button>' +
        '<button class="copy" onclick="editRoute(' + JSON.stringify(r) + ')" title="Edit">✏️</button>' +
        '<button class="del" onclick="deleteRoute(\'' + r.id + '\')" title="Delete">✕</button>' +
      '</div>' +
    '</div>';
  }).join("");
}

function openModal() {
  editingId = null;
  document.getElementById("modal-title").textContent = "Add Mock Route";
  document.getElementById("f-method").value = "GET";
  document.getElementById("f-path").value = "";
  document.getElementById("f-status").value = "200";
  document.getElementById("f-delay").value = "0";
  document.getElementById("f-desc").value = "";
  document.getElementById("f-body").value = "";
  document.getElementById("f-headers").value = "";
  document.getElementById("f-match-headers").value = "";
  document.getElementById("f-match-body").value = "";
  document.getElementById("modal").style.display = "";
}

function editRoute(r) {
  editingId = r.id;
  document.getElementById("modal-title").textContent = "Edit Route";
  document.getElementById("f-method").value = r.method;
  document.getElementById("f-path").value = r.path;
  document.getElementById("f-status").value = r.status;
  document.getElementById("f-delay").value = r.delay || 0;
  document.getElementById("f-desc").value = r.description || "";
  document.getElementById("f-body").value = r.body || "";
  document.getElementById("f-headers").value = r.headers ? JSON.stringify(r.headers, null, 2) : "";
  document.getElementById("f-match-headers").value = r.match_headers ? JSON.stringify(r.match_headers, null, 2) : "";
  document.getElementById("f-match-body").value = r.match_body || "";
  document.getElementById("modal").style.display = "";
}

function closeModal() {
  document.getElementById("modal").style.display = "none";
}

async function saveRoute() {
  let headers = {};
  const h = document.getElementById("f-headers").value.trim();
  if (h) {
    try { headers = JSON.parse(h); } catch(e) { alert("Invalid headers JSON"); return; }
  }

  let matchHeaders = {};
  const mh = document.getElementById("f-match-headers").value.trim();
  if (mh) {
    try { matchHeaders = JSON.parse(mh); } catch(e) { alert("Invalid match headers JSON"); return; }
  }

  const route = {
    method: document.getElementById("f-method").value,
    path: document.getElementById("f-path").value,
    status: parseInt(document.getElementById("f-status").value),
    delay: parseInt(document.getElementById("f-delay").value) || 0,
    description: document.getElementById("f-desc").value.trim(),
    body: document.getElementById("f-body").value,
    headers: Object.keys(headers).length ? headers : undefined,
    match_headers: Object.keys(matchHeaders).length ? matchHeaders : undefined,
    match_body: document.getElementById("f-match-body").value.trim() || undefined,
  };

  if (!route.path) { alert("Path is required"); return; }

  const httpMethod = editingId ? "PUT" : "POST";
  if (editingId) route.id = editingId;

  await fetch(API, {
    method: httpMethod,
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(route),
  });

  closeModal();
  loadRoutes();
}

async function deleteRoute(id) {
  if (!confirm("Delete this route?")) return;
  await fetch(API + "?id=" + id, { method: "DELETE" });
  loadRoutes();
}

function copyCurl(r) {
  const url = location.origin + MOCK_BASE + r.path;
  const cmd = "curl -X " + r.method + " " + url;
  navigator.clipboard.writeText(cmd);
}

// --- Import / Export ---
function exportRoutes() { location.href = "/_api/export"; }

async function importRoutes(e) {
  const file = e.target.files[0];
  if (!file) return;
  const text = await file.text();
  try { JSON.parse(text); } catch(err) { alert("Invalid JSON file"); return; }
  const res = await fetch("/_api/import", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: text,
  });
  const data = await res.json();
  alert("Imported " + data.imported + " routes!");
  loadRoutes();
  e.target.value = "";
}

// --- Templates ---
async function openTemplates() {
  document.getElementById("templates-modal").style.display = "";
  const res = await fetch("/_api/templates");
  const templates = await res.json();
  document.getElementById("template-list").innerHTML = templates.map(t => {
    const mc = t.method.toLowerCase();
    const cond = (t.match_headers || t.match_body) ? ' <span class="condition-badge">⚡</span>' : '';
    return '<div class="template-item" onclick="useTemplate(this)" ' +
      'data-method="' + t.method + '" data-path="' + escapeAttr(t.path) + '" ' +
      'data-status="' + t.status + '" data-delay="' + (t.delay||0) + '" ' +
      'data-desc="' + escapeAttr(t.description||'') + '" ' +
      'data-body="' + escapeAttr(t.body) + '" ' +
      'data-headers="' + escapeAttr(t.headers ? JSON.stringify(t.headers) : '') + '" ' +
      'data-match-headers="' + escapeAttr(t.match_headers ? JSON.stringify(t.match_headers) : '') + '" ' +
      'data-match-body="' + escapeAttr(t.match_body||'') + '">' +
      '<div class="template-top">' +
        '<span class="method ' + mc + '">' + t.method + '</span>' +
        '<code>' + escapeHtml(t.path) + '</code>' + cond +
        '<span class="status">' + t.status + (t.delay ? ' · ' + t.delay + 'ms' : '') + '</span>' +
      '</div>' +
      (t.description ? '<div class="template-desc">' + escapeHtml(t.description) + '</div>' : '') +
    '</div>';
  }).join("");
}

function closeTemplates() { document.getElementById("templates-modal").style.display = "none"; }

function useTemplate(el) {
  const d = el.dataset;
  document.getElementById("f-method").value = d.method;
  document.getElementById("f-path").value = d.path;
  document.getElementById("f-status").value = d.status;
  document.getElementById("f-delay").value = d.delay;
  document.getElementById("f-desc").value = d.desc;
  document.getElementById("f-body").value = d.body;
  document.getElementById("f-headers").value = d.headers || "";
  document.getElementById("f-match-headers").value = d.matchHeaders || "";
  document.getElementById("f-match-body").value = d.matchBody || "";
  closeTemplates();
  document.getElementById("modal").style.display = "";
  document.getElementById("modal-title").textContent = "Add Mock Route";
  editingId = null;
}

// --- Logs ---
async function loadLogs() {
  const res = await fetch("/_api/logs");
  const logs = await res.json();
  const el = document.getElementById("logs");
  const empty = document.getElementById("logs-empty");
  document.getElementById("log-count").textContent = logs.length + " requests";

  if (logs.length === 0) { el.innerHTML = ""; empty.style.display = ""; return; }
  empty.style.display = "none";

  el.innerHTML = logs.slice().reverse().map(l => {
    const sc = l.status >= 200 && l.status < 300 ? 'ok' : l.status >= 400 ? 'err' : '';
    const px = l.proxied ? '<span class="proxy-badge">PROXY</span>' : '';
    return '<div class="log-item">' +
      '<span class="log-time">' + l.timestamp + '</span>' +
      px +
      '<span class="method ' + l.method.toLowerCase() + '">' + l.method + '</span>' +
      '<code>' + escapeHtml(l.path) + '</code>' +
      '<span class="log-status ' + sc + '">' + l.status + '</span>' +
      '<span class="log-delay">' + l.delay + 'ms</span>' +
    '</div>';
  }).join("");
}

async function clearLogs() {
  if (!confirm("Clear all logs?")) return;
  await fetch("/_api/clear-logs", { method: "POST" });
  loadLogs();
}

// --- Settings ---
async function loadSettings() {
  const res = await fetch(CONFIG_API);
  const cfg = await res.json();
  document.getElementById("s-proxy").value = cfg.proxy_url || "";
  document.getElementById("s-cors").checked = cfg.cors_enabled;
  document.getElementById("s-maxlogs").value = cfg.max_logs || 500;
}

async function saveSettings() {
  await fetch(CONFIG_API, {
    method: "PUT",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      proxy_url: document.getElementById("s-proxy").value.trim(),
      cors_enabled: document.getElementById("s-cors").checked,
      max_logs: parseInt(document.getElementById("s-maxlogs").value) || 500,
    }),
  });
  alert("Settings saved!");
}

// --- Utils ---
function escapeHtml(s) {
  return (s||"").replace(/&/g,"&amp;").replace(/</g,"&lt;").replace(/>/g,"&gt;").replace(/"/g,"&quot;");
}
function escapeAttr(s) {
  return (s||"").replace(/&/g,"&amp;").replace(/"/g,"&quot;").replace(/</g,"&lt;").replace(/>/g,"&gt;");
}

// Auto-refresh logs
setInterval(() => {
  if (document.getElementById("tab-logs").classList.contains("active")) loadLogs();
}, 3000);

loadRoutes();
`
