package server

const styleCSS = `
* { margin: 0; padding: 0; box-sizing: border-box; }

body {
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
  background: #0f1117;
  color: #c9d1d9;
  min-height: 100vh;
}

.app { max-width: 860px; margin: 0 auto; padding: 40px 24px; }

header { margin-bottom: 24px; }
header h1 { font-size: 28px; font-weight: 700; color: #f0f6fc; }
header p { color: #8b949e; margin-top: 4px; }

/* Tabs */
.tabs { display: flex; gap: 0; margin-bottom: 24px; border-bottom: 1px solid #21262d; }
.tab {
  background: none; border: none; color: #8b949e; padding: 10px 20px;
  font-size: 14px; cursor: pointer; border-bottom: 2px solid transparent;
  transition: all .15s;
}
.tab:hover { color: #c9d1d9; }
.tab.active { color: #f0f6fc; border-bottom-color: #58a6ff; }
.tab-content { display: none; }
.tab-content.active { display: block; }

/* Toolbar */
.toolbar {
  display: flex; align-items: center; gap: 10px; margin-bottom: 20px; flex-wrap: wrap;
}
.toolbar-left { display: flex; gap: 8px; align-items: center; }
.toolbar button {
  background: #238636; color: #fff; border: none; padding: 7px 14px;
  border-radius: 6px; font-size: 13px; cursor: pointer; font-weight: 500;
  white-space: nowrap;
}
.toolbar button:hover { background: #2ea043; }
.toolbar button.secondary {
  background: transparent; border: 1px solid #30363d; color: #8b949e;
}
.toolbar button.secondary:hover { border-color: #58a6ff; color: #58a6ff; }

.search {
  flex: 1; min-width: 150px; background: #0d1117; border: 1px solid #30363d;
  border-radius: 6px; padding: 7px 12px; color: #c9d1d9; font-size: 13px;
}
.search:focus { border-color: #58a6ff; outline: none; }
.search::placeholder { color: #484f58; }

.hint { color: #8b949e; font-size: 13px; }
.hint code { background: #161b22; padding: 2px 8px; border-radius: 4px; color: #79c0ff; }

/* Routes */
.routes { display: flex; flex-direction: column; gap: 6px; }

.route {
  display: flex; justify-content: space-between; align-items: center;
  background: #161b22; border: 1px solid #30363d; border-radius: 8px; padding: 10px 16px;
  transition: border-color .15s;
}
.route:hover { border-color: #484f58; }
.route-info { display: flex; align-items: center; gap: 10px; flex-wrap: wrap; }
.route-actions { display: flex; gap: 4px; }

.method {
  display: inline-block; padding: 2px 8px; border-radius: 4px;
  font-size: 11px; font-weight: 700; text-transform: uppercase;
}
.method.get { background: #1f6feb33; color: #58a6ff; }
.method.post { background: #23863633; color: #3fb950; }
.method.put { background: #d2992233; color: #d29922; }
.method.patch { background: #d2992233; color: #d29922; }
.method.delete { background: #f8514933; color: #f85149; }
.method.all { background: #a371f733; color: #a371f7; }
.method.ws { background: #1f6feb33; color: #58a6ff; }
.method.gql { background: #e91e6333; color: #e91e63; }
.method.grpc { background: #673ab733; color: #673ab7; }

code { font-size: 13px; color: #c9d1d9; }
.status { font-size: 12px; color: #8b949e; }
.desc { font-size: 11px; color: #6e7681; max-width: 200px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }

button.copy, button.edit {
  background: transparent; border: 1px solid #30363d; color: #8b949e;
  padding: 3px 8px; border-radius: 4px; font-size: 12px; cursor: pointer;
}
button.copy:hover { border-color: #58a6ff; color: #58a6ff; }
button.del {
  background: transparent; border: 1px solid #30363d; color: #f85149;
  padding: 3px 8px; border-radius: 4px; cursor: pointer; font-size: 12px;
}
button.del:hover { background: #f8514922; }

.empty {
  text-align: center; padding: 60px 20px; color: #8b949e;
  border: 1px dashed #30363d; border-radius: 8px;
}

.import-btn {
  display: inline-flex; align-items: center; background: transparent;
  border: 1px solid #30363d; color: #8b949e; padding: 7px 14px;
  border-radius: 6px; font-size: 13px; cursor: pointer; font-weight: 500;
  transition: all .15s;
}
.import-btn:hover { border-color: #58a6ff; color: #58a6ff; }

/* Logs */
.logs { display: flex; flex-direction: column; gap: 4px; }
.log-item {
  display: flex; align-items: center; gap: 10px;
  background: #161b22; border: 1px solid #21262d; border-radius: 6px; padding: 8px 14px;
  font-size: 12px;
}
.log-time { color: #484f58; font-family: "SF Mono", Monaco, monospace; min-width: 65px; }
.log-status { font-weight: 600; min-width: 30px; text-align: center; }
.log-status.ok { color: #3fb950; }
.log-status.err { color: #f85149; }
.log-delay { color: #484f58; margin-left: auto; font-family: "SF Mono", Monaco, monospace; }

/* Footer */
.footer {
  margin-top: 40px; padding-top: 20px; border-top: 1px solid #21262d;
  color: #484f58; font-size: 13px; text-align: center;
}
.footer code { background: #161b22; padding: 2px 8px; border-radius: 4px; color: #79c0ff; }

/* Modal */
.modal {
  position: fixed; inset: 0; background: rgba(0,0,0,0.6);
  display: flex; align-items: center; justify-content: center; z-index: 100;
}
.modal-content {
  background: #161b22; border: 1px solid #30363d; border-radius: 12px;
  padding: 24px; width: 560px; max-height: 90vh; overflow-y: auto;
}
.modal-content.wide { width: 700px; }
.modal-content h2 { font-size: 18px; margin-bottom: 16px; color: #f0f6fc; }
.form { display: grid; grid-template-columns: 1fr 1fr; gap: 12px; }
.row.full { grid-column: 1 / -1; }
.row label { display: block; font-size: 12px; color: #8b949e; margin-bottom: 4px; }
.help-text { display: block; font-size: 10px; color: #484f58; margin-top: 2px; }
.row input, .row select, .row textarea {
  width: 100%; background: #0d1117; border: 1px solid #30363d;
  border-radius: 6px; padding: 8px 12px; color: #c9d1d9; font-size: 14px;
  font-family: inherit;
}
.row textarea { font-family: "SF Mono", Monaco, monospace; font-size: 13px; }
.actions { margin-top: 20px; display: flex; gap: 8px; }
.actions button {
  padding: 8px 20px; border-radius: 6px; border: none; cursor: pointer;
  font-size: 14px; font-weight: 500;
}
.actions button:first-child { background: #238636; color: #fff; }
.actions button:first-child:hover { background: #2ea043; }
.actions button.secondary { background: transparent; border: 1px solid #30363d; color: #8b949e; }

/* Templates */
.template-hint { font-size: 13px; color: #8b949e; margin-bottom: 16px; }
.template-list { display: grid; grid-template-columns: 1fr 1fr; gap: 8px; max-height: 55vh; overflow-y: auto; }
.template-item {
  background: #0d1117; border: 1px solid #30363d; border-radius: 8px;
  padding: 12px; cursor: pointer; transition: all .15s;
}
.template-item:hover { border-color: #58a6ff; background: #161b22; }
.template-top { display: flex; align-items: center; gap: 8px; margin-bottom: 4px; }
.template-desc { font-size: 12px; color: #8b949e; }

/* Settings */
.settings { max-width: 600px; }
.setting-group {
  background: #161b22; border: 1px solid #30363d; border-radius: 8px;
  padding: 20px; margin-bottom: 16px;
}
.setting-group h3 { font-size: 16px; margin-bottom: 8px; }
.setting-desc { font-size: 13px; color: #8b949e; margin-bottom: 16px; }
.setting-row {
  display: flex; align-items: center; gap: 12px; margin-bottom: 12px;
}
.setting-row label { min-width: 100px; font-size: 14px; color: #c9d1d9; }
.setting-row input[type="text"], .setting-row input[type="number"] {
  background: #0d1117; border: 1px solid #30363d; border-radius: 6px;
  padding: 8px 12px; color: #c9d1d9; font-size: 14px;
}
.save-btn {
  background: #238636; color: #fff; border: none; padding: 7px 16px;
  border-radius: 6px; font-size: 13px; cursor: pointer; margin-top: 8px;
}
.save-btn:hover { background: #2ea043; }

/* Condition badge */
.condition-badge {
  display: inline-block; background: #a371f733; color: #a371f7;
  font-size: 10px; padding: 1px 6px; border-radius: 3px; font-weight: 700;
}
.script-badge {
  display: inline-block; background: #f0883e33; color: #f0883e;
  font-size: 10px; padding: 1px 6px; border-radius: 3px; font-weight: 700;
}
.proxy-badge {
  display: inline-block; background: #1f6feb33; color: #58a6ff;
  font-size: 10px; padding: 1px 6px; border-radius: 3px; font-weight: 700;
}
.conditions-header { margin-top: 8px; padding-top: 8px; border-top: 1px solid #21262d; }
`
