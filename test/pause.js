const fetch = require("node-fetch")
const pauses = []
const count = process.env.COUNT || 1000 

function sleep(ms) {
  return new Promise(resolve => setTimeout(resolve, ms));
}

(async () => {
  for (let i = 0; i < count; i++) {
    let id = Math.floor(Math.random() * 10) + 1
    pauses.push(fetch("http://localhost:8080/pause/" + id))
    await sleep(10)
  }

  console.log("added all requests")
  Promise.all(pauses)
    .then((resp) => console.log(resp.length))
    .catch(err => console.log(err))
})()
