package server

const styleCSS = `
* { margin: 0; padding: 0; box-sizing: border-box; }

body {
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
  background: #0f1117;
  color: #c9d1d9;
  min-height: 100vh;
}

.app { max-width: 800px; margin: 0 auto; padding: 40px 24px; }

header { margin-bottom: 32px; }
header h1 { font-size: 28px; font-weight: 700; color: #f0f6fc; }
header p { color: #8b949e; margin-top: 4px; }

.toolbar {
  display: flex; align-items: center; gap: 16px; margin-bottom: 24px;
}
.toolbar button {
  background: #238636; color: #fff; border: none; padding: 8px 16px;
  border-radius: 6px; font-size: 14px; cursor: pointer; font-weight: 500;
}
.toolbar button:hover { background: #2ea043; }
.hint { color: #8b949e; font-size: 13px; }
.hint code { background: #161b22; padding: 2px 8px; border-radius: 4px; color: #79c0ff; }

.routes { display: flex; flex-direction: column; gap: 8px; }

.route {
  display: flex; justify-content: space-between; align-items: center;
  background: #161b22; border: 1px solid #30363d; border-radius: 8px; padding: 12px 16px;
}
.route-info { display: flex; align-items: center; gap: 12px; }
.route-actions { display: flex; gap: 8px; }

.method {
  display: inline-block; padding: 2px 8px; border-radius: 4px;
  font-size: 11px; font-weight: 700; text-transform: uppercase;
}
.method.get { background: #1f6feb33; color: #58a6ff; }
.method.post { background: #23863633; color: #3fb950; }
.method.put { background: #d2992233; color: #d29922; }
.method.patch { background: #d2992233; color: #d29922; }
.method.delete { background: #f8514933; color: #f85149; }

code { font-size: 13px; color: #c9d1d9; }
.status { font-size: 12px; color: #8b949e; }

button.copy {
  background: transparent; border: 1px solid #30363d; color: #8b949e;
  padding: 4px 10px; border-radius: 4px; font-size: 12px; cursor: pointer;
}
button.copy:hover { border-color: #58a6ff; color: #58a6ff; }
button.del {
  background: transparent; border: 1px solid #30363d; color: #f85149;
  padding: 4px 8px; border-radius: 4px; cursor: pointer; font-size: 12px;
}
button.del:hover { background: #f8514922; }

.empty {
  text-align: center; padding: 60px 20px; color: #8b949e;
  border: 1px dashed #30363d; border-radius: 8px;
}

/* Modal */
.modal {
  position: fixed; inset: 0; background: rgba(0,0,0,0.6);
  display: flex; align-items: center; justify-content: center; z-index: 100;
}
.modal-content {
  background: #161b22; border: 1px solid #30363d; border-radius: 12px;
  padding: 24px; width: 520px; max-height: 90vh; overflow-y: auto;
}
.modal-content h2 { font-size: 18px; margin-bottom: 20px; color: #f0f6fc; }
.form { display: grid; grid-template-columns: 1fr 1fr; gap: 12px; }
.row.full { grid-column: 1 / -1; }
.row label { display: block; font-size: 12px; color: #8b949e; margin-bottom: 4px; }
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
`
