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

  Promise.all(pauses)
    .then((resp) => console.log(resp.length))
    .catch(err => console.log(err))


  const continues = []
  for (let i = 1; i <= 11; i++) {
    continues.push(fetch("http://localhost:8080/continue/" + i))
    // await sleep(20)
  }

  Promise.all(continues)
    .then((resp) => console.log(resp.length))
    .catch(err => console.log(err))

})()