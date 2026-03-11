package server

const indexHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>MockAPI — Local API Mock Server</title>
  <link rel="stylesheet" href="/web/style.css">
</head>
<body>
  <div class="app">
    <header>
      <h1>🦞 MockAPI</h1>
      <p>Lightweight local API mock server</p>
    </header>

    <!-- Tabs -->
    <div class="tabs">
      <button class="tab active" data-tab="routes">Routes</button>
      <button class="tab" data-tab="logs">Request Log</button>
    </div>

    <!-- Routes Tab -->
    <div id="tab-routes" class="tab-content active">
      <div class="toolbar">
        <div class="toolbar-left">
          <button onclick="openModal()">+ Add Route</button>
          <button class="secondary" onclick="openTemplates()">📋 Templates</button>
          <button class="secondary" onclick="exportRoutes()">⬇ Export</button>
          <label class="import-btn">
            ⬆ Import
            <input type="file" accept=".json" onchange="importRoutes(event)" hidden>
          </label>
        </div>
        <input class="search" id="route-search" placeholder="🔍 Filter routes..." oninput="loadRoutes()" />
      </div>

      <div id="routes" class="routes"></div>
      <div id="empty" class="empty" style="display:none;">
        <p>No routes yet. Click <strong>+ Add Route</strong> or use a <strong>Template</strong>.</p>
      </div>
    </div>

    <!-- Logs Tab -->
    <div id="tab-logs" class="tab-content">
      <div class="toolbar">
        <span class="hint" id="log-count"></span>
        <button class="secondary" onclick="clearLogs()">🗑 Clear Logs</button>
        <button class="secondary" onclick="loadLogs()">↻ Refresh</button>
      </div>
      <div id="logs" class="logs"></div>
      <div id="logs-empty" class="empty" style="display:none;">
        <p>No requests logged yet. Hit a mock endpoint to see activity.</p>
      </div>
    </div>

    <div class="footer">
      Mock base: <code id="mock-url"></code>
    </div>
  </div>

  <!-- Add/Edit Route Modal -->
  <div id="modal" class="modal" style="display:none;">
    <div class="modal-content">
      <h2 id="modal-title">Add Mock Route</h2>
      <div class="form">
        <div class="row">
          <label>Method</label>
          <select id="f-method">
            <option>GET</option><option>POST</option><option>PUT</option>
            <option>PATCH</option><option>DELETE</option><option>ALL</option>
          </select>
        </div>
        <div class="row">
          <label>Path</label>
          <input id="f-path" placeholder="/users/:id" />
          <span class="help-text">Use :param for dynamic, * for wildcard</span>
        </div>
        <div class="row">
          <label>Status</label>
          <input id="f-status" type="number" value="200" />
        </div>
        <div class="row">
          <label>Delay (ms)</label>
          <input id="f-delay" type="number" value="0" placeholder="0" />
        </div>
        <div class="row full">
          <label>Description</label>
          <input id="f-desc" placeholder="Optional description" />
        </div>
        <div class="row full">
          <label>Response Body</label>
          <textarea id="f-body" rows="8" placeholder='{"id": 1, "name": "test"}&#10;Use {{param}} for path params'></textarea>
        </div>
        <div class="row full">
          <label>Custom Headers (JSON, optional)</label>
          <textarea id="f-headers" rows="3" placeholder='{"X-Custom": "value"}'></textarea>
        </div>
      </div>
      <div class="actions">
        <button onclick="saveRoute()">Save</button>
        <button class="secondary" onclick="closeModal()">Cancel</button>
      </div>
    </div>
  </div>

  <!-- Templates Modal -->
  <div id="templates-modal" class="modal" style="display:none;">
    <div class="modal-content wide">
      <h2>📋 Quick Templates</h2>
      <p class="template-hint">Click to add. Supports REST API, Auth, Health Check, and more.</p>
      <div id="template-list" class="template-list"></div>
      <div class="actions">
        <button class="secondary" onclick="closeTemplates()">Close</button>
      </div>
    </div>
  </div>

  <script src="/web/app.js"></script>
</body>
</html>
`
