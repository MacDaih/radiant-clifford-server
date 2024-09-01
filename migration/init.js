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
