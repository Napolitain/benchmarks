const express = require('express');

const app = express();

app.get('/', (req, res) => {
  res.json({ message: 'Hello, World!' });
});

app.use((req, res) => {
  res.status(404).json({ error: 'Not Found' });
});

const PORT = 8080;
app.listen(PORT, () => {
  console.log(`Express server listening on :${PORT}`);
});
