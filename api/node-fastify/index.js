const fastify = require('fastify')({ logger: false });

fastify.get('/', async (request, reply) => {
  return { message: 'Hello, World!' };
});

fastify.setNotFoundHandler(async (request, reply) => {
  reply.code(404);
  return { error: 'Not Found' };
});

const PORT = 8080;
fastify.listen({ port: PORT, host: '0.0.0.0' }, (err) => {
  if (err) {
    console.error(err);
    process.exit(1);
  }
  console.log(`Fastify server listening on :${PORT}`);
});
