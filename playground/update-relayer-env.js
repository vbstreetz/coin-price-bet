const fs = require('fs');
const readline = require('readline');
const logfile = __dirname + '/relayer-create.log';

const re = new RegExp('Channel created: \\[band-consumer\\]chan\{(\\w+)\}port\{coinpricebet\} -> \\[ibc-bandchain\\]chan\{(\\w+)\}port\{oracle\}');

readline.createInterface({
    input: fs.createReadStream(logfile)
}).on('line', function (line) {
    const m = line.match(re);
    if (m) {
        const betchain_oracle_channel = m[1];
        const bandchain_oracle_channel = m[2];

        const f = __dirname + '/relayers.json';
        const j = JSON.parse(fs.readFileSync(f, 'utf8'));
        j.betchain_oracle_channel = betchain_oracle_channel;
        j.bandchain_oracle_channel = bandchain_oracle_channel;
        fs.writeFileSync(f, JSON.stringify(j, null, 2), 'utf8');

        let e = '';
        Object.entries(j).forEach(([k, v]) => {
            e += `${k}=${v}\n`;
        });
        fs.writeFileSync(__dirname + '/relayers.env', e, 'utf8');
    }
});

