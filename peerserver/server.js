const express = require("express");
const app = express();

app.set('view engine', 'ejs')

const server = require('http').Server(app);

const { ExpressPeerServer } = require("peer");

const peerServer = ExpressPeerServer(server, {
    debug: true,
    path: "/myapp",
});

peerServer.on('connection', (client) => { console.log(`[connect] client: ${client.getId()} time: ${new Date().toLocaleString()}`) });
peerServer.on('disconnect', (client) => { console.log(`[disconnect] client: ${client.getId()} time: ${new Date().toLocaleString()}`) });

app.use("/peerjs", peerServer);

app.get("/", (req, res) => {
    res.status(200).send("Hello World");
});
server.listen(3000);
