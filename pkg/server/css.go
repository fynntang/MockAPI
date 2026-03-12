package server

const styleCSS = `
* { margin: 0; padding: 0; box-sizing: border-box; }

:root {
  --bg-primary: #0f1117;
  --bg-secondary: #161b22;
  --bg-tertiary: #0d1117;
  --border-color: #30363d;
  --border-light: #21262d;
  --text-primary: #f0f6fc;
  --text-secondary: #c9d1d9;
  --text-muted: #8b949e;
  --text-faint: #484f58;
  --accent-blue: #58a6ff;
  --accent-green: #3fb950;
  --accent-orange: #d29922;
  --accent-purple: #a371f7;
  --accent-red: #f85149;
  --btn-primary: #238636;
  --btn-primary-hover: #2ea043;
}

.light {
  --bg-primary: #ffffff;
  --bg-secondary: #f6f8fa;
  --bg-tertiary: #ffffff;
  --border-color: #d0d7de;
  --border-light: #d8dee4;
  --text-primary: #24292f;
  --text-secondary: #57606a;
  --text-muted: #6e7781;
  --text-faint: #8c959f;
  --accent-blue: #0969da;
  --accent-green: #1a7f37;
  --accent-orange: #bf8700;
  --accent-purple: #8250df;
  --accent-red: #cf222e;
  --btn-primary: #1f883d;
  --btn-primary-hover: #1a7f37;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
  background: var(--bg-primary);
  color: var(--text-secondary);
  min-height: 100vh;
}

.app { max-width: 860px; margin: 0 auto; padding: 40px 24px; }

header { margin-bottom: 24px; display: flex; justify-content: space-between; align-items: flex-start; }
header .title { }
header h1 { font-size: 28px; font-weight: 700; color: var(--text-primary); }
header p { color: var(--text-muted); margin-top: 4px; }

/* Theme Toggle */
.theme-toggle {
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 6px;
  padding: 6px 10px;
  cursor: pointer;
  font-size: 18px;
  transition: all 0.15s;
}
.theme-toggle:hover {
  border-color: var(--accent-blue);
}

/* Tabs */
.tabs { display: flex; gap: 0; margin-bottom: 24px; border-bottom: 1px solid var(--border-light); }
.tab {
  background: none; border: none; color: var(--text-muted); padding: 10px 20px;
  font-size: 14px; cursor: pointer; border-bottom: 2px solid transparent;
  transition: all .15s;
}
.tab:hover { color: var(--text-secondary); }
.tab.active { color: var(--text-primary); border-bottom-color: var(--accent-blue); }
.tab-content { display: none; }
.tab-content.active { display: block; }

/* Toolbar */
.toolbar {
  display: flex; align-items: center; gap: 10px; margin-bottom: 20px; flex-wrap: wrap;
}
.toolbar-left { display: flex; gap: 8px; align-items: center; }
.toolbar button {
  background: var(--btn-primary); color: #fff; border: none; padding: 7px 14px;
  border-radius: 6px; font-size: 13px; cursor: pointer; font-weight: 500;
  white-space: nowrap;
}
.toolbar button:hover { background: var(--btn-primary-hover); }
.toolbar button.secondary {
  background: transparent; border: 1px solid var(--border-color); color: var(--text-muted);
}
.toolbar button.secondary:hover { border-color: var(--accent-blue); color: var(--accent-blue); }

.search {
  flex: 1; min-width: 150px; background: var(--bg-tertiary); border: 1px solid var(--border-color);
  border-radius: 6px; padding: 7px 12px; color: var(--text-secondary); font-size: 13px;
}
.search:focus { border-color: var(--accent-blue); outline: none; }
.search::placeholder { color: var(--text-faint); }

.hint { color: var(--text-muted); font-size: 13px; }
.hint code { background: var(--bg-secondary); padding: 2px 8px; border-radius: 4px; color: var(--accent-blue); }

/* Routes */
.routes { display: flex; flex-direction: column; gap: 6px; }

.route {
  display: flex; justify-content: space-between; align-items: center;
  background: var(--bg-secondary); border: 1px solid var(--border-color); border-radius: 8px; padding: 10px 16px;
  transition: border-color .15s;
}
.route:hover { border-color: var(--text-faint); }
.route-info { display: flex; align-items: center; gap: 10px; flex-wrap: wrap; }
.route-actions { display: flex; gap: 4px; }

.method {
  display: inline-block; padding: 2px 8px; border-radius: 4px;
  font-size: 11px; font-weight: 700; text-transform: uppercase;
}
.method.get { background: rgba(9,105,218,0.15); color: var(--accent-blue); }
.method.post { background: rgba(26,127,55,0.15); color: var(--accent-green); }
.method.put { background: rgba(191,135,0,0.15); color: var(--accent-orange); }
.method.patch { background: rgba(191,135,0,0.15); color: var(--accent-orange); }
.method.delete { background: rgba(207,34,46,0.15); color: var(--accent-red); }
.method.all { background: rgba(130,80,223,0.15); color: var(--accent-purple); }
.method.ws { background: rgba(9,105,218,0.15); color: var(--accent-blue); }
.method.gql { background: rgba(233,30,99,0.15); color: #e91e63; }
.method.grpc { background: rgba(103,58,183,0.15); color: #673ab7; }

code { font-size: 13px; color: var(--text-secondary); }
.status { font-size: 12px; color: var(--text-muted); }
.desc { font-size: 11px; color: var(--text-faint); max-width: 200px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }

button.copy, button.edit {
  background: transparent; border: 1px solid var(--border-color); color: var(--text-muted);
  padding: 3px 8px; border-radius: 4px; font-size: 12px; cursor: pointer;
}
button.copy:hover { border-color: var(--accent-blue); color: var(--accent-blue); }
button.del {
  background: transparent; border: 1px solid var(--border-color); color: var(--accent-red);
  padding: 3px 8px; border-radius: 4px; cursor: pointer; font-size: 12px;
}
button.del:hover { background: rgba(207,34,46,0.1); }

.empty {
  text-align: center; padding: 60px 20px; color: var(--text-muted);
  border: 1px dashed var(--border-color); border-radius: 8px;
}

.import-btn {
  display: inline-flex; align-items: center; background: transparent;
  border: 1px solid var(--border-color); color: var(--text-muted); padding: 7px 14px;
  border-radius: 6px; font-size: 13px; cursor: pointer; font-weight: 500;
  transition: all .15s;
}
.import-btn:hover { border-color: var(--accent-blue); color: var(--accent-blue); }

/* Logs */
.logs { display: flex; flex-direction: column; gap: 4px; }
.log-item {
  display: flex; align-items: center; gap: 10px;
  background: var(--bg-secondary); border: 1px solid var(--border-light); border-radius: 6px; padding: 8px 14px;
  font-size: 12px;
}
.log-time { color: var(--text-faint); font-family: "SF Mono", Monaco, monospace; min-width: 65px; }
.log-status { font-weight: 600; min-width: 30px; text-align: center; }
.log-status.ok { color: var(--accent-green); }
.log-status.err { color: var(--accent-red); }
.log-delay { color: var(--text-faint); margin-left: auto; font-family: "SF Mono", Monaco, monospace; }

/* Footer */
.footer {
  margin-top: 40px; padding-top: 20px; border-top: 1px solid var(--border-light);
  color: var(--text-faint); font-size: 13px; text-align: center;
}
.footer code { background: var(--bg-secondary); padding: 2px 8px; border-radius: 4px; color: var(--accent-blue); }

/* Modal */
.modal {
  position: fixed; inset: 0; background: rgba(0,0,0,0.6);
  display: flex; align-items: center; justify-content: center; z-index: 100;
}
.modal-content {
  background: var(--bg-secondary); border: 1px solid var(--border-color); border-radius: 12px;
  padding: 24px; width: 560px; max-height: 90vh; overflow-y: auto;
}
.modal-content.wide { width: 700px; }
.modal-content h2 { font-size: 18px; margin-bottom: 16px; color: var(--text-primary); }
.form { display: grid; grid-template-columns: 1fr 1fr; gap: 12px; }
.row.full { grid-column: 1 / -1; }
.row label { display: block; font-size: 12px; color: var(--text-muted); margin-bottom: 4px; }
.help-text { display: block; font-size: 10px; color: var(--text-faint); margin-top: 2px; }
.row input, .row select, .row textarea {
  width: 100%; background: var(--bg-tertiary); border: 1px solid var(--border-color);
  border-radius: 6px; padding: 8px 12px; color: var(--text-secondary); font-size: 14px;
  font-family: inherit;
}
.row textarea { font-family: "SF Mono", Monaco, monospace; font-size: 13px; }
.actions { margin-top: 20px; display: flex; gap: 8px; }
.actions button {
  padding: 8px 20px; border-radius: 6px; border: none; cursor: pointer;
  font-size: 14px; font-weight: 500;
}
.actions button:first-child { background: var(--btn-primary); color: #fff; }
.actions button:first-child:hover { background: var(--btn-primary-hover); }
.actions button.secondary { background: transparent; border: 1px solid var(--border-color); color: var(--text-muted); }

/* Templates */
.template-hint { font-size: 13px; color: var(--text-muted); margin-bottom: 16px; }
.template-list { display: grid; grid-template-columns: 1fr 1fr; gap: 8px; max-height: 55vh; overflow-y: auto; }
.template-item {
  background: var(--bg-tertiary); border: 1px solid var(--border-color); border-radius: 8px;
  padding: 12px; cursor: pointer; transition: all .15s;
}
.template-item:hover { border-color: var(--accent-blue); background: var(--bg-secondary); }
.template-top { display: flex; align-items: center; gap: 8px; margin-bottom: 4px; }
.template-desc { font-size: 12px; color: var(--text-muted); }

/* Settings */
.settings { max-width: 600px; }
.setting-group {
  background: var(--bg-secondary); border: 1px solid var(--border-color); border-radius: 8px;
  padding: 20px; margin-bottom: 16px;
}
.setting-group h3 { font-size: 16px; margin-bottom: 8px; color: var(--text-primary); }
.setting-desc { font-size: 13px; color: var(--text-muted); margin-bottom: 16px; }
.setting-row {
  display: flex; align-items: center; gap: 12px; margin-bottom: 12px;
}
.setting-row label { min-width: 100px; font-size: 14px; color: var(--text-secondary); }
.setting-row input[type="text"], .setting-row input[type="number"] {
  background: var(--bg-tertiary); border: 1px solid var(--border-color); border-radius: 6px;
  padding: 8px 12px; color: var(--text-secondary); font-size: 14px;
}
.save-btn {
  background: var(--btn-primary); color: #fff; border: none; padding: 7px 16px;
  border-radius: 6px; font-size: 13px; cursor: pointer; margin-top: 8px;
}
.save-btn:hover { background: var(--btn-primary-hover); }

/* Condition badge */
.condition-badge {
  display: inline-block; background: rgba(130,80,223,0.15); color: var(--accent-purple);
  font-size: 10px; padding: 1px 6px; border-radius: 3px; font-weight: 700;
}
.script-badge {
  display: inline-block; background: rgba(240,136,62,0.15); color: #f0883e;
  font-size: 10px; padding: 1px 6px; border-radius: 3px; font-weight: 700;
}
.proxy-badge {
  display: inline-block; background: rgba(9,105,218,0.15); color: var(--accent-blue);
  font-size: 10px; padding: 1px 6px; border-radius: 3px; font-weight: 700;
}
.conditions-header { margin-top: 8px; padding-top: 8px; border-top: 1px solid var(--border-light); }
`
