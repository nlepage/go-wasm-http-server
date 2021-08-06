const http = require("http");
const fs = require("fs");
const path = require("path");

const app = http.createServer((req, res) => {
  let url = req.url.replace(/\?.*$/, "");
  if (url === "/") {
    url = "index.html";
  }
  url = path.join(__dirname, url);
  console.log("读取文件", url);
  res.statusCode = 200;
  if (/\.js$/.test(url)) {
    res.setHeader("Content-Type", "text/javascript");
  }
  if (/\.wasm$/.test(url)) {
    res.setHeader("Content-Type", "application/wasm");
  }
  try {
    res.end(fs.readFileSync(url));
  } catch (e) {
  }
});

app.listen(3001);
