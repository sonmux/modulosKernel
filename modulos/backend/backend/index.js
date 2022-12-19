const express = require('express');
const mysql = requiere('mysql');

const db = mysql.createConnection({
    host: 'localhost:3306',
    user: 'root',
    password: 'password',
    database: 'tarea2'
})

const app = express();



app.get("/usuarios", (req, res) => {
    db.query(`SELECT * FROM usuarios`, (error, data) => {
     if (error) {
       return res.json({ status: "ERROR", error });
     }
     return res.json(data);
   });
  });

  app.get("/registrados", (req, res) => {
    db.query(`SELECT nombre FROM usuarios`, (error, data) => {
     if (error) {
       return res.json({ status: "ERROR", error });
     }
     return res.json(data);
   });
  });


app.listen(3000, () => {
    console.log("servidor escuchando en el puerto 3000")
});