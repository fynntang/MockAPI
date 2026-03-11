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

    <div class="toolbar">
      <button onclick="openModal()">+ Add Route</button>
      <span class="hint">Mock base URL: <code id="mock-url"></code></span>
    </div>

    <div id="routes" class="routes"></div>

    <div id="empty" class="empty" style="display:none;">
      <p>No routes yet. Click <strong>+ Add Route</strong> to create one.</p>
    </div>
  </div>

  <div id="modal" class="modal" style="display:none;">
    <div class="modal-content">
      <h2>Add Mock Route</h2>
      <div class="form">
        <div class="row">
          <label>Method</label>
          <select id="f-method">
            <option>GET</option><option>POST</option><option>PUT</option><option>PATCH</option><option>DELETE</option>
          </select>
        </div>
        <div class="row">
          <label>Path</label>
          <input id="f-path" placeholder="/users/:id" />
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
          <label>Response Body (JSON)</label>
          <textarea id="f-body" rows="8" placeholder='{"id": 1, "name": "test"}'></textarea>
        </div>
        <div class="row full">
          <label>Custom Headers (JSON)</label>
          <textarea id="f-headers" rows="3" placeholder='{"X-Custom": "value"}'></textarea>
        </div>
      </div>
      <div class="actions">
        <button onclick="addRoute()">Create</button>
        <button class="secondary" onclick="closeModal()">Cancel</button>
      </div>
    </div>
  </div>

  <script src="/web/app.js"></script>
</body>
</html>
`
