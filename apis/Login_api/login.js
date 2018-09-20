const express = require('express');
const querystring = require('querystring');

const bodyParser = require('body-parser');
const app = express();
const http = require('http');

var callback = function(re) {
var str = '';
re.on('data', function(chunk) {
str += chunk;
});
re.on('end', function() {
console.log(str);
});
};


app.use(bodyParser.urlencoded({ extended: false }));

app.get('/', (request, response) =>  response.sendFile(`${__dirname}/login.html`));
app.get('/login', (request, response) =>  response.sendFile(`${__dirname}/login.html`));
app.get('/register', (request, response) =>  response.sendFile(`${__dirname}/login.html`));

app.post('/login', (req, response) => {
  const username = req.body.Username;
  const password = req.body.Password;
  console.log(username,"",password);
  response.writeHead(200, {'Content-Type': 'text/html'});
   
  var bodyString = JSON.stringify({
    Username: username,
    password: password
});

  var option1 = {
    host: "localhost",
    port: 8080,
    path: "/login",
    method: "POST",
    headers: {
        "Content-Type": "application/json",
        'Content-Length': Buffer.byteLength(bodyString)
    }
  };
  console.log(option1);
  var post_req=http.request(option1, callback).write(bodyString);
  
});

app.post('/register', (req, response) => {
  const username = req.body.Username;
  const password = req.body.Password;
  console.log(username,"",password);
  response.writeHead(200, {'Content-Type': 'text/html'});
   
  var bodyString = JSON.stringify({
    Username: username,
    password: password
});

  var option2 = {
    host: "localhost",
    port: 8080,
    path: "/register",
    method: "POST",
    headers: {
        "Content-Type": "application/json",
        'Content-Length': Buffer.byteLength(bodyString)

    }
  };
  console.log(option2);
  var post_req=http.request(option2, callback).write(bodyString);
  
});

app.listen(8081, () => console.info('Application running on port 8081'));
