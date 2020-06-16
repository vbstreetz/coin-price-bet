const fs = require("fs")
const readline = require("readline")
const logfile = __dirname + "/relayer-create.log"

// const re1 = new RegExp(
//   "Channel created: \\[band-consumer\\]chan{(\\w+)}port{transfer} -> \\[band-cosmoshub\\]chan{(\\w+)}port{transfer}"
// )
const re2 = new RegExp(
  "Channel created: \\[band-consumer\\]chan{(\\w+)}port{coinpricebet} -> \\[ibc-bandchain\\]chan{(\\w+)}port{oracle}"
)

const env = {}

readline
  .createInterface({
    input: fs.createReadStream(logfile),
  })
  .on("line", function (line) {
    // const m1 = line.match(re1)
    // if (m1) {
    //   env.betchain_transfer_channel = m1[1]
    //   env.gaia_transfer_channel = m1[2]
    // }

    const m2 = line.match(re2)
    if (m2) {
      env.betchain_oracle_channel = m2[1]
      env.bandchain_oracle_channel = m2[2]
    }
  })
  .on("close", function () {
    let e = ""
    Object.entries(env).forEach(([k, v]) => {
      e += `${k}=${v}\n`
    })
    fs.writeFileSync(__dirname + "/relayers.env", e, "utf8")
  })
