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
      <div class="title">
        <h1>🦞 MockAPI</h1>
        <p>Lightweight local API mock server</p>
      </div>
      <button class="theme-toggle" onclick="toggleTheme()" title="Toggle theme">🌙</button>
    </header>

    <div class="tabs">
      <button class="tab active" data-tab="routes">Routes</button>
      <button class="tab" data-tab="graphql">GraphQL</button>
      <button class="tab" data-tab="grpc">gRPC</button>
      <button class="tab" data-tab="ws">WebSocket</button>
      <button class="tab" data-tab="logs">Request Log</button>
      <button class="tab" data-tab="settings">Settings</button>
    </div>

    <!-- Routes Tab -->
    <div id="tab-routes" class="tab-content active">
      <div class="toolbar">
        <div class="toolbar-left">
          <button onclick="openModal()">+ Add Route</button>
          <button class="secondary" onclick="openTemplates()">📋 Templates</button>
          <button class="secondary" onclick="openSwaggerModal()">📥 Swagger</button>
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

    <!-- WebSocket Tab -->
    <div id="tab-ws" class="tab-content">
      <div class="toolbar">
        <div class="toolbar-left">
          <button onclick="openWSModal()">+ Add WS Handler</button>
        </div>
      </div>
      <div id="ws-handlers" class="routes"></div>
      <div id="ws-empty" class="empty" style="display:none;">
        <p>No WebSocket handlers. Click <strong>+ Add WS Handler</strong> to create one.</p>
      </div>
    </div>

    <!-- GraphQL Tab -->
    <div id="tab-graphql" class="tab-content">
      <div class="toolbar">
        <div class="toolbar-left">
          <button onclick="openGraphQLModal()">+ Add GraphQL Mock</button>
        </div>
        <span class="hint">Endpoint: <code id="graphql-url"></code></span>
      </div>
      <div id="graphql-handlers" class="routes"></div>
      <div id="graphql-empty" class="empty" style="display:none;">
        <p>No GraphQL mocks. Click <strong>+ Add GraphQL Mock</strong> to create one.</p>
      </div>
    </div>

    <!-- gRPC Tab -->
    <div id="tab-grpc" class="tab-content">
      <div class="toolbar">
        <div class="toolbar-left">
          <button onclick="openGRPCModal()">+ Add gRPC Mock</button>
          <button class="secondary" onclick="openProtoModal()">📥 Import Proto</button>
        </div>
        <span class="hint">Endpoint: <code id="grpc-url"></code></span>
      </div>
      <div id="grpc-handlers" class="routes"></div>
      <div id="grpc-empty" class="empty" style="display:none;">
        <p>No gRPC mocks. Click <strong>+ Add gRPC Mock</strong> or <strong>Import Proto</strong>.</p>
      </div>
    </div>

    <!-- Logs Tab -->
    <div id="tab-logs" class="tab-content">
      <div class="toolbar">
        <span class="hint" id="log-count"></span>
        <button class="secondary" onclick="clearLogs()">🗑 Clear</button>
        <button class="secondary" onclick="loadLogs()">↻ Refresh</button>
      </div>
      <div id="logs" class="logs"></div>
      <div id="logs-empty" class="empty" style="display:none;">
        <p>No requests logged yet. Hit a mock endpoint to see activity.</p>
      </div>
    </div>

    <!-- Settings Tab -->
    <div id="tab-settings" class="tab-content">
      <div class="settings">
        <div class="setting-group">
          <h3>🌐 Proxy Mode</h3>
          <p class="setting-desc">Forward unmatched requests to a real backend. Mock specific routes while proxying the rest.</p>
          <div class="setting-row">
            <label>Backend URL</label>
            <input id="s-proxy" placeholder="http://localhost:3000" style="width:400px" />
          </div>
          <button onclick="saveSettings()" class="save-btn">Save Settings</button>
        </div>
        <div class="setting-group">
          <h3>🔧 General</h3>
          <div class="setting-row">
            <label>CORS</label>
            <input type="checkbox" id="s-cors" checked />
          </div>
          <div class="setting-row">
            <label>Max Logs</label>
            <input type="number" id="s-maxlogs" value="500" style="width:100px" />
          </div>
          <button onclick="saveSettings()" class="save-btn">Save Settings</button>
        </div>
      </div>
    </div>

    <div class="footer">
      Mock base: <code id="mock-url"></code> | WS: <code id="ws-url"></code>
    </div>
  </div>

  <!-- Add/Edit Route Modal -->
  <div id="modal" class="modal" style="display:none;">
    <div class="modal-content wide">
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
          <label>Response Type</label>
          <select id="f-type" onchange="toggleResponseType()">
            <option value="static">Static JSON</option>
            <option value="script">JavaScript (Dynamic)</option>
          </select>
        </div>
        
        <div id="static-response">
          <div class="row full">
            <label>Response Body</label>
            <textarea id="f-body" rows="6" placeholder='{"id": 1}&#10;Use {{param}} for path params'></textarea>
          </div>
          <div class="row full">
            <label>Custom Headers (optional)</label>
            <textarea id="f-headers" rows="2" placeholder='{"X-Token": "abc"}'></textarea>
          </div>
        </div>
        
        <div id="script-response" style="display:none;">
          <div class="row full">
            <label>JavaScript Script</label>
            <textarea id="f-script" rows="10" placeholder='// Available: method, path, headers, body, params, query
// Use respond({status, body, headers}) to return
var id = params.id || "unknown";
respond({
  status: 200,
  body: JSON.stringify({id: id, time: Date.now()})
});'></textarea>
          </div>
        </div>

        <div class="row full conditions-header">
          <label>⚡ Conditional Matching (optional)</label>
          <p class="help-text">Only respond when conditions are met. Otherwise falls through to next route or proxy/404.</p>
        </div>
        <div class="row full">
          <label>Match Headers</label>
          <textarea id="f-match-headers" rows="2" placeholder='{"Authorization": "Bearer test"}'></textarea>
        </div>
        <div class="row full">
          <label>Match Body (substring)</label>
          <input id="f-match-body" placeholder='e.g. "error" to match error responses' />
        </div>
      </div>
      <div class="actions">
        <button onclick="saveRoute()">Save</button>
        <button class="secondary" onclick="closeModal()">Cancel</button>
      </div>
    </div>
  </div>

  <!-- Swagger Import Modal -->
  <div id="swagger-modal" class="modal" style="display:none;">
    <div class="modal-content">
      <h2>📥 Import Swagger/OpenAPI</h2>
      <p class="template-hint">Paste your OpenAPI 2.0 or 3.x spec (JSON or YAML) to auto-generate mock routes.</p>
      <div class="form">
        <div class="row full">
          <label>OpenAPI Spec</label>
          <textarea id="swagger-input" rows="15" placeholder='{
  "openapi": "3.0.0",
  "paths": {
    "/users": {
      "get": { ... }
    }
  }
}'></textarea>
        </div>
        <div class="row full">
          <label>Or upload file</label>
          <input type="file" id="swagger-file" accept=".json,.yaml,.yml" onchange="loadSwaggerFile(event)" />
        </div>
      </div>
      <div class="actions">
        <button onclick="importSwagger()">Import</button>
        <button class="secondary" onclick="closeSwaggerModal()">Cancel</button>
      </div>
    </div>
  </div>

  <!-- WebSocket Modal -->
  <div id="ws-modal" class="modal" style="display:none;">
    <div class="modal-content">
      <h2>Add WebSocket Handler</h2>
      <div class="form">
        <div class="row">
          <label>Path</label>
          <input id="ws-path" placeholder="/chat" />
        </div>
        <div class="row full">
          <label>Description</label>
          <input id="ws-desc" placeholder="Optional description" />
        </div>
        <div class="row">
          <label>Delay (ms)</label>
          <input id="ws-delay" type="number" value="0" />
        </div>
        <div class="row full">
          <label>Auto Reply (JSON)</label>
          <textarea id="ws-auto-reply" rows="3" placeholder='{"type": "echo", "data": "received"}'></textarea>
        </div>
        <div class="row full">
          <label>On Connect Message</label>
          <textarea id="ws-on-connect" rows="2" placeholder='{"type": "connected"}'></textarea>
        </div>
        <div class="row full">
          <label>On Message Script (JS, optional)</label>
          <textarea id="ws-on-message" rows="5" placeholder='// body = incoming message
respond({body: JSON.stringify({echo: body})});'></textarea>
        </div>
      </div>
      <div class="actions">
        <button onclick="saveWSHandler()">Save</button>
        <button class="secondary" onclick="closeWSModal()">Cancel</button>
      </div>
    </div>
  </div>

  <!-- GraphQL Modal -->
  <div id="graphql-modal" class="modal" style="display:none;">
    <div class="modal-content">
      <h2>Add GraphQL Mock</h2>
      <div class="form">
        <div class="row full">
          <label>Operation Name</label>
          <input id="gql-op" placeholder="GetUsers, CreateUser, or leave empty for catch-all" />
        </div>
        <div class="row full">
          <label>Description</label>
          <input id="gql-desc" placeholder="Optional description" />
        </div>
        <div class="row">
          <label>Delay (ms)</label>
          <input id="gql-delay" type="number" value="0" />
        </div>
        <div class="row full">
          <label>Response Data (JSON)</label>
          <textarea id="gql-response" rows="8" placeholder='{"users": [{"id": 1, "name": "Alice"}]}'></textarea>
        </div>
      </div>
      <div class="actions">
        <button onclick="saveGraphQLHandler()">Save</button>
        <button class="secondary" onclick="closeGraphQLModal()">Cancel</button>
      </div>
    </div>
  </div>

  <!-- gRPC Modal -->
  <div id="grpc-modal" class="modal" style="display:none;">
    <div class="modal-content">
      <h2>Add gRPC Mock</h2>
      <div class="form">
        <div class="row">
          <label>Service</label>
          <input id="grpc-service" placeholder="UserService" />
        </div>
        <div class="row">
          <label>Method</label>
          <input id="grpc-method" placeholder="GetUser" />
        </div>
        <div class="row full">
          <label>Description</label>
          <input id="grpc-desc" placeholder="Optional description" />
        </div>
        <div class="row">
          <label>Delay (ms)</label>
          <input id="grpc-delay" type="number" value="0" />
        </div>
        <div class="row full">
          <label>Mock Response (JSON)</label>
          <textarea id="grpc-response" rows="8" placeholder='{"id": 1, "name": "John"}'></textarea>
        </div>
      </div>
      <div class="actions">
        <button onclick="saveGRPCHandler()">Save</button>
        <button class="secondary" onclick="closeGRPCModal()">Cancel</button>
      </div>
    </div>
  </div>

  <!-- Import Proto Modal -->
  <div id="proto-modal" class="modal" style="display:none;">
    <div class="modal-content">
      <h2>📥 Import Proto File</h2>
      <p class="template-hint">Paste your .proto file to auto-generate gRPC mock handlers.</p>
      <div class="form">
        <div class="row full">
          <label>Proto Content</label>
          <textarea id="proto-input" rows="15" placeholder='syntax = "proto3";

service UserService {
  rpc GetUser(GetUserRequest) returns (User);
  rpc CreateUser(CreateUserRequest) returns (User);
}'></textarea>
        </div>
      </div>
      <div class="actions">
        <button onclick="importProto()">Import</button>
        <button class="secondary" onclick="closeProtoModal()">Cancel</button>
      </div>
    </div>
  </div>

  <!-- Templates Modal -->
  <div id="templates-modal" class="modal" style="display:none;">
    <div class="modal-content wide">
      <h2>📋 Quick Templates</h2>
      <p class="template-hint">Click to add. Includes REST CRUD, Auth, conditional responses, scripts, and more.</p>
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