const fetch = require("node-fetch");

(async () => {
  const continues = []
  for (let i = 1; i <= 11; i++) {
    continues.push(fetch("http://localhost:8080/continue/" + i))
  }

  Promise.all(continues)
    .then((resp) => console.log(resp.length))
    .catch(err => console.log(err))
})()