const commander = require('commander')
const grpc = require('grpc')
const { range } = require('lodash')
const { resolve } = require('path')

const port = 9090

const descriptor = grpc.load(resolve(__dirname, 'hellow.proto'))
const Hellow = descriptor.com.github.nlepage.grpc_hellow.service.Hellow

const sayHellow = (name, count, { host }) => {
  const client = new Hellow(`${host}:${port}`, grpc.credentials.createInsecure())
  client.sayHellow({ name, count }, (err, { message } = {}) => {
    if (err) {
      console.error(err)
      return
    }
    console.log(message)
  })
}

const serve = () => {
  const server = new grpc.Server()
  server.addService(Hellow.service, {
    sayHellow: ({ request: { name, count } }, callback) => callback(null, {
      message: `Hellow ${name}, this is node ! ${range(1, Number(count) + 1).join(' ')}`
    }),
  })
  server.bind(`0.0.0.0:${port}`, grpc.ServerCredentials.createInsecure());
  server.start()
}

commander
  .name('hellow')
  .description('hellow is a gRPC hello world')

commander
  .command('client [name] [count]')
  .option('--host [host]', 'Host to call', 'localhost')
  .description('invokes the hellow client')
  .action(sayHellow)

commander
  .command('server')
  .description('starts the hellow server')
  .action(serve)

commander.parse(process.argv)
