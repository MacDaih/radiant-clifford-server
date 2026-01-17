use radiant;
db.createCollection("report");

db.createUser({
    user: "clifford",
    pwd: "1234567",
    roles: [{
        role: "readWrite",
        db: "radiant"
    }],
    mechanisms: ["SCRAM-SHA-256"],
});

db.createUser({
    user: "root",
    pwd: "root",
    roles: [{ role: "userAdminAnyDatabase", db: "admin" }, "readWriteAnyDatabase"]
});

