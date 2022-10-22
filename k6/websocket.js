import ws from 'k6/ws';
import { check } from 'k6';
import { uuidv4 } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';
import { sleep } from 'k6';

export default function () {
  const params = { tags: { my_tag: 'hello' } };
  const randomUUID = uuidv4();
  const url = 'ws://localhost/test/ws/' + randomUUID;

  const res = ws.connect(url, null, function (socket) {
    socket.on('open', () =>  {
        // Set a timeout to close
        socket.setInterval(function timeout() {
            socket.close();
          }, 1000);

        socket.send(Date.now() + ": Hi from " + randomUUID);
    });

    // For testing, close when we get a message
    socket.on('message', (data) => {
        // sleep 10 seconds 
        sleep(10);
        socket.close();
    });
    
   // socket.close();
  });

  check(res, { 'status is 101': (r) => r && r.status === 101 });
}
