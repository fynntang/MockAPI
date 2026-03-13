package server

const appJS = `
const MOCK_BASE = "/mock";
const WS_BASE = "/ws";
const API = "/_api/routes";
const CONFIG_API = "/_api/config";
const WS_API = "/_api/ws";

document.getElementById("mock-url").textContent = location.origin + MOCK_BASE;
document.getElementById("ws-url").textContent = location.origin + WS_BASE;
document.getElementById("graphql-url").textContent = location.origin + "/graphql";
document.getElementById("grpc-url").textContent = location.origin + "/grpc";

// --- Tabs ---
document.querySelectorAll(".tab").forEach(tab => {
  tab.addEventListener("click", () => {
    document.querySelectorAll(".tab").forEach(t => t.classList.remove("active"));
    document.querySelectorAll(".tab-content").forEach(c => c.classList.remove("active"));
    tab.classList.add("active");
    document.getElementById("tab-" + tab.dataset.tab).classList.add("active");
    if (tab.dataset.tab === "logs") loadLogs();
    if (tab.dataset.tab === "settings") loadSettings();
    if (tab.dataset.tab === "ws") loadWSHandlers();
  });
});

// --- Routes ---
let editingId = null;
let editingWSPath = null;

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

  if (filtered.length === 0) { el.innerHTML = ""; empty.style.display = ""; return; }
  empty.style.display = "none";

  el.innerHTML = filtered.map(r => {
    const mc = r.method.toLowerCase();
    const desc = r.description ? '<span class="desc">' + escapeHtml(r.description) + '</span>' : '';
    const cond = (r.match_headers || r.match_body) ? '<span class="condition-badge" title="Conditional">⚡</span>' : '';
    const scriptBadge = r.script ? '<span class="script-badge" title="Script">📜</span>' : '';
    const routeData = encodeURIComponent(JSON.stringify(r));
    return '<div class="route">' +
      '<div class="route-info">' +
        '<span class="method ' + mc + '">' + r.method + '</span>' +
        '<code>' + escapeHtml(r.path) + '</code>' +
        scriptBadge + cond + desc +
        '<span class="status">' + r.status + (r.delay ? ' · ' + r.delay + 'ms' : '') + '</span>' +
      '</div>' +
      '<div class="route-actions">' +
        '<button class="copy" onclick="copyCurl(\'' + routeData + '\')" title="Copy curl">📋</button>' +
        '<button class="copy" onclick="editRoute(\'' + routeData + '\')" title="Edit">✏️</button>' +
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
  document.getElementById("f-script").value = "";
  document.getElementById("f-match-headers").value = "";
  document.getElementById("f-match-body").value = "";
  document.getElementById("f-type").value = "static";
  toggleResponseType();
  document.getElementById("modal").style.display = "";
}

function editRoute(data) {
  const r = JSON.parse(decodeURIComponent(data));
  editingId = r.id;
  document.getElementById("modal-title").textContent = "Edit Route";
  document.getElementById("f-method").value = r.method;
  document.getElementById("f-path").value = r.path;
  document.getElementById("f-status").value = r.status;
  document.getElementById("f-delay").value = r.delay || 0;
  document.getElementById("f-desc").value = r.description || "";
  document.getElementById("f-body").value = r.body || "";
  document.getElementById("f-headers").value = r.headers ? JSON.stringify(r.headers, null, 2) : "";
  document.getElementById("f-script").value = r.script || "";
  document.getElementById("f-match-headers").value = r.match_headers ? JSON.stringify(r.match_headers, null, 2) : "";
  document.getElementById("f-match-body").value = r.match_body || "";
  document.getElementById("f-type").value = r.script ? "script" : "static";
  toggleResponseType();
  document.getElementById("modal").style.display = "";
}

function closeModal() { document.getElementById("modal").style.display = "none"; }

function toggleResponseType() {
  const type = document.getElementById("f-type").value;
  document.getElementById("static-response").style.display = type === "static" ? "" : "none";
  document.getElementById("script-response").style.display = type === "script" ? "" : "none";
}

async function saveRoute() {
  let headers = {};
  const h = document.getElementById("f-headers").value.trim();
  if (h) { try { headers = JSON.parse(h); } catch(e) { alert("Invalid headers JSON"); return; } }

  let matchHeaders = {};
  const mh = document.getElementById("f-match-headers").value.trim();
  if (mh) { try { matchHeaders = JSON.parse(mh); } catch(e) { alert("Invalid match headers JSON"); return; } }

  const type = document.getElementById("f-type").value;
  
  const route = {
    method: document.getElementById("f-method").value,
    path: document.getElementById("f-path").value,
    status: parseInt(document.getElementById("f-status").value),
    delay: parseInt(document.getElementById("f-delay").value) || 0,
    description: document.getElementById("f-desc").value.trim(),
    headers: Object.keys(headers).length ? headers : undefined,
    match_headers: Object.keys(matchHeaders).length ? matchHeaders : undefined,
    match_body: document.getElementById("f-match-body").value.trim() || undefined,
  };

  if (type === "script") {
    route.script = document.getElementById("f-script").value;
    route.body = "";
  } else {
    route.body = document.getElementById("f-body").value;
  }

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

function copyCurl(data) {
  const r = JSON.parse(decodeURIComponent(data));
  const url = location.origin + MOCK_BASE + r.path;
  const cmd = "curl -X " + r.method + " " + url;
  navigator.clipboard.writeText(cmd);
}

// --- Swagger Import ---
function openSwaggerModal() { document.getElementById("swagger-modal").style.display = ""; }
function closeSwaggerModal() { document.getElementById("swagger-modal").style.display = "none"; }

function loadSwaggerFile(e) {
  const file = e.target.files[0];
  if (!file) return;
  const reader = new FileReader();
  reader.onload = (ev) => { document.getElementById("swagger-input").value = ev.target.result; };
  reader.readAsText(file);
}

async function importSwagger() {
  const input = document.getElementById("swagger-input").value.trim();
  if (!input) { alert("Please paste or upload an OpenAPI spec"); return; }
  
  const res = await fetch("/_api/import-swagger", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: input,
  });
  
  if (!res.ok) { const err = await res.text(); alert("Import failed: " + err); return; }
  
  const data = await res.json();
  alert("Imported " + data.imported + " routes from OpenAPI spec!");
  closeSwaggerModal();
  loadRoutes();
}

// --- WebSocket ---
async function loadWSHandlers() {
  const res = await fetch(WS_API);
  const handlers = await res.json();
  const el = document.getElementById("ws-handlers");
  const empty = document.getElementById("ws-empty");

  if (handlers.length === 0) { el.innerHTML = ""; empty.style.display = ""; return; }
  empty.style.display = "none";

  el.innerHTML = handlers.map(h => {
    const desc = h.description ? '<span class="desc">' + escapeHtml(h.description) + '</span>' : '';
    const streamBadge = h.stream_enabled ? '<span class="stream-badge" title="Stream Mode">📡</span>' : '';
    const delayBadge = h.delay ? '<span class="status">' + h.delay + 'ms</span>' : '';
    const handlerData = encodeURIComponent(JSON.stringify(h));
    return '<div class="route">' +
      '<div class="route-info">' +
        '<span class="method ws">WS</span>' +
        '<code>' + escapeHtml(h.path) + '</code>' +
        streamBadge + desc + delayBadge +
      '</div>' +
      '<div class="route-actions">' +
        '<button class="copy" onclick="copyWSUrl(\'' + escapeHtml(h.path) + '\')" title="Copy URL">📋</button>' +
        '<button class="copy" onclick="editWSHandler(\'' + handlerData + '\')" title="Edit">✏️</button>' +
        '<button class="del" onclick="deleteWSHandler(\'' + escapeHtml(h.path) + '\')" title="Delete">✕</button>' +
      '</div>' +
    '</div>';
  }).join("");
}

function toggleWSMode() {
  const mode = document.getElementById("ws-mode").value;
  document.getElementById("ws-reply-mode").style.display = mode === "reply" ? "" : "none";
  document.getElementById("ws-stream-mode").style.display = mode === "stream" ? "" : "none";
}

function toggleStreamInterval() {
  const type = document.getElementById("ws-stream-interval-type").value;
  document.getElementById("ws-fixed-interval").style.display = type === "fixed" ? "" : "none";
  document.getElementById("ws-random-interval").style.display = type === "random" ? "" : "none";
}

function openWSModal() {
  editingWSPath = null;
  document.getElementById("ws-modal-title").textContent = "Add WebSocket Handler";
  document.getElementById("ws-path").value = "";
  document.getElementById("ws-desc").value = "";
  document.getElementById("ws-mode").value = "reply";
  document.getElementById("ws-delay").value = "0";
  document.getElementById("ws-auto-reply").value = "";
  document.getElementById("ws-on-connect").value = "";
  document.getElementById("ws-on-message").value = "";
  // Stream mode fields
  document.getElementById("ws-stream-messages").value = "";
  document.getElementById("ws-stream-interval-type").value = "fixed";
  document.getElementById("ws-stream-interval").value = "1000";
  document.getElementById("ws-stream-min-delay").value = "500";
  document.getElementById("ws-stream-max-delay").value = "3000";
  document.getElementById("ws-stream-format").value = "json";
  document.getElementById("ws-stream-loop").checked = true;
  document.getElementById("ws-stream-on-connect").value = "";
  toggleWSMode();
  toggleStreamInterval();
  document.getElementById("ws-modal").style.display = "";
}

function editWSHandler(data) {
  const h = JSON.parse(decodeURIComponent(data));
  editingWSPath = h.path;
  document.getElementById("ws-modal-title").textContent = "Edit WebSocket Handler";
  document.getElementById("ws-path").value = h.path;
  document.getElementById("ws-desc").value = h.description || "";
  
  // Determine mode
  const isStream = h.stream_enabled;
  document.getElementById("ws-mode").value = isStream ? "stream" : "reply";
  
  // Reply mode fields
  document.getElementById("ws-delay").value = h.delay_ms || h.delay || 0;
  document.getElementById("ws-auto-reply").value = h.auto_reply || "";
  document.getElementById("ws-on-connect").value = isStream ? (h.stream_on_connect || "") : (h.on_connect || "");
  document.getElementById("ws-on-message").value = h.on_message || "";
  
  // Stream mode fields
  document.getElementById("ws-stream-messages").value = (h.stream_messages || []).join("\n");
  document.getElementById("ws-stream-interval").value = h.stream_interval_ms || 1000;
  document.getElementById("ws-stream-min-delay").value = h.stream_min_delay_ms || 500;
  document.getElementById("ws-stream-max-delay").value = h.stream_max_delay_ms || 3000;
  document.getElementById("ws-stream-format").value = h.stream_format || "json";
  document.getElementById("ws-stream-loop").checked = h.stream_loop !== false;
  
  // Set interval type
  if (h.stream_random) {
    document.getElementById("ws-stream-interval-type").value = "random";
  } else {
    document.getElementById("ws-stream-interval-type").value = "fixed";
  }
  
  toggleWSMode();
  toggleStreamInterval();
  document.getElementById("ws-modal").style.display = "";
}

function closeWSModal() { document.getElementById("ws-modal").style.display = "none"; }

async function saveWSHandler() {
  const mode = document.getElementById("ws-mode").value;
  
  const h = {
    path: document.getElementById("ws-path").value,
    description: document.getElementById("ws-desc").value,
  };
  
  if (!h.path) { alert("Path is required"); return; }
  
  if (mode === "stream") {
    // Stream mode
    const messagesText = document.getElementById("ws-stream-messages").value;
    const messages = messagesText.split("\n").map(m => m.trim()).filter(m => m);
    
    if (messages.length === 0) { alert("At least one message is required for stream mode"); return; }
    
    h.stream_enabled = true;
    h.stream_messages = messages;
    h.stream_format = document.getElementById("ws-stream-format").value;
    h.stream_loop = document.getElementById("ws-stream-loop").checked;
    
    const intervalType = document.getElementById("ws-stream-interval-type").value;
    if (intervalType === "random") {
      h.stream_random = true;
      h.stream_min_delay_ms = parseInt(document.getElementById("ws-stream-min-delay").value) || 500;
      h.stream_max_delay_ms = parseInt(document.getElementById("ws-stream-max-delay").value) || 3000;
    } else {
      h.stream_random = false;
      h.stream_interval_ms = parseInt(document.getElementById("ws-stream-interval").value) || 1000;
    }
    
    const onConnect = document.getElementById("ws-stream-on-connect").value.trim();
    if (onConnect) h.on_connect = onConnect;
  } else {
    // Reply mode
    h.delay = parseInt(document.getElementById("ws-delay").value) || 0;
    h.auto_reply = document.getElementById("ws-auto-reply").value;
    const onConnect = document.getElementById("ws-on-connect").value.trim();
    if (onConnect) h.on_connect = onConnect;
    h.on_message = document.getElementById("ws-on-message").value;
  }

  if (editingWSPath) {
    await fetch(WS_API + "?old_path=" + encodeURIComponent(editingWSPath), {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(h),
    });
  } else {
    await fetch(WS_API, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(h),
    });
  }

  closeWSModal();
  loadWSHandlers();
}

async function deleteWSHandler(path) {
  if (!confirm("Delete this WebSocket handler?")) return;
  await fetch(WS_API + "?path=" + encodeURIComponent(path), { method: "DELETE" });
  loadWSHandlers();
}

function copyWSUrl(path) {
  const url = location.origin.replace("http", "ws") + WS_BASE + path;
  navigator.clipboard.writeText(url);
  alert("Copied: " + url);
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
    const scr = t.script ? ' <span class="script-badge">📜</span>' : '';
    return '<div class="template-item" onclick="useTemplate(this)" ' +
      'data-method="' + t.method + '" data-path="' + escapeAttr(t.path) + '" ' +
      'data-status="' + t.status + '" data-delay="' + (t.delay||0) + '" ' +
      'data-desc="' + escapeAttr(t.description||'') + '" ' +
      'data-body="' + escapeAttr(t.body) + '" ' +
      'data-headers="' + escapeAttr(t.headers ? JSON.stringify(t.headers) : '') + '" ' +
      'data-script="' + escapeAttr(t.script||'') + '" ' +
      'data-match-headers="' + escapeAttr(t.match_headers ? JSON.stringify(t.match_headers) : '') + '" ' +
      'data-match-body="' + escapeAttr(t.match_body||'') + '">' +
      '<div class="template-top">' +
        '<span class="method ' + mc + '">' + t.method + '</span>' +
        '<code>' + escapeHtml(t.path) + '</code>' + scr + cond +
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
  document.getElementById("f-script").value = d.script || "";
  document.getElementById("f-match-headers").value = d.matchHeaders || "";
  document.getElementById("f-match-body").value = d.matchBody || "";
  document.getElementById("f-type").value = d.script ? "script" : "static";
  toggleResponseType();
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

async function clearLog() {
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

// --- GraphQL ---
async function loadGraphQLHandlers() {
  const res = await fetch("/_api/graphql");
  const handlers = await res.json();
  const el = document.getElementById("graphql-handlers");
  const empty = document.getElementById("graphql-empty");

  if (handlers.length === 0) { el.innerHTML = ""; empty.style.display = ""; return; }
  empty.style.display = "none";

  el.innerHTML = handlers.map(h => {
    const desc = h.description ? '<span class="desc">' + escapeHtml(h.description) + '</span>' : '';
    return '<div class="route">' +
      '<div class="route-info">' +
        '<span class="method gql">GQL</span>' +
        '<code>' + escapeHtml(h.operationName || "(catch-all)") + '</code>' +
        desc +
        (h.delay ? '<span class="status">' + h.delay + 'ms</span>' : '') +
      '</div>' +
      '<div class="route-actions">' +
        '<button class="copy" onclick="copyGQLQuery(\'' + escapeHtml(h.operationName || "") + '\')" title="Copy query">📋</button>' +
        '<button class="del" onclick="deleteGraphQLHandler(\'' + h.id + '\')" title="Delete">✕</button>' +
      '</div>' +
    '</div>';
  }).join("");
}

function openGraphQLModal() {
  document.getElementById("gql-op").value = "";
  document.getElementById("gql-desc").value = "";
  document.getElementById("gql-delay").value = "0";
  document.getElementById("gql-response").value = "";
  document.getElementById("graphql-modal").style.display = "";
}

function closeGraphQLModal() { document.getElementById("graphql-modal").style.display = "none"; }

async function saveGraphQLHandler() {
  let response = {};
  const respText = document.getElementById("gql-response").value.trim();
  if (respText) {
    try { response = JSON.parse(respText); } catch(e) { alert("Invalid JSON response"); return; }
  }

  const h = {
    operationName: document.getElementById("gql-op").value.trim(),
    description: document.getElementById("gql-desc").value.trim(),
    delay: parseInt(document.getElementById("gql-delay").value) || 0,
    response: response,
  };

  await fetch("/_api/graphql", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(h),
  });

  closeGraphQLModal();
  loadGraphQLHandlers();
}

async function deleteGraphQLHandler(id) {
  if (!confirm("Delete this GraphQL mock?")) return;
  await fetch("/_api/graphql?id=" + id, { method: "DELETE" });
  loadGraphQLHandlers();
}

function copyGQLQuery(opName) {
  const query = opName ? "{ " + opName + " }" : "{ __typename }";
  navigator.clipboard.writeText(query);
  alert("Copied: " + query);
}

// --- gRPC ---
async function loadGRPCHandlers() {
  const res = await fetch("/_api/grpc");
  const handlers = await res.json();
  const el = document.getElementById("grpc-handlers");
  const empty = document.getElementById("grpc-empty");

  if (handlers.length === 0) { el.innerHTML = ""; empty.style.display = ""; return; }
  empty.style.display = "none";

  el.innerHTML = handlers.map(h => {
    const desc = h.description ? '<span class="desc">' + escapeHtml(h.description) + '</span>' : '';
    return '<div class="route">' +
      '<div class="route-info">' +
        '<span class="method grpc">gRPC</span>' +
        '<code>' + escapeHtml(h.service + "." + h.method) + '</code>' +
        desc +
        (h.delay ? '<span class="status">' + h.delay + 'ms</span>' : '') +
      '</div>' +
      '<div class="route-actions">' +
        '<button class="copy" onclick="copyGRPCCurl(\'' + escapeHtml(h.service) + '\', \'' + escapeHtml(h.method) + '\')" title="Copy curl">📋</button>' +
        '<button class="del" onclick="deleteGRPCHandler(\'' + h.id + '\')" title="Delete">✕</button>' +
      '</div>' +
    '</div>';
  }).join("");
}

function openGRPCModal() {
  document.getElementById("grpc-service").value = "";
  document.getElementById("grpc-method").value = "";
  document.getElementById("grpc-desc").value = "";
  document.getElementById("grpc-delay").value = "0";
  document.getElementById("grpc-response").value = "";
  document.getElementById("grpc-modal").style.display = "";
}

function closeGRPCModal() { document.getElementById("grpc-modal").style.display = "none"; }

async function saveGRPCHandler() {
  let response = {};
  const respText = document.getElementById("grpc-response").value.trim();
  if (respText) {
    try { response = JSON.parse(respText); } catch(e) { alert("Invalid JSON response"); return; }
  }

  const h = {
    service: document.getElementById("grpc-service").value.trim(),
    method: document.getElementById("grpc-method").value.trim(),
    description: document.getElementById("grpc-desc").value.trim(),
    delay: parseInt(document.getElementById("grpc-delay").value) || 0,
    mock_response: response,
  };

  if (!h.service || !h.method) { alert("Service and Method are required"); return; }

  await fetch("/_api/grpc", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(h),
  });

  closeGRPCModal();
  loadGRPCHandlers();
}

async function deleteGRPCHandler(id) {
  if (!confirm("Delete this gRPC mock?")) return;
  await fetch("/_api/grpc?id=" + id, { method: "DELETE" });
  loadGRPCHandlers();
}

function copyGRPCCurl(service, method) {
  const cmd = "grpcurl -plaintext -d '{}' localhost:9099 " + service + "/" + method;
  navigator.clipboard.writeText(cmd);
  alert("Copied: " + cmd);
}

// --- Import Proto ---
function openProtoModal() {
  document.getElementById("proto-input").value = "";
  document.getElementById("proto-modal").style.display = "";
}

function closeProtoModal() { document.getElementById("proto-modal").style.display = "none"; }

async function importProto() {
  const input = document.getElementById("proto-input").value.trim();
  if (!input) { alert("Please paste proto content"); return; }
  
  const res = await fetch("/_api/import-proto", {
    method: "POST",
    headers: { "Content-Type": "text/plain" },
    body: input,
  });
  
  if (!res.ok) { const err = await res.text(); alert("Import failed: " + err); return; }
  
  const data = await res.json();
  alert("Imported " + data.imported + " gRPC methods!");
  closeProtoModal();
  loadGRPCHandlers();
}

// --- Theme Toggle ---
function toggleTheme() {
  const html = document.documentElement;
  const btn = document.querySelector('.theme-toggle');
  const isLight = html.classList.contains('light');
  
  if (isLight) {
    html.classList.remove('light');
    btn.textContent = '🌙';
    localStorage.setItem('mockapi-theme', 'dark');
  } else {
    html.classList.add('light');
    btn.textContent = '☀️';
    localStorage.setItem('mockapi-theme', 'light');
  }
}

// Initialize theme from localStorage
(function() {
  const saved = localStorage.getItem('mockapi-theme');
  const btn = document.querySelector('.theme-toggle');
  if (saved === 'light') {
    document.documentElement.classList.add('light');
    if (btn) btn.textContent = '☀️';
  }
})();

loadRoutes();
`