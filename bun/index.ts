import { networkInterfaces } from "os";
import data from "../data.json";

const PORT = 3000;

const nets = networkInterfaces();
const localIP = nets.en0?.find((n) => n.family === "IPv4")?.address;

const server = Bun.serve({
  port: PORT,
  fetch(req) {
    const start = Bun.nanoseconds();
    console.log(
      `${new Date().toISOString()} ${req.headers.get("User-Agent")} - ${req.method} ${req.url} ${Bun.nanoseconds() - start}`,
    );

    if (new URL(req.url).pathname === "/") {
      return Response.json(data);
    }

    return new Response("Page not found", { status: 404 });
  },
});

console.log(`\n   Local:   ${server.url}`);
console.log(`   Network: http://${localIP}:${3000}/\n`);
