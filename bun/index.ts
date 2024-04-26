import { networkInterfaces } from "os";
import data from "../data.json";

const RESET = "\x1b[0m";
const style = (c: string) => (s: string) => c + s + RESET;
const bold = style(`\x1b[1m`);
const blue = style(`\x1b[34m`);
const cyan = style(`\x1b[36m`);
const magenta = style(`\x1b[35m`);
// const green = style(`\x1b[32m`);
// const red = style(`\x1b[31m`);
// const yellow = style(`\x1b[33m`);
// const white = style(`\x1b[37m`);

const PORT = 3000;

const nets = networkInterfaces();
const localIP = nets.en0?.find((n) => n.family === "IPv4")?.address;

const server = Bun.serve({
  port: PORT,
  fetch(req) {
    const start = Bun.nanoseconds();
    console.log(
      `${new Date().toISOString()} ${magenta(req.headers.get("User-Agent"))} ${cyan(req.method)} ${req.url} ${Bun.nanoseconds() - start}ns`,
    );

    if (new URL(req.url).pathname === "/") {
      return Response.json(data);
    }

    return new Response("Page not found", { status: 404 });
  },
});

const localNetwork = `http://${localIP}:${3000}/`;
console.log(`\n   ${bold("Local:")}   ${blue(server.url)}`);
console.log(`   ${bold("Network:")} ${blue(localNetwork)}\n`);
