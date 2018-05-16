const commander = require('commander')
const grpc = require('grpc')
const { range } = require('lodash')
const { resolve } = require('path')

const port = 9090

const descriptor = grpc.load(resolve(__dirname, 'hellow.proto'))
const Hellow = descriptor.com.github.nlepage.grpc_hellow.service.Hellow

const sayHellow = (name, count, { host, stream: isStream }) => {
  const client = new Hellow(`${host}:${port}`, grpc.credentials.createInsecure())
  if (isStream) {
    const stream = client.streamHellow({ name, count })
    stream.on('data', ({ message }) => console.log(message))
    stream.on('error', err => console.error(err))
  } else {
    client.sayHellow({ name, count }, (err, { message } = {}) => {
      if (err) {
        console.error(err)
        return
      }
      console.log(message)
    })
  }
}

const sleep = async time => new Promise(resolve => setTimeout(resolve, time))

const serve = () => {
  const server = new grpc.Server()
  server.addService(Hellow.service, {
    sayHellow: ({ request: { name, count } }, callback) => callback(null, {
      message: Array(Number(count)).fill(`Hellow ${name}, this is node !`).join('\n')
    }),
    streamHellow: async stream => {
      const message = `Hellow ${stream.request.name}, this is node !`
      for (let i = 0; i < stream.request.count; i++) {
        await sleep(1000)
        stream.write(message)
      }
      stream.end()
    }
  })
  server.bind(`0.0.0.0:${port}`, grpc.ServerCredentials.createInsecure())
  server.start()
}

commander
  .name('hellow')
  .description('hellow is a gRPC hello world')

commander
  .command('client <name> <count>')
  .option('--host [host]', 'Host to call', 'localhost')
  .option('-s, --stream', 'Ask for a streamed response', false)
  .description('invokes the hellow client')
  .action(sayHellow)

commander
  .command('server')
  .description('starts the hellow server')
  .action(serve)

commander.parse(process.argv)
